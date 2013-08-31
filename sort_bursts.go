package gobdb

import (
	"sort"
)

// It implements sort.Interface.
type burstIdSlice []BurstId

func (s burstIdSlice) Len() int {
	return len(s)
}

func (s burstIdSlice) Less(i, j int) bool {
	diff := s[i].First() - s[j].First()
	if diff != 0 {
		return diff < 0
	}
	return s[i].Last() > s[j].Last()
}

func (s burstIdSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Sort BurstIds in ascending order of their First() and then descending order of their Last().
func SortBursts(bursts []BurstId) {
	sort.Sort(burstIdSlice(bursts))
}
