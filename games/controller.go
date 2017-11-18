package games

import (
	"encoding/json"
	"fmt"
	"github.com/glindstrom/betting/db"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"time"
	"gopkg.in/mgo.v2/bson"
)

type GameController struct{}

type JsonGame struct {
	ID bson.ObjectId `json:"id"`
	DateTime time.Time `json:"dateTime"`
	League string `json:"league"`
	Team1 string `json:"team1"`
	Team2 string `json:"team2"`
	Odds1 float64 `json:"odds1"`
	Odds2 float64 `json:"odds2"`
	OfferedOdds1 float64 `json:"offeredOdds1"`
	OfferedOdds2 float64 `json:"offeredOdds2"`
}

func (gc GameController) GetGames(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	saveAllGames(db.Games)

	gms, err := AllGames()
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Failed get all games: ", err)
		return
	}

	respBody, err := json.Marshal(gamesAsJson(gms))
	if err != nil {
		log.Fatal(err)
	}

	ResponseWithJSON(w, respBody, http.StatusOK)
}

func (gc GameController) UpdateGame(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var jg JsonGame
	err := decoder.Decode(&jg)
	if err != nil {
		ErrorWithJSON(w, "Error decoding JSON", http.StatusInternalServerError)
		return
	}
	g, err := OneGame(jg.ID)
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		return
	}
	g.OfferedOdds1 = jg.OfferedOdds1
	g.OfferedOdds2 = jg.OfferedOdds2
	err = UpdateGame(g)
	if err != nil {
		ErrorWithJSON(w, "Error updating game", http.StatusInternalServerError)
		return
	}
	ErrorWithJSON(w, "OK", http.StatusOK)
}

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(json)
}

func gamesAsJson(gms []Game) []JsonGame {
	var rs []JsonGame
	for _,g := range gms {
		if !IsTodayOrTomorrow(g) {
			continue
		}
		r := JsonGame{
			DateTime: g.DateTime,
			League: g.League,
			Team1:g.Team1,
			Team2:g.Team2,
			Odds1:g.Odds1(),
			Odds2:g.Odds2(),
		}
		rs = append(rs, r)
	}
	return rs
}

func saveAllGames(c *mgo.Collection) {
	for _, g := range fetchAllGames() {
		var oldGame Game
		err := c.Find(bson.M{"id538": g.ID538}).One(&oldGame)
		if err != nil && g.Status == "pre" {
			c.Insert(g)
		} else {
			oldGame.DateTime = g.DateTime
			oldGame.Score1 = g.Score1
			oldGame.Score2 = g.Score2
			oldGame.Prob1 = g.Prob1
			oldGame.Prob2 = g.Prob2
			c.Update(bson.M{"id538": g.ID538}, &oldGame)
		}
	}
}
