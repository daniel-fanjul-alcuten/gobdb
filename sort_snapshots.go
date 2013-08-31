package gobdb

import (
	"sort"
)

// It implements sort.Interface.
type snapshotIdSlice []SnapshotId

func (s snapshotIdSlice) Len() int {
	return len(s)
}

func (s snapshotIdSlice) Less(i, j int) bool {
	return s[i].Id() > s[j].Id()
}

func (s snapshotIdSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Sort SnapshotIds in descending order of their Id().
func SortSnapshots(snapshots []SnapshotId) {
	sort.Sort(snapshotIdSlice(snapshots))
}
