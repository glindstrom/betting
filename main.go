package main

import (
	"github.com/glindstrom/betting/games"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func main() {
	r := httprouter.New()
	gc := games.GameController{}
	r.GET("/games", gc.GetUpcomingGames)
	r.POST("/game", gc.UpdateGame)
	r.POST("/games/bet", gc.Bet)
	r.OPTIONS("/game", gc.Options)
	r.OPTIONS("/games/bet", gc.Options)
	http.ListenAndServe("localhost:8080", r)
}
