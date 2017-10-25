package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type wrapper struct {
	Leagues  []league `json:"leagues"`
	Matches  []match  `json:"matches"`
	Matches2 []match  `json:"matches2"`
}

type league struct {
	LongName string
	Country  string
	ID       int
}

type match struct {
	ID       int `json: "id"`
	LeagueID int `json:"league_id"`
	Team1    string
	Team2    string
	Prob1    float64
	Prob2    float64
	Probtie  float64
	DateTime time.Time
}

func (m match) odds1() string {
	return odds(m.Prob1)
}

func (m match) odds2() string {
	return odds(m.Prob2)
}

func (m match) oddsTie() string {
	return odds(m.Probtie)
}

func (m match) time() time.Time {
	return m.DateTime.Local()
}

func odds(p float64) string {
	x, y := big.NewFloat(1), big.NewFloat(p)
	res, _ := new(big.Float).Quo(x, y).Float64()
	return strings.Replace(fmt.Sprintf("%.2f", res), ".", ",", 1)
}

var mLeague map[int]league

var mMatch map[int]match

type matchesSort []match

func (s matchesSort) Len() int {
	return len(s)
}

func (s matchesSort) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s matchesSort) Less(i, j int) bool {
	if s[i].time().Before(s[j].time()) {
		return true
	}
	if s[i].time().After(s[j].time()) {
		return false
	}
	c1 := mLeague[s[i].LeagueID].Country
	c2 := mLeague[s[j].LeagueID].Country
	if c1 != c2 {
		return strings.Compare(c1, c2) == -1
	}
	return strings.Compare(s[i].Team1, s[j].Team1) == -1
}

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

	mLeague = make(map[int]league)
	addToLeagueMap(w.Leagues)

	mMatch = make(map[int]match)
	addMatchesToMap(w.Matches)
	addMatchesToMap(w.Matches2)

	v := make([]match, 0, len(mMatch))
	for _, value := range mMatch {
		v = append(v, value)
	}

	// printMatches(v, isTomorrow)
	printMatches(v, isToday)

}
func addMatchesToMap(matches []match) {
	for _, m := range matches {
		mMatch[m.ID] = m
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln("json decode error:", err)
	}
}

func printMatches(matches []match, fn shouldPrint) {
	sort.Sort(matchesSort(matches))
	for _, m := range matches {
		if fn(m) {
			printMatch(m)
		}
	}
}

type shouldPrint func(m match) bool

func isToday(m match) bool {
	t1 := truncateDate(time.Now())
	t2 := truncateDate(m.time())
	return t1.Equal(t2) && m.time().After(time.Now())
}

func isTomorrow(m match) bool {
	tomorrow := time.Now().Add(24 * time.Hour)
	return truncateDate(tomorrow).Equal(truncateDate(m.time()))
}

func printMatch(m match) {
	t := m.time().Format("02.01 15:04")
	fmt.Println(mLeague[m.LeagueID].Country+";", t+";", m.Team1+";", m.Team2+";", m.odds1()+";", m.oddsTie()+";", m.odds2()+";", floatToString(m.Prob1)+";", floatToString(m.Probtie)+";", floatToString(m.Prob2))

}

func floatToString(f float64) string {
	return strings.Replace(strconv.FormatFloat(f, 'f', -1, 64), ".", ",", 1)
}

func addToLeagueMap(ls []league) {
	for _, l := range ls {
		mLeague[l.ID] = l
	}
}

func truncateDate(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}
