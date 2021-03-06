package gobdb

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func ExampleDirBurstRepository() {

	dir, _ := ioutil.TempDir("", "gobdb-")
	defer os.RemoveAll(dir)

	// the testRoot object keeps a counter
	bursts := NewDirBurstRepository(dir)
	{
		dispatcher := NewDefaultBurstDispatcher(bursts)
		database := NewDefaultDatabase(&testRoot{0}, 0, dispatcher)

		// the testWriter increments the counter
		_, _, _ = database.Write(&testWriter{3})
		// the testWriter decrements the counter
		_, _, _ = database.Write(&testWriter{-1})
		// the testReader reads the counter
		result := database.Read(&testReader{})
		fmt.Println("before burst:", result)

		_ = dispatcher.Close()
	}

	{
		root := &testRoot{0}
		burstIds, _ := bursts.Bursts()
		var id TransactionId
		_ = ApplyBursts(root, 0, &id, burstIds)
		database := NewDefaultDatabase(root, id, nil)

		// the testReader reads the counter
		result := database.Read(&testReader{})
		fmt.Println("after burst:", result)
	}
	// Output: before burst: 2
	// after burst: 2
}

func TestDirBurstRepositoryInterface(t *testing.T) {

	var i interface{} = NewDirBurstRepository("")
	if _, ok := i.(BurstRepository); !ok {
		t.Error(i)
	}
	if _, ok := i.(WriteBurstRepository); !ok {
		t.Error(i)
	}
}

func TestDirBurstRepositoryEmpty(t *testing.T) {

	dir, err := ioutil.TempDir("", "gobdb-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	repository := NewDirBurstRepository(dir)
	if repository == nil {
		t.Fatal(repository)
	}

	bursts, err := repository.Bursts()
	if len(bursts) != 0 {
		t.Error(len(bursts))
	}
	if err != nil {
		t.Error(err)
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
	if err := wburst.Close(); err != nil {
		t.Error(err)
	}

	bursts, err = repository.Bursts()
	if len(bursts) != 0 {
		t.Error(len(bursts))
	}
	if err != nil {
		t.Error(err)
	}
}

func TestDirBurstRepositoryBursts(t *testing.T) {

	dir, err := ioutil.TempDir("", "gobdb-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	repository := NewDirBurstRepository(dir)
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
	bursts, err := repository.Bursts()
	if err != nil {
		t.Error(err)
	}
	if len(bursts) != 1 || bursts[0].First() != 1 || bursts[0].Last() != 2 {
		t.Fatal(bursts)
	} else {
		id = bursts[0]
		if bursts[0].Repository() != repository {
			t.Error(bursts[0].Repository())
		}
	}

	rburst, err := id.Read()
	if err != nil {
		t.Error(err)
	}
	if rburst == nil {
		t.Fatal(rburst)
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
