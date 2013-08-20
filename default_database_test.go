package gobdb

import (
	"io"
	"testing"
)

func TestDefaultDatabaseInterface(t *testing.T) {

	var i interface{}
	i, _ = NewDefaultDatabase(nil, nil, nil)
	if _, ok := i.(Database); !ok {
		t.Error(i)
	}
	if _, ok := i.(WriteDatabase); !ok {
		t.Error(i)
	}
}

func TestDefaultDatabaseEmpty(t *testing.T) {

	database, err := NewDefaultDatabase(&testRoot{}, nil, nil)
	if database == nil {
		t.Fatal(database)
	}
	if err != nil {
		t.Error(err)
	}

	result := database.Read(&testReader{})
	if value, ok := result.(int); !ok {
		t.Error(result)
	} else if value != 0 {
		t.Error(value)
	}
}

func TestDefaultDatabaseWrite(t *testing.T) {

	database, err := NewDefaultDatabase(&testRoot{}, nil, nil)
	if database == nil {
		t.Fatal(database)
	}
	if err != nil {
		t.Error(err)
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
	database, err := NewDefaultDatabase(&testRoot{}, nil, dispatcher)
	if database == nil {
		t.Fatal(database)
	}
	if err != nil {
		t.Error(err)
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

	bursts := repository.Bursts()
	if len(bursts) != 1 {
		t.Fatal(bursts)
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

func TestDefaultDatabaseWithSnapshotReader(t *testing.T) {

	snapshots := NewMemSnapshotRepository()

	wsnapshot, err := snapshots.WriteSnapshot(1)
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

	snapshotId := snapshots.Snapshots()[0]
	bursts := NewMemBurstRepository()
	dispatcher := NewDefaultBurstDispatcher(bursts)
	defer dispatcher.Close()
	database, err := NewDefaultDatabase(&testRoot{}, snapshotId, dispatcher)
	if database == nil {
		t.Fatal(database)
	}
	if err != nil {
		t.Error(err)
	}

	result, err := database.Write(&testWriter{12})
	if value, ok := result.(int); !ok {
		t.Error(result)
	} else if value != 23 {
		t.Error(value)
	}
	if err != nil {
		t.Error(err)
	}

	if err := dispatcher.Close(); err != nil {
		t.Error(err)
	}

	rbursts := bursts.Bursts()
	if len(rbursts) != 1 {
		t.Fatal(rbursts)
	}
	burst, err := bursts.ReadBurst(rbursts[0])
	if err != nil {
		t.Error(err)
	}
	defer burst.Close()

	transaction, err := burst.Read()
	if err != nil {
		t.Error(err)
	}
	if transaction.Id != 2 {
		t.Error(transaction.Id)
	}
	tw, ok := transaction.Writer.(*testWriter)
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
