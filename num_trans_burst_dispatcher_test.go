package gobdb

import (
	"io"
	"testing"
)

func TestNumTransactionsBurstDispatcherInterface(t *testing.T) {

	var i interface{} = NewNumTransactionsBurstDispatcher(0, nil)
	if _, ok := i.(BurstDispatcher); !ok {
		t.Error(i)
	}
}

func TestNumTransactionsBurstDispatcherWrite(t *testing.T) {

	repository := NewMemBurstRepository()
	dispatcher := NewNumTransactionsBurstDispatcher(2, NewDefaultBurstDispatcher(repository))
	if dispatcher == nil {
		t.Fatal(dispatcher)
	}
	defer dispatcher.Close()

	if err := dispatcher.Write(Transaction{1, &testWriter{11}}); err != nil {
		t.Error(err)
	}
	bursts, err := repository.Bursts()
	if len(bursts) != 0 {
		t.Error(len(bursts))
	}
	if err != nil {
		t.Error(err)
	}

	if err := dispatcher.Rotate(); err != nil {
		t.Error(err)
	}
	bursts, err = repository.Bursts()
	if len(bursts) != 1 {
		t.Error(len(bursts))
	}
	if err != nil {
		t.Error(err)
	}

	if err := dispatcher.Write(Transaction{2, &testWriter{12}}); err != nil {
		t.Error(err)
	}
	bursts, err = repository.Bursts()
	if len(bursts) != 1 {
		t.Error(len(bursts))
	}
	if err != nil {
		t.Error(err)
	}

	if err := dispatcher.Write(Transaction{3, &testWriter{13}}); err != nil {
		t.Error(err)
	}
	bursts, err = repository.Bursts()
	if len(bursts) != 2 {
		t.Error(len(bursts))
	}
	if err != nil {
		t.Error(err)
	}

	if err := dispatcher.Write(Transaction{4, &testWriter{14}}); err != nil {
		t.Error(err)
	}
	bursts, err = repository.Bursts()
	if len(bursts) != 2 {
		t.Error(len(bursts))
	}
	if err != nil {
		t.Error(err)
	}

	if err := dispatcher.Close(); err != nil {
		t.Error(err)
	}
	bursts, err = repository.Bursts()
	if len(bursts) != 3 {
		t.Fatal(len(bursts))
	}
	if err != nil {
		t.Error(err)
	}
	SortBursts(bursts)

	burst, err := repository.ReadBurst(bursts[0])
	if err != nil {
		t.Error(err)
	}
	defer burst.Close()

	transaction, err := burst.Read()
	if err != nil {
		t.Error(err)
	}
	if transaction.Id != 1 {
		t.Error(transaction.Id)
	}
	tw, ok := transaction.Writer.(*testWriter)
	if !ok {
		t.Fatalf("%#v", transaction.Writer)
	}
	if tw.Increment != 11 {
		t.Error(tw.Increment)
	}

	if _, err := burst.Read(); err != io.EOF {
		t.Error(err)
	}

	burst, err = repository.ReadBurst(bursts[1])
	if err != nil {
		t.Error(err)
	}
	defer burst.Close()

	transaction, err = burst.Read()
	if err != nil {
		t.Error(err)
	}
	if transaction.Id != 2 {
		t.Error(transaction.Id)
	}
	tw, ok = transaction.Writer.(*testWriter)
	if !ok {
		t.Fatalf("%#v", transaction.Writer)
	}
	if tw.Increment != 12 {
		t.Error(tw.Increment)
	}

	transaction, err = burst.Read()
	if err != nil {
		t.Error(err)
	}
	if transaction.Id != 3 {
		t.Error(transaction.Id)
	}
	tw, ok = transaction.Writer.(*testWriter)
	if !ok {
		t.Fatalf("%#v", transaction.Writer)
	}
	if tw.Increment != 13 {
		t.Error(tw.Increment)
	}

	if _, err := burst.Read(); err != io.EOF {
		t.Error(err)
	}

	burst, err = repository.ReadBurst(bursts[2])
	if err != nil {
		t.Error(err)
	}
	defer burst.Close()

	transaction, err = burst.Read()
	if err != nil {
		t.Error(err)
	}
	if transaction.Id != 4 {
		t.Error(transaction.Id)
	}
	tw, ok = transaction.Writer.(*testWriter)
	if !ok {
		t.Fatalf("%#v", transaction.Writer)
	}
	if tw.Increment != 14 {
		t.Error(tw.Increment)
	}

	if _, err := burst.Read(); err != io.EOF {
		t.Error(err)
	}
}
