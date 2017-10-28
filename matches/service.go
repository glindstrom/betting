package matches

import "time"

func PrintMatches(ms []Match, fn shouldPrint) {
	SortMatches(ms)
	for _, m := range ms {
		if fn(m) {
			m.PrintCSV()
		}
	}
}

type shouldPrint func(m Match) bool

func IsToday(m Match) bool {
	t1 := truncateDate(time.Now())
	t2 := truncateDate(m.Time())
	return t1.Equal(t2) && m.Time().After(time.Now())
}

func IsTomorrow(m Match) bool {
	tomorrow := time.Now().Add(24 * time.Hour)
	return truncateDate(tomorrow).Equal(truncateDate(m.Time()))
}

func truncateDate(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}
