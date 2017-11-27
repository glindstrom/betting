package games

import (
	"github.com/glindstrom/betting/db"
	"gopkg.in/mgo.v2/bson"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const (
	maxBetFraction = 0.05
	kellyFraction  = 0.25
)

type Game struct {
	ID              bson.ObjectId `json:"id" bson:"_id"`
	ID538           int           `json:"id538" bson:"id538"`
	DateTime        time.Time     `json:"dateTime" bson:"dateTime"`
	Status          string        `json:"status" bson:"status"`
	Team1           string        `json:"team1" bson:"team1"`
	Team2           string        `json:"team2" bson:"team2"`
	Score1          *int          `json:"score1" bson:"score1,omitempty"`
	Score2          *int          `json:"score2" bson:"score2,omitempty"`
	Prob1           float64       `json:"prob1" bson:"prob1"`
	Prob2           float64       `json:"prob2" bson:"prob2"`
	League          string        `json:"league" bson:"league"`
	OfferedOdds1    float64       `json:"offeredOdds1" bson:"offeredOdds1,omitempty"`
	OfferedOdds2    float64       `json:"offeredOdds2" bson:"offeredOdds2,omitempty"`
	PredictedWinner string        `json:"predictedWinner" bson:"predictedWinner,omitempty"`
	BetAmount       float64       `json:"betAmount" bson:"betAmount,omitempty"`
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

func (g Game) BetPlaced() bool {
	return g.PredictedWinner != ""
}

func odds(p float64) float64 {
	x, y := big.NewFloat(1), big.NewFloat(p)
	res, _ := new(big.Float).Quo(x, y).Float64()
	return res
}

func (g Game) OptimalBetSize() float64 {
	if !shouldBetOnGame(g) {
		return 0
	}

	var fraction float64
	if shouldBet1(g) {
		fraction = fractionToBet(g.OfferedOdds1, g.Prob1)
	} else {
		fraction = fractionToBet(g.OfferedOdds2, g.Prob2)
	}

	betSize := new(big.Float).Mul(big.NewFloat(fraction), big.NewFloat(bankroll()))
	betSize = new(big.Float).Mul(betSize, big.NewFloat(kellyFraction))
	betSizeAsFloat, _ := betSize.Float64()
	return math.Min(betSizeAsFloat, maxBetFraction*bankroll())
}

func fractionToBet(offeredOdds float64, prob float64) float64 {
	bp := new(big.Float).Mul(big.NewFloat(offeredOdds-1), big.NewFloat(prob))
	bpp := new(big.Float).Sub(bp, big.NewFloat(1-prob))
	f, _ := new(big.Float).Quo(bpp, big.NewFloat(offeredOdds-1)).Float64()
	return f
}

func shouldBetOnGame(g Game) bool {
	return shouldBet1(g) || shouldBet2(g)
}

func shouldBet1(g Game) bool {
	return shouldBet(g.Odds1(), g.OfferedOdds1)
}

func shouldBet2(g Game) bool {
	return shouldBet(g.Odds2(), g.OfferedOdds2)
}

func shouldBet(oddsPredicted float64, oddsOffered float64) bool {
	if oddsPredicted == 0 || oddsOffered == 0 {
		return false
	}
	return oddsOffered > oddsPredicted
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

func UpcomingGames() ([]Game, error) {

	tomorrow := time.Now().Add(48 * time.Hour).UTC().Truncate(24 * time.Hour)
	var gms []Game
	err := db.Games.Find(bson.M{"status": "pre", "dateTime": bson.M{"$lte": tomorrow}}).All(&gms)
	if err != nil {
		return nil, err
	}
	return gms, nil
}

func OneGame(id bson.ObjectId) (Game, error) {
	var g Game
	err := db.Games.Find(bson.M{"_id": id}).One(&g)
	if err != nil {
		return g, err
	}
	return g, nil
}

func UpdateGame(g Game) error {
	return db.Games.Update(bson.M{"_id": g.ID}, &g)
}

func floatOrZero(f float64) float64 {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0
	}
	return f
}

func bankroll() float64 {
	return 69.71
}
