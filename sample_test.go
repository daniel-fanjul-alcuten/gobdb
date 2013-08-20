package gobdb

import (
	"encoding/gob"
)

type testRoot struct {
	counter int
}

type testReader struct {
}

func (op *testReader) Read(root Root) interface{} {
	r := root.(*testRoot)
	return r.counter
}

type testWriter struct {
	Increment int
}

func (op *testWriter) Write(root Root) (interface{}, error) {
	r := root.(*testRoot)
	r.counter += op.Increment
	return r.counter, nil
}

func testSnapshooter(root Root, write func(...Writer) error) error {
	r := root.(*testRoot)
	return write(&testWriter{r.counter})
}

func init() {
	gob.Register(&testWriter{})
}
