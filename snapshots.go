package gobdb

// A Snapshot is a gob stream of Writers. To apply these Writers in sequence
// must yield the same result than to apply the Writers of all Transactions
// starting from the first one until the given one.
type SnapshotId interface {
	// The last transaction.
	Id() TransactionId
	// The Repository it comes from.
	Repository() SnapshotRepository
	// Get a SnapshotReader of a Snapshot.
	Read() (SnapshotReader, error)
}

// It reads the Writers of a Snapshot.
type SnapshotReader interface {
	// The id of the Snapshot.
	Id() SnapshotId
	// It reads a Writer until io.EOF.
	Read() (Writer, error)
	// It closes the stream.
	Close() error
}

// A container that can read Snapshots.
type SnapshotRepository interface {
	// List of all Snapshots.
	Snapshots() ([]SnapshotId, error)
}

// It writes the Writers of a Snapshot.
type SnapshotWriter interface {
	Id() TransactionId
	Write(Writer) error
	Close() error
}

// A container that can write Snapshots.
type WriteSnapshotRepository interface {
	// Get a SnapshotWriter of a Snapshot.
	WriteSnapshot(TransactionId) (SnapshotWriter, error)
}
