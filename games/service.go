package games

import (
	"time"
)

func PrintGames(games []Game, fn filter) {
	SortGames(games)
	for _, g := range games {
		if fn(g) {
			g.PrintCSV()
		}
	}
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
