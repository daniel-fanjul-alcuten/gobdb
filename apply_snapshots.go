package gobdb

import (
	"errors"
	"io"
	"sort"
)

type SnapshotIdSlice []SnapshotId

func (s SnapshotIdSlice) Len() int {
	return len(s)
}

func (s SnapshotIdSlice) Less(i, j int) bool {
	return s[i].Id() < s[j].Id()
}

func (s SnapshotIdSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// It selects the best snapshot and applies the Writers to the Root.
func ApplySnapshot(root Root, repository SnapshotRepository) (id TransactionId, err error) {

	snapshots := repository.Snapshots()
	if len(snapshots) == 0 {
		return
	}

	sort.Sort(SnapshotIdSlice(snapshots))
	snapshotId := snapshots[len(snapshots)-1]
	id = snapshotId.Id()

	var reader SnapshotReader
	reader, err = repository.ReadSnapshot(snapshotId)
	if err != nil {
		return
	}
	defer reader.Close()

	for {
		var writer Writer
		writer, err = reader.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
		if writer == nil {
			err = errors.New("gobdb: decoded nil Writer on ApplySnapshot")
			return
		}
		_, err = writer.Write(root)
		if err != nil {
			return
		}
	}
}
