package games

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const (
	GamesCollection = "games"
)

type Game struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	ID538    int           `json:"id538" bson:"id538"`
	DateTime time.Time     `json:"dateTime" bson:"dateTime"`
	Status   string        `json:"status" bson:"status"`
	Team1    string        `json:"team1" bson:"team1"`
	Team2    string        `json:"team2" bson:"team2"`
	Score1   *int          `json:"score1" bson:"score1,omitempty"`
	Score2   *int          `json:"score2" bson:"score2,omitempty"`
	Prob1    float64       `json:"prob1" bson:"prob1"`
	Prob2    float64       `json:"prob2" bson:"prob2"`
	League   string        `json:"league" bson:"league"`
}

func (g Game) Time() time.Time {
	return g.DateTime.Local()
}

func (g Game) Odds1() string {
	return odds(g.Prob1)
}

func (g Game) Odds2() string {
	return odds(g.Prob2)
}

func odds(p float64) string {
	x, y := big.NewFloat(1), big.NewFloat(p)
	res, _ := new(big.Float).Quo(x, y).Float64()
	return strings.Replace(fmt.Sprintf("%.2f", res), ".", ",", 1)
}

func (m Game) PrintCSV() {
	t := m.Time().Format("02.01 15:04")
	fmt.Println(t+";", m.League+";", m.Team2+";", m.Team1+";", m.Odds2()+";", m.Odds1()+";", floatToString(m.Prob2)+";", floatToString(m.Prob1))
}

func floatToString(f float64) string {
	return strings.Replace(strconv.FormatFloat(f, 'f', -1, 64), ".", ",", 1)
}
