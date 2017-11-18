package main

import (
	"github.com/glindstrom/betting/games"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func main() {
	r := httprouter.New()
	gc := games.GameController{}
	r.GET("/games", gc.GetGames)
	r.POST("/game", gc.UpdateGame)
	http.ListenAndServe("localhost:8080", r)
}
