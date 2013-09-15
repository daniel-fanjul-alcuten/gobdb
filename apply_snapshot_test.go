package gobdb

import (
	"testing"
)

func TestApplySnapshot(t *testing.T) {

	repository := NewMemSnapshotRepository()
	if repository == nil {
		t.Fatal(repository)
	}

	wsnapshot, err := repository.WriteSnapshot(1)
	if wsnapshot == nil {
		t.Fatal(wsnapshot)
	}
	if err != nil {
		t.Error(err)
	}
	defer wsnapshot.Close()

	if err := wsnapshot.Write(&testWriter{11}); err != nil {
		t.Error(err)
	}

	if err := wsnapshot.Close(); err != nil {
		t.Error(err)
	}

	root := &testRoot{}
	snapshots, err := repository.Snapshots()
	if len(snapshots) != 1 {
		t.Fatal(len(snapshots))
	}
	if err != nil {
		t.Error(err)
	}
	snapshotId := snapshots[0]
	if err := ApplySnapshot(root, snapshotId); err != nil {
		t.Error(err)
	}
	if root.counter != 11 {
		t.Error(root.counter)
	}
}
