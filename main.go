package main

import (
	"encoding/json"
	"github.com/glindstrom/betting/games"
	"io/ioutil"
	"log"
	"net/http"
)

type nbaGameWrapper struct {
	Games []games.NBAGameDTO `json:"games"`
}

type nflGameWrapper struct {
	Games []games.NFLGameDTO `json:"games"`
}

func main() {
	resp, err := http.Get("https://projects.fivethirtyeight.com/2018-nba-predictions/data.json")
	if err != nil {
		log.Fatalln("http get error:", err)
	}
	defer resp.Body.Close()

	buf, _ := ioutil.ReadAll(resp.Body)

	var nba nbaGameWrapper
	err = json.Unmarshal(buf, &nba)
	checkErr(err)
	allGames := make([]games.GameDTO, 0)
	for _, value := range nba.Games {
		allGames = append(allGames, value)
	}

	resp, err = http.Get("https://projects.fivethirtyeight.com/2017-nfl-predictions/data.json")
	if err != nil {
		log.Fatalln("http get error:", err)
	}
	defer resp.Body.Close()

	buf, _ = ioutil.ReadAll(resp.Body)

	var nfl nflGameWrapper
	err = json.Unmarshal(buf, &nfl)
	checkErr(err)

	for _, value := range nfl.Games {
		allGames = append(allGames, value)
	}

	v := make([]games.Game, 0, len(allGames))
	for _, value := range allGames {
		v = append(v, value.ToGame())
	}

	games.PrintGames(v, games.IsToday)

}

func checkErr(err error) {
	if err != nil {
		log.Fatalln("json decode error:", err)
	}
}
