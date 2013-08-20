package gobdb

// A Snapshot is a gob stream of Writers. To apply these Writers in sequence
// must yield the same result than to apply the Writers of all Transactions
// starting from the first one until the given one.
type SnapshotId interface {
	Id() TransactionId
	Repository() SnapshotRepository
}

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
	Id() TransactionId
	Write(Writer) error
	Close() error
}

// A container that can write Snapshots.
type WriteSnapshotRepository interface {
	// Get a SnapshotWriter of a Snapshot.
	WriteSnapshot(TransactionId) (SnapshotWriter, error)
}
