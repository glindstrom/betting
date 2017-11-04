package games

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type NFLGameDTO struct {
	ID       int       `json:"id"`
	DateTime time.Time `json:"datetime"`
	Status   string    `json:"status"`
	Team1    string    `json:"team1"`
	Team2    string    `json:"team2"`
	Score1   *int      `json:"score1"`
	Score2   *int      `json:"score2"`
	Prob1    float64   `json:"prob1"`
	Prob2    float64   `json:"prob2"`
}

func (g NFLGameDTO) ToGame() Game {
	return Game{
		ID:       bson.NewObjectId(),
		ID538:    g.ID,
		DateTime: g.DateTime,
		Status:   g.Status,
		Team1:    g.Team1,
		Team2:    g.Team2,
		Prob1:    g.Prob1,
		Prob2:    g.Prob2,
		Score1:   g.Score1,
		Score2:   g.Score2,
		League:   "NFL",
	}
}
