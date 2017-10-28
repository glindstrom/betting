package matches

import (
	"sort"
	"strings"
)

type mSort []Match

func (s mSort) Len() int {
	return len(s)
}

func (s mSort) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s mSort) Less(i, j int) bool {
	if s[i].Time().Before(s[j].Time()) {
		return true
	}
	if s[i].Time().After(s[j].Time()) {
		return false
	}
	c1 := s[i].Country
	c2 := s[j].Country
	if c1 != c2 {
		return strings.Compare(c1, c2) == -1
	}
	l1 := s[i].LeagueName
	l2 := s[j].LeagueName
	if l1 != l2 {
		return strings.Compare(l1, l2) == -1
	}
	return strings.Compare(s[i].Team1, s[j].Team1) == -1
}

func SortMatches(matches []Match) {
	sort.Sort(mSort(matches))
}
