package gobdb

import (
	"encoding/gob"
)

type testRoot struct {
	counter int
}

type testReader struct {
}

func (r *testReader) Read(Root) interface{} {
	return nil
}

type testWriter struct {
	Increment int
}

func (w *testWriter) Write(Root) interface{} {
	return nil
}

func init() {
	gob.Register(&testWriter{})
}
