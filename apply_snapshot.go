package gobdb

import (
	"errors"
	"io"
)

// It applies the Writers of a Snapshot to a Root.
func ApplySnapshot(root Root, snapshotId SnapshotId) error {

	reader, err := snapshotId.Read()
	if err != nil {
		return err
	}
	defer reader.Close()

	for {
		writer, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return reader.Close()
			}
			return err
		}
		if writer == nil {
			return errors.New("gobdb: decoded nil Writer on ApplySnapshot")
		}
		if _, err := writer.Write(root); err != nil {
			return err
		}
	}
}
