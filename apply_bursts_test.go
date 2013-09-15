package gobdb

import (
	"testing"
)

func TestApplyBursts(t *testing.T) {

	repository := NewMemBurstRepository()

	wburst, err := repository.WriteBurst()
	if err != nil {
		t.Error(err)
	}
	if err := wburst.Write(Transaction{1, &testWriter{11}}); err != nil {
		t.Error(err)
	}
	if err := wburst.Write(Transaction{3, &testWriter{13}}); err != nil {
		t.Error(err)
	}
	if err := wburst.Close(); err != nil {
		t.Error(err)
	}

	wburst, err = repository.WriteBurst()
	if err != nil {
		t.Error(err)
	}
	if err := wburst.Write(Transaction{2, &testWriter{12}}); err != nil {
		t.Error(err)
	}
	if err := wburst.Write(Transaction{3, &testWriter{13}}); err != nil {
		t.Error(err)
	}
	if err := wburst.Write(Transaction{4, &testWriter{14}}); err != nil {
		t.Error(err)
	}
	if err := wburst.Write(Transaction{6, &testWriter{16}}); err != nil {
		t.Error(err)
	}
	if err := wburst.Close(); err != nil {
		t.Error(err)
	}

	wburst, err = repository.WriteBurst()
	if err != nil {
		t.Error(err)
	}
	if err := wburst.Write(Transaction{6, &testWriter{16}}); err != nil {
		t.Error(err)
	}
	if err := wburst.Close(); err != nil {
		t.Error(err)
	}

	root := &testRoot{0}
	bursts, err := repository.Bursts()
	if err != nil {
		t.Error(err)
	}
	var id TransactionId
	if err := ApplyBursts(root, 0, &id, bursts); err != nil {
		t.Error(err)
	}
	if id != 4 {
		t.Error(id)
	}
	if root.counter != 50 {
		t.Error(root.counter)
	}
}
