package gobdb

import (
	"errors"
	"io"
)

// It selects the best snapshot and applies the Writers to the Root.
func ApplySnapshot(root Root, repository SnapshotRepository) (id TransactionId, err error) {

	snapshots := repository.Snapshots()
	if len(snapshots) == 0 {
		return
	}

	SortSnapshots(snapshots)
	snapshotId := snapshots[0]
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
