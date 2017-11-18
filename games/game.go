package games

import (
	"fmt"
	"github.com/glindstrom/betting/db"
	"gopkg.in/mgo.v2/bson"
	"math/big"
	"strconv"
	"strings"
	"time"
	"math"
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
	OfferedOdds1 float64    `json:"offeredOdds1" bson:"offeredOdds1,omitempty"`
	OfferedOdds2 float64    `json:"offeredOdds2" bson:"offeredOdds2,omitempty"`
}

func (g Game) Time() time.Time {
	return g.DateTime.Local()
}

func (g Game) Odds1() float64 {
	return floatOrZero(odds(g.Prob1))
}

func (g Game) Odds2() float64 {
	return floatOrZero(odds(g.Prob2))
}

func (g Game) Odds1String() string {
	return floatToString(odds(g.Prob1))
}

func (g Game) Odds2String() string {
	return floatToString(odds(g.Prob2))
}

func odds(p float64) float64 {
	x, y := big.NewFloat(1), big.NewFloat(p)
	res, _ := new(big.Float).Quo(x, y).Float64()
	return res
	//return strings.Replace(fmt.Sprintf("%.2f", res), ".", ",", 1)
}


func (m Game) PrintCSV() {
	t := m.Time().Format("02.01 15:04")
	fmt.Println(t+";", m.League+";", m.Team2+";", m.Team1+";", m.Odds2String()+";", m.Odds1String()+";", floatToString(m.Prob2)+";", floatToString(m.Prob1))
}

func floatToString(f float64) string {
	return strings.Replace(strconv.FormatFloat(f, 'f', -1, 64), ".", ",", 1)
}

func AllGames() ([]Game, error) {
	var gms []Game
	err := db.Games.Find(bson.M{}).All(&gms)
	if err != nil {
		return nil, err
	}
	return gms, nil
}

func OneGame(id bson.ObjectId ) (Game, error) {
	var g Game
	err := db.Games.Find(bson.M{"_id": id}).One(&g)
	if err != nil {
		return g, err
	}
	return g, nil
}

func UpdateGame(g Game) (error) {
	return db.Games.Update(bson.M{"_id": g.ID}, &g)
}

func floatOrZero(f float64) float64 {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0
	}
	return f
}
