package games

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
)

type Game struct {
	ID       int `json:"id"`
	DateTime time.Time
	Status   string `json:"status"`
	Team1    string `json:"team1"`
	Team2    string `json:"team2"`
	Score1   int    `json:"score1"`
	Score2   int    `json:"score2"`
	Prob1    float64
	Prob2    float64
	League   string
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

func (g Game) PrintCSV() {
	t := g.Time().Format("02.01 15:04")
	fmt.Println(t+";", g.League+";", g.Team1+";", g.Team2+";", g.Odds1()+";", g.Odds2()+";", floatToString(g.Prob1)+";", floatToString(g.Prob2))
}

func floatToString(f float64) string {
	return strings.Replace(strconv.FormatFloat(f, 'f', -1, 64), ".", ",", 1)
}
