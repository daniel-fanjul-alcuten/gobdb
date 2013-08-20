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

	root := &testRoot{}
	id, err := ApplySnapshot(root, repository)
	if root.counter != 12 {
		t.Error(root.counter)
	}
	if id != 2 {
		t.Error(id)
	}
	if err != nil {
		t.Error(err)
	}
}
