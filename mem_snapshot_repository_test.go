package gobdb

import (
	"io"
	"testing"
)

func TestMemSnapshotRepositoryInterface(t *testing.T) {

	var i interface{} = NewMemSnapshotRepository()
	if _, ok := i.(SnapshotRepository); !ok {
		t.Error(i)
	}
	if _, ok := i.(WriteSnapshotRepository); !ok {
		t.Error(i)
	}
}

func TestMemSnapshotRepositoryEmpty(t *testing.T) {

	repository := NewMemSnapshotRepository()
	if repository == nil {
		t.Fatal(repository)
	}

	snapshots := repository.Snapshots()
	if len(snapshots) != 0 {
		t.Error(snapshots)
	}
}

func TestMemSnapshotRepositorySnapshots(t *testing.T) {

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
	if wsnapshot.Id() != 1 {
		t.Error(wsnapshot.Id())
	}
	defer wsnapshot.Close()

	if err := wsnapshot.Write(&testWriter{11}); err != nil {
		t.Error(err)
	}

	if err := wsnapshot.Close(); err != nil {
		t.Error(err)
	}

	var id SnapshotId
	if snapshots := repository.Snapshots(); len(snapshots) != 1 || snapshots[0].Id() != 1 {
		t.Error(snapshots)
	} else {
		id = snapshots[0]
	}

	rsnapshot, err := repository.ReadSnapshot(id)
	if rsnapshot == nil {
		t.Fatal(rsnapshot)
	}
	if err != nil {
		t.Error(err)
	}
	if rsnapshot.Id() != id {
		t.Error(rsnapshot.Id)
	}
	defer rsnapshot.Close()

	writer, err := rsnapshot.Read()
	if err != nil {
		t.Error(err)
	}
	tw, ok := writer.(*testWriter)
	if !ok {
		t.Fatalf("%#v", writer)
	}
	if tw.Increment != 11 {
		t.Error(tw.Increment)
	}

	if _, err := rsnapshot.Read(); err != io.EOF {
		t.Error(err)
	}
}
