package gobdb

import (
	"testing"
)

func TestSortSnapshot(t *testing.T) {

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

	wsnapshot, err = repository.WriteSnapshot(2)
	if wsnapshot == nil {
		t.Fatal(wsnapshot)
	}
	if err != nil {
		t.Error(err)
	}
	defer wsnapshot.Close()

	if err := wsnapshot.Write(&testWriter{12}); err != nil {
		t.Error(err)
	}

	if err := wsnapshot.Close(); err != nil {
		t.Error(err)
	}

	snapshots, err := repository.Snapshots()
	if len(snapshots) != 2 {
		t.Fatal(len(snapshots))
	}
	if err != nil {
		t.Error(err)
	}
	SortSnapshots(snapshots)
	if snapshots[0].Id() != 2 {
		t.Error(snapshots[0].Id())
	}
	if snapshots[1].Id() != 1 {
		t.Error(snapshots[1].Id())
	}
}
