package main

import (
	"github.com/glindstrom/betting/db"
	"github.com/glindstrom/betting/games"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func init() {
	db.Connect()
}

func main() {
	r := httprouter.New()
	// Get a UserController instance
	gc := games.GameController{}
	r.GET("/games", gc.GetGames)
	http.ListenAndServe("localhost:8080", r)
}
