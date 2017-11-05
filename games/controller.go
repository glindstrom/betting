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
)

type GameController struct{}

type Response struct {
	DateTime time.Time `json:"dateTime"`
	League string `json:"league"`
	Team1 string `json:"team1"`
	Team2 string `json:"team2"`
	Odds1 float64 `json:"odds1"`
	Odds2 float64 `json:"odds2"`
}

func (gc GameController) GetGames(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(json)
}

func gamesAsJson(gms []Game) []Response {
	var rs []Response
	for _,g := range gms {
		if !IsTodayOrTomorrow(g) {
			continue
		}
		r := Response{
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
		if g.Status == "pre" {
			c.Insert(g)
		}
	}
}
