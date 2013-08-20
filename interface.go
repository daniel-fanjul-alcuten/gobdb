// In-memory non-relational transaction-based snapshot-based database
package gobdb

import (
	"io"
)

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

// A Burst is a gob stream of Transactions.
// It is required to contain them in order, but not to be consecutive.
type BurstId struct {
	// The first TransactionId of the stream.
	First TransactionId
	// The last TransactionId of the stream.
	Last TransactionId
}

// An io.ReadCloser of a Burst.
type BurstReader struct {
	Id BurstId
	io.ReadCloser
}

// An io.WriteCloser of a Burst.
type BurstWriter struct {
	First TransactionId
	io.WriteCloser
}

// A Snapshot is a gob stream of Writers. To apply these Writers in sequence
// must yield the same result than to apply the Writers of all Transactions
// starting from the first one until the give one.
type SnapshotId TransactionId

// An io.ReadCloser of a Snapshot.
type SnapshotReader struct {
	Id SnapshotId
	io.ReadCloser
}

// An io.WriteCloser of a Snapshot.
type SnapshotWriter struct {
	Id SnapshotId
	io.WriteCloser
}

// A container of Bursts and Snapshots.
type Repository interface {
	// List of all Bursts.
	Bursts() []BurstId
	// Get an io.ReadCloser of a Burst.
	ReadBurst(BurstId) (BurstReader, error)
	// Get an io.WriteCloser of a Burst.
	WriteBurst(TransactionId) (BurstWriter, error)
	// List of all Snapshots.
	Snapshots() []SnapshotId
	// Get an io.ReadCloser of a Snapshot.
	ReadSnapshot(SnapshotId) (SnapshotReader, error)
	// Get an io.WriteCloser of a Snapshot.
	WriteSnapshot(TransactionId) (SnapshotWriter, error)
}

// A controller of a Root object that keep it synced with a Repository.
type Database interface {
	Read(Reader) interface{}
	Write(Writer) interface{}
}
