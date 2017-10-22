package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"
	"math/big"
	"strings"
	"strconv"
)

type leagues struct {
	Leagues []league
}

type matches struct {
	Matches []match
}

type matches2 struct {
	Matches []match
}

type league struct {
	LongName string
	Country  string
	ID       int
}

type match struct {
	LeagueID int `json:"league_id"`
	Team1    string
	Team2    string
	Prob1    float64
	Prob2    float64
	Probtie  float64
	DateTime time.Time
}

func(m match) odds1() string {
	return odds(m.Prob1)
}

func(m match) odds2() string {
	return odds(m.Prob2)
}

func(m match) oddsTie() string {
	return odds(m.Probtie)
}

func(m match) time() time.Time {
	return m.DateTime.Local()
}

func odds(p float64) string {
	x, y := big.NewFloat(1), big.NewFloat(p)
	res, _ := new(big.Float).Quo(x, y).Float64()
	return strings.Replace(fmt.Sprintf("%.2f", res), ".", ",", 1)
}

var mLeague map[int]league

type matchesSort matches

func (s matchesSort) Len() int {
	return len(s.Matches)
}

func (s matchesSort) Swap(i, j int) { s.Matches[i], s.Matches[j] = s.Matches[j], s.Matches[i] }

func (s matchesSort) Less(i, j int) bool { return s.Matches[i].time().Before(s.Matches[j].time()) }

func main() {
	resp, err := http.Get("https://projects.fivethirtyeight.com/soccer-predictions/data.json")
	if err != nil {
		log.Fatalln("http get error:", err)
	}
	defer resp.Body.Close()

	buf, _ := ioutil.ReadAll(resp.Body)

	var leagues leagues
	err = json.Unmarshal(buf, &leagues)
	checkErr(err)
	addToLeagueMap(leagues.Leagues)

	var matches1 matches
	err = json.Unmarshal(buf, &matches1)
	checkErr(err)

	var matches2 matches2
	err = json.Unmarshal(buf, &matches2)
	checkErr(err)

	fmt.Println(matches2)

	var allMatches matches
	//allMatches.Matches = append(allMatches.Matches, matches1.Matches...)
	allMatches.Matches = append(allMatches.Matches, matches2.Matches...)

	printTodaysMatches(allMatches)
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln("json decode error:", err)
	}
}

func printTodaysMatches(matches matches) {
	sort.Sort(matchesSort(matches))
	for _, m := range matches.Matches {
		if isToday(m) {
			printMatch(m)
		}
	}
}

func isToday(m match) bool {
	t1 := time.Now().Truncate(24 * time.Hour)
	t2 := m.time().Truncate(24 * time.Hour)
	return t1.Equal(t2) && m.time().After(time.Now())
}

func printMatch(m match) {
	t := m.time().Format("02.01 15:04")
	fmt.Println(mLeague[m.LeagueID].Country+";", t+";", m.Team1+";", m.Team2+";", m.odds1()+";", m.oddsTie()+";", m.odds2()+";", floatToString(m.Prob1)+";", floatToString(m.Probtie)+";", floatToString(m.Prob2))

}

func floatToString(f float64) string {
	return strings.Replace(strconv.FormatFloat(f, 'f', -1, 64), ".", ",", 1)
}

func addToLeagueMap(ls[]league) {
	mLeague = make(map[int]league)
	for _,l := range ls {
		mLeague[l.ID] = l
	}
}
