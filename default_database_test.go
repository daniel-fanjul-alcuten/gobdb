package gobdb

import (
	"fmt"
	"io"
	"testing"
)

func ExampleDefaultDatabase() {

	bursts := NewMemBurstRepository()
	dispatcher := NewDefaultBurstDispatcher(bursts)
	snapshots := NewMemSnapshotRepository()
	{
		// the testRoot object keeps a counter
		database := NewDefaultDatabase(&testRoot{0}, 0, dispatcher)

		// the testWriter increments the counter
		result1, _ := database.Write(&testWriter{3})
		fmt.Println("first write:", result1)

		// testSnapshooter is a Snapshooter for the type testRoot
		_ = database.TakeSnapshot(testSnapshooter, snapshots)

		// the testWriter decrements the counter
		result2, _ := database.Write(&testWriter{-1})
		fmt.Println("second write:", result2)

		_ = dispatcher.Close()
	}

	snapshotIds, _ := snapshots.Snapshots()
	snapshotId := snapshotIds[0]
	{
		root := &testRoot{0}
		_ = ApplySnapshot(root, snapshotId)
		database := NewDefaultDatabase(root, snapshotId.Id(), nil)

		// the testReader reads the counter
		result := database.Read(&testReader{})
		fmt.Println("after snapshot:", result)
	}
	// Output: first write: 3
	// second write: 2
	// after snapshot: 3
}

func TestDefaultDatabaseInterface(t *testing.T) {

	var i interface{}
	i = NewDefaultDatabase(nil, 0, nil)
	if _, ok := i.(Database); !ok {
		t.Error(i)
	}
	if _, ok := i.(WriteDatabase); !ok {
		t.Error(i)
	}
	if _, ok := i.(SnapshotDatabase); !ok {
		t.Error(i)
	}
}

func TestDefaultDatabaseEmpty(t *testing.T) {

	database := NewDefaultDatabase(&testRoot{}, 0, nil)
	if database == nil {
		t.Fatal(database)
	}

	result := database.Read(&testReader{})
	if value, ok := result.(int); !ok {
		t.Error(result)
	} else if value != 0 {
		t.Error(value)
	}
}

func TestDefaultDatabaseWrite(t *testing.T) {

	database := NewDefaultDatabase(&testRoot{}, 0, nil)
	if database == nil {
		t.Fatal(database)
	}

	result, err := database.Write(&testWriter{3})
	if value, ok := result.(int); !ok {
		t.Error(result)
	} else if value != 3 {
		t.Error(value)
	}
	if err != nil {
		t.Error(err)
	}

	result, err = database.Write(&testWriter{-1})
	if value, ok := result.(int); !ok {
		t.Error(result)
	} else if value != 2 {
		t.Error(value)
	}
	if err != nil {
		t.Error(err)
	}

	result = database.Read(&testReader{})
	if value, ok := result.(int); !ok {
		t.Error(result)
	} else if value != 2 {
		t.Error(value)
	}
}

func TestDefaultDatabaseWithBurstDispather(t *testing.T) {

	repository := NewMemBurstRepository()
	dispatcher := NewDefaultBurstDispatcher(repository)
	defer dispatcher.Close()
	database := NewDefaultDatabase(&testRoot{}, 0, dispatcher)
	if database == nil {
		t.Fatal(database)
	}

	result, err := database.Write(&testWriter{11})
	if value, ok := result.(int); !ok {
		t.Error(result)
	} else if value != 11 {
		t.Error(value)
	}
	if err != nil {
		t.Error(err)
	}

	result, err = database.Write(&testWriter{12})
	if value, ok := result.(int); !ok {
		t.Error(result)
	} else if value != 23 {
		t.Error(value)
	}
	if err != nil {
		t.Error(err)
	}

	result = database.Read(&testReader{})
	if value, ok := result.(int); !ok {
		t.Error(result)
	} else if value != 23 {
		t.Error(value)
	}

	if err := dispatcher.Close(); err != nil {
		t.Error(err)
	}

	bursts, err := repository.Bursts()
	if len(bursts) != 1 {
		t.Fatal(len(bursts))
	}
	if err != nil {
		t.Error(err)
	}

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

	if _, err := burst.Read(); err != io.EOF {
		t.Error(err)
	}
}

func TestDefaultDatabaseTakeSnapshot(t *testing.T) {

	database := NewDefaultDatabase(&testRoot{}, 0, nil)
	if database == nil {
		t.Fatal(database)
	}

	result, err := database.Write(&testWriter{11})
	if value, ok := result.(int); !ok {
		t.Error(result)
	} else if value != 11 {
		t.Error(value)
	}
	if err != nil {
		t.Error(err)
	}

	result, err = database.Write(&testWriter{12})
	if value, ok := result.(int); !ok {
		t.Error(result)
	} else if value != 23 {
		t.Error(value)
	}
	if err != nil {
		t.Error(err)
	}

	repository := NewMemSnapshotRepository()
	if err := database.TakeSnapshot(testSnapshooter, repository); err != nil {
		t.Error(err)
	}

	var id SnapshotId
	snapshots, err := repository.Snapshots()
	if err != nil {
		t.Error(err)
	}
	if len(snapshots) != 1 {
		t.Error(snapshots)
	} else {
		id = snapshots[0]
		if id.Id() != 2 {
			t.Error(id.Id())
		}
	}

	snapshot, err := repository.ReadSnapshot(id)
	if err != nil {
		t.Error(err)
	}
	defer snapshot.Close()

	writer, err := snapshot.Read()
	if err != nil {
		t.Error(err)
	}
	tw, ok := writer.(*testWriter)
	if !ok {
		t.Fatalf("%#v", writer)
	}
	if tw.Increment != 23 {
		t.Error(tw.Increment)
	}

	if _, err := snapshot.Read(); err != io.EOF {
		t.Error(err)
	}
}
