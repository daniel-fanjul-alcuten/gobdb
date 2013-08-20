package gobdb

import (
	"fmt"
	"io"
	"testing"
)

func ExampleMemBurstRepository() {

	// the testRoot object keeps a counter
	bursts := NewMemBurstRepository()
	dispatcher := NewDefaultBurstDispatcher(bursts)
	database := NewDefaultDatabase(&testRoot{0}, 0, dispatcher)

	// the testWriter increments the counter
	result1, _ := database.Write(&testWriter{3})
	fmt.Println("first write:", result1)

	// the testWriter decrements the counter
	result2, _ := database.Write(&testWriter{-1})
	fmt.Println("second write:", result2)

	// the testReader reads the counter
	result3 := database.Read(&testReader{})
	fmt.Println("read:", result3)

	_ = dispatcher.Close()
	fmt.Println("bursts:", len(bursts.Bursts()))
	// Output: first write: 3
	// second write: 2
	// read: 2
	// bursts: 1
}

func TestMemBurstRepositoryInterface(t *testing.T) {

	var i interface{} = NewMemBurstRepository()
	if _, ok := i.(BurstRepository); !ok {
		t.Error(i)
	}
	if _, ok := i.(WriteBurstRepository); !ok {
		t.Error(i)
	}
}

func TestMemBurstRepositoryEmpty(t *testing.T) {

	repository := NewMemBurstRepository()
	if repository == nil {
		t.Fatal(repository)
	}

	bursts := repository.Bursts()
	if len(bursts) != 0 {
		t.Error(bursts)
	}
}

func TestMemBurstRepositoryBursts(t *testing.T) {

	repository := NewMemBurstRepository()
	if repository == nil {
		t.Fatal(repository)
	}

	wburst, err := repository.WriteBurst()
	if wburst == nil {
		t.Fatal(wburst)
	}
	if err != nil {
		t.Error(err)
	}
	if wburst.First() != 0 {
		t.Error(wburst.First())
	}
	if wburst.Last() != 0 {
		t.Error(wburst.Last())
	}
	defer wburst.Close()

	if err := wburst.Write(Transaction{1, &testWriter{11}}); err != nil {
		t.Error(err)
	}
	if wburst.First() != 1 {
		t.Error(wburst.First())
	}
	if wburst.Last() != 1 {
		t.Error(wburst.Last())
	}

	if err := wburst.Write(Transaction{2, &testWriter{12}}); err != nil {
		t.Error(err)
	}
	if wburst.First() != 1 {
		t.Error(wburst.First())
	}
	if wburst.Last() != 2 {
		t.Error(wburst.Last())
	}

	if err := wburst.Close(); err != nil {
		t.Error(err)
	}

	var id BurstId
	if bursts := repository.Bursts(); len(bursts) != 1 || bursts[0].First() != 1 || bursts[0].Last() != 2 {
		t.Error(bursts)
	} else {
		id = bursts[0]
		if bursts[0].Repository() != repository {
			t.Error(bursts[0].Repository())
		}
	}

	rburst, err := repository.ReadBurst(id)
	if rburst == nil {
		t.Fatal(rburst)
	}
	if err != nil {
		t.Error(err)
	}
	if rburst.Id() != id {
		t.Error(rburst.Id)
	}
	defer rburst.Close()

	transaction, err := rburst.Read()
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

	transaction, err = rburst.Read()
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

	if _, err := rburst.Read(); err != io.EOF {
		t.Error(err)
	}
}
