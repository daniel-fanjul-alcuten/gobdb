package gobdb

// The type of objects that will be held in memory.
type Root interface{}

// An operation that will read some data from a Root.
// It must be deterministic.
// It must be gob encodable.
type Reader interface {
	Read(Root) interface{}
}

// An operation that will update and read some data from a Root.
// It must be deterministic.
// It must be gob encodable.
type Writer interface {
	Write(Root) interface{}
}
