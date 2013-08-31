package gobdb

import (
	"testing"
)

func TestSortBursts(t *testing.T) {

	id1 := &memBurstId{1, 2, nil}
	id2 := &memBurstId{1, 3, nil}
	id3 := &memBurstId{2, 3, nil}
	id4 := &memBurstId{2, 4, nil}
	ids := []BurstId{id1, id2, id3, id4}

	SortBursts(ids)

	if ids[0] != id2 {
		t.Error(ids[0])
	}
	if ids[1] != id1 {
		t.Error(ids[1])
	}
	if ids[2] != id4 {
		t.Error(ids[2])
	}
	if ids[3] != id3 {
		t.Error(ids[3])
	}
}
