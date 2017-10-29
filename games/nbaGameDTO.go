package games

import (
	"time"
)

type NBAGameDTO struct {
	ID     int     `json:"id"`
	Date   string  `json:"date"`
	TimeEt string  `json:"time_et"`
	Status string  `json:"status"`
	Team1  string  `json:"team1"`
	Team2  string  `json:"team2"`
	Score1 int     `json:"score1"`
	Score2 int     `json:"score2"`
	Prob1  float64 `json:"carmelo_prob1"`
	Prob2  float64 `json:"carmelo_prob2"`
}

func (g NBAGameDTO) ToGame() Game {
	dt := g.Date + " " + g.TimeEt
	loc, _ := time.LoadLocation("America/New_York")
	layout := "2006-01-02 15:04"
	date, _ := time.ParseInLocation(layout, dt, loc)

	return Game{
		ID:       g.ID,
		DateTime: date.UTC(),
		Status:   g.Status,
		Team1:    g.Team1,
		Team2:    g.Team2,
		Prob1:    g.Prob1,
		Prob2:    g.Prob2,
		Score1:   g.Score1,
		Score2:   g.Score2,
		League:   "NBA",
	}
}
