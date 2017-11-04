package matches

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
)

type Match struct {
	ID         int    `json: "id"`
	LeagueID   int    `json:"league_id"`
	LeagueName string `json: ",omitempty"`
	Country    string `json: ",omitempty"`
	Team1      string
	Team2      string
	Prob1      float64
	Prob2      float64
	Probtie    float64
	DateTime   time.Time
	Score1     *int `json: ",omitempty"`
	Score2     *int `json: ",omitempty"`
}

func (m Match) Odds1() string {
	return odds(m.Prob1)
}

func (m Match) Odds2() string {
	return odds(m.Prob2)
}

func (m Match) OddsTie() string {
	return odds(m.Probtie)
}

func (m Match) Time() time.Time {
	return m.DateTime.Local()
}

func odds(p float64) string {
	x, y := big.NewFloat(1), big.NewFloat(p)
	res, _ := new(big.Float).Quo(x, y).Float64()
	return strings.Replace(fmt.Sprintf("%.2f", res), ".", ",", 1)
}

func (m Match) PrintCSV() {
	t := m.Time().Format("02.01 15:04")
	fmt.Println(m.Country+";", t+";", m.Team1+";", m.Team2+";", m.Odds1()+";", m.OddsTie()+";", m.Odds2()+";", floatToString(m.Prob1)+";", floatToString(m.Probtie)+";", floatToString(m.Prob2))
}

func floatToString(f float64) string {
	return strings.Replace(strconv.FormatFloat(f, 'f', -1, 64), ".", ",", 1)
}
