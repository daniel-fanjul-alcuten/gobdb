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

// It defines the order in which the Writers must be reapplied.
type TransactionId uint64

// The Writers must be reapplied in order.
type Transaction struct {
	Id TransactionId
	Writer
}
