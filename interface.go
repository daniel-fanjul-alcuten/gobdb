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
	Write(Root) interface{}
}

// It defines the order in which the Writers must be reapplied.
// The sequence starts from 1 and the zero value is considered like a nil.
type TransactionId uint64

// The Writers must be reapplied in order.
type Transaction struct {
	Id TransactionId
	Writer
}

// A Burst is a gob stream of Transactions.
// It is required to contain them in order, but not to be consecutive.
type BurstId interface {
	First() TransactionId
	Last() TransactionId
}

// It reads the Transactions of a Burst.
type BurstReader interface {
	Id() BurstId
	Read() (Transaction, error)
	Close() error
}

// A container that can read Bursts.
type BurstRepository interface {
	// List of all Bursts.
	Bursts() []BurstId
	// Get a BurstReader of a Burst.
	ReadBurst(BurstId) (BurstReader, error)
}

// It writes the Transactions of a Burst.
type BurstWriter interface {
	First() TransactionId
	Last() TransactionId
	Write(Transaction) error
	Close() error
}

// A container that can write Bursts.
type WriteBurstRepository interface {
	// Get a BurstWriter of a Burst.
	WriteBurst() (BurstWriter, error)
}

// A Snapshot is a gob stream of Writers. To apply these Writers in sequence
// must yield the same result than to apply the Writers of all Transactions
// starting from the first one until the given one.
type SnapshotId TransactionId

// It reads the Writers of a Snapshot.
type SnapshotReader interface {
	Id() SnapshotId
	Read() (Writer, error)
	Close() error
}

// A container that can read Snapshots.
type SnapshotRepository interface {
	// List of all Snapshots.
	Snapshots() []SnapshotId
	// Get a SnapshotReader of a Snapshot.
	ReadSnapshot(SnapshotId) (SnapshotReader, error)
}

// It writes the Writers of a Snapshot.
type SnapshotWriter interface {
	Id() SnapshotId
	Write(Writer) error
	Close() error
}

// A container that can write Snapshots.
type WriteSnapshotRepository interface {
	// Get a SnapshotWriter of a Snapshot.
	WriteSnapshot(TransactionId) (SnapshotWriter, error)
}

// A controller of one instance of a Root object.
// On startup, it may access a SnapshotRepository, read a Snapshot and apply
// the Writers to the Root object.
// Then, it may access a BurstRepository, read some Bursts and apply the
// Writers to the Root object.
type Database interface {
	// It applies the Reader to the Root object and returns its result.
	Read(Reader) interface{}
}

// It dispatches the writes of Writers to one or more Bursts of a
// WriteBurstRepository. It implements the rotation of Bursts.
type BurstDispatcher interface {
	// It creates a Burst or reuses a previous one and writes the Transaction
	// to it.
	// It may close and create new Bursts when some threshold is reached based
	// on criteria like time, number of bytes, number of Transactions and so on.
	Write(Transaction) error
	// Close the last BurstWriter if needed.
	Close() error
}

// A database that can update the Root object.
// On every Write(), it updates the Root object and then writes the Writer into
// a BurstDispatcher.
type WriteDatabase interface {
	Database
	// It applies the Writer to the Root object and returns its result. It writes
	// the Writer to the BurstDispatcher.
	Write(Writer) (interface{}, error)
}

// A Root that can take snapshots
type SnapshotRoot interface {
	Root
	// It invokes the function as many times as needed with the sequence of
	// Writers that are enough to recover the same state of the Root.
	// It the given function returns an errors, it must be returned as soon as
	// possible.
	Snapshot(func(...Writer) error) error
}

// A database that can take snapshots of the SnapshotRoot object.
type SnapshotDatabase interface {
	Database
	// It invokes Root.Snapshot() and writes all its Writers into the
	// WriteSnapshotRepository.
	Snapshot(WriteSnapshotRepository) error
}
