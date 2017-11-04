package games

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	nbaDataURL = "https://projects.fivethirtyeight.com/2018-nba-predictions/data.json"
	nflDataURL = "https://projects.fivethirtyeight.com/2017-nfl-predictions/data.json"
	mlbDataURL = "https://projects.fivethirtyeight.com/2017-mlb-predictions/data.json"
)

type nbaGameWrapper struct {
	Games []NBAGameDTO `json:"games"`
}

type nflGameWrapper struct {
	Games []NFLGameDTO `json:"games"`
}

type mlbGameWrapper struct {
	Games []MLBGameDTO `json:"games"`
}

func PrintGames(games []Game, fn filter) {
	SortGames(games)
	for _, g := range games {
		if fn(g) {
			g.PrintCSV()
		}
	}
}

func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	return s
}

type filter func(g Game) bool

func IsToday(g Game) bool {
	t1 := truncateDate(time.Now())
	t2 := truncateDate(g.Time())
	return t1.Equal(t2) && g.Time().After(time.Now())
}

func IsTomorrow(g Game) bool {
	tomorrow := time.Now().Add(24 * time.Hour)
	return truncateDate(tomorrow).Equal(truncateDate(g.Time()))
}

func truncateDate(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func PrintAllGames() {
	PrintGames(fetchAllGames(), IsToday)
}

func fetchAllGames() []Game {
	allGames := make([]GameDTO, 0)
	allGames = append(allGames, importNBAGames()...)
	allGames = append(allGames, importNFLGames()...)
	allGames = append(allGames, importMLBGames()...)
	v := make([]Game, 0, len(allGames))
	for _, value := range allGames {
		g := value.ToGame()
		v = append(v, g)
	}
	return v
}

func importNBAGames() []GameDTO {
	buf := fetchGames(nbaDataURL)
	var nba nbaGameWrapper
	err := json.Unmarshal(buf, &nba)
	checkErr(err)
	gs := make([]GameDTO, 0, len(nba.Games))
	for _, value := range nba.Games {
		gs = append(gs, value)
	}
	return gs
}

func importNFLGames() []GameDTO {
	buf := fetchGames(nflDataURL)
	var nfl nflGameWrapper
	err := json.Unmarshal(buf, &nfl)
	checkErr(err)
	gs := make([]GameDTO, 0, len(nfl.Games))
	for _, value := range nfl.Games {
		gs = append(gs, value)
	}
	return gs
}

func importMLBGames() []GameDTO {
	buf := fetchGames(mlbDataURL)
	var mlb mlbGameWrapper
	err := json.Unmarshal(buf, &mlb)
	checkErr(err)
	gs := make([]GameDTO, 0, len(mlb.Games))
	for _, value := range mlb.Games {
		gs = append(gs, value)
	}
	return gs
}

func fetchGames(s string) []byte {
	resp, err := http.Get(s)
	if err != nil {
		log.Fatalln("http get error:", err)
	}
	defer resp.Body.Close()

	buf, _ := ioutil.ReadAll(resp.Body)
	return buf
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln("json decode error:", err)
	}
}
