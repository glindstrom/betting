package games

import (
	"sort"
	"strings"
)

type mSort []Game

func (s mSort) Len() int {
	return len(s)
}

func (s mSort) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s mSort) Less(i, j int) bool {
	if s[i].DateTime.Before(s[j].DateTime) {
		return true
	}
	if s[i].DateTime.After(s[j].DateTime) {
		return false
	}
	l1 := s[i].League
	l2 := s[j].League
	if l1 != l2 {
		return strings.Compare(l1, l2) == -1
	}
	return strings.Compare(s[i].Team1, s[j].Team1) == -1
}

func SortGames(games []Game) {
	sort.Sort(mSort(games))
}
