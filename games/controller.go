package games

import (
	"encoding/json"
	"fmt"
	"github.com/glindstrom/betting/db"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

type GameController struct{}

func (gc GameController) GetGames(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := db.Session.Copy()
	defer session.Close()
	c := session.DB(db.DataBase).C(GamesCollection)
	saveAllGames(c)

	var games []Game
	err := c.Find(bson.M{}).All(&games)
	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Failed get all games: ", err)
		return
	}

	respBody, err := json.Marshal(games)
	if err != nil {
		log.Fatal(err)
	}

	ResponseWithJSON(w, respBody, http.StatusOK)
}

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

func saveAllGames(c *mgo.Collection) {
	for _, g := range fetchAllGames() {
		if g.Status == "pre" {
			c.Insert(g)
		}
	}
}
