package main

import (
	"encoding/json"
	"github.com/glindstrom/betting/leagues"
	"github.com/glindstrom/betting/matches"
	"io/ioutil"
	"log"
	"net/http"
)

type wrapper struct {
	Leagues  []leagues.League `json:"leagues"`
	Matches  []matches.Match  `json:"matches"`
	Matches2 []matches.Match  `json:"matches2"`
}

var mLeague map[int]leagues.League

var mMatch map[int]matches.Match

func main() {
	resp, err := http.Get("https://projects.fivethirtyeight.com/soccer-predictions/data.json")
	if err != nil {
		log.Fatalln("http get error:", err)
	}
	defer resp.Body.Close()

	buf, _ := ioutil.ReadAll(resp.Body)

	var w wrapper
	err = json.Unmarshal(buf, &w)
	checkErr(err)

	mLeague = make(map[int]leagues.League)
	addToLeagueMap(w.Leagues)

	mMatch = make(map[int]matches.Match)
	addMatchesToMap(w.Matches)
	addMatchesToMap(w.Matches2)

	v := make([]matches.Match, 0, len(mMatch))
	for _, value := range mMatch {
		value.LeagueName = mLeague[value.LeagueID].LongName
		value.Country = mLeague[value.LeagueID].Country
		v = append(v, value)
	}

	matches.PrintMatches(v, matches.IsTomorrow)
	//printMatches(v, matches.IsToday)

}
func addMatchesToMap(matches []matches.Match) {
	for _, m := range matches {
		mMatch[m.ID] = m
	}
}

func addToLeagueMap(ls []leagues.League) {
	for _, l := range ls {
		mLeague[l.ID] = l
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln("json decode error:", err)
	}
}
