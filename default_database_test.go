package gobdb

import (
	"testing"
)

func TestDefaultDatabaseInterface(t *testing.T) {

	var i interface{} = NewDefaultDatabase(nil)
	if _, ok := i.(Database); !ok {
		t.Error(i)
	}
	if _, ok := i.(WriteDatabase); !ok {
		t.Error(i)
	}
}

func TestDefaultDatabaseEmpty(t *testing.T) {

	database := NewDefaultDatabase(&testRoot{})
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

func TestDefaultDatabaseSimple(t *testing.T) {

	database := NewDefaultDatabase(&testRoot{})
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
