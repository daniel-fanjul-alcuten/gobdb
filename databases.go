package gobdb

// A controller of one instance of a Root object.
// On startup, it may read a Snapshot and apply the Writers to the Root object.
// Then, it may read some Bursts and apply the Transactions.
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

// The user defined function to take a Snapshot of a Root.
// It invokes the given function as many times as needed with the sequence of
// Writers that are enough to recover the same state of the Root.
// It the given function returns an errors, it must be returned immediately.
type Snapshooter func(Root, func(...Writer) error) error

// A database that can take snapshots of the SnapshotRoot object.
type SnapshotDatabase interface {
	Database
	// It invokes Root.Snapshot() and writes all its Writers into the
	// WriteSnapshotRepository.
	TakeSnapshot(Snapshooter, WriteSnapshotRepository) error
}
