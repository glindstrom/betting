package main

import (
	"encoding/json"
	"github.com/glindstrom/betting/games"
	"io/ioutil"
	"log"
	"net/http"
)

type wrapper struct {
	Games []games.GameDTO `json:"games"`
}

func main() {
	resp, err := http.Get("https://projects.fivethirtyeight.com/2018-nba-predictions/data.json")
	if err != nil {
		log.Fatalln("http get error:", err)
	}
	defer resp.Body.Close()

	buf, _ := ioutil.ReadAll(resp.Body)

	var w wrapper
	err = json.Unmarshal(buf, &w)
	checkErr(err)

	v := make([]games.Game, 0, len(w.Games))
	for _, value := range w.Games {
		v = append(v, value.ToGame("NBA"))
	}

	games.PrintGames(v, games.IsToday)

}

func checkErr(err error) {
	if err != nil {
		log.Fatalln("json decode error:", err)
	}
}
