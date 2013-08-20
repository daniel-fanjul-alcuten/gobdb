package gobdb

// The object that will be kept in memory.
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
	Write(Root) (interface{}, error)
}

// It defines the order in which the Writers must be reapplied.
// The sequence starts from 1 and the zero value is considered like a nil.
type TransactionId uint64

// The Writers must be reapplied in order.
type Transaction struct {
	Id TransactionId
	Writer
}
