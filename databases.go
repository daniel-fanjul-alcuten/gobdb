package gobdb

// A controller of one instance of a Root object.
// On startup, it may read a Snapshot and apply the Writers to the Root object.
// Then, it may read some Bursts and apply the Transactions.
type Database interface {
	// It applies the Reader to the Root object and returns its result.
	Read(Reader) interface{}
}

// It dispatches the writes of the Transactions to one or more Bursts.
// Bursts can be created and closed when needed,
// so rotation of files can be implemented.
type BurstDispatcher interface {
	// It writes the Transaction in a previous or new Burst.
	Write(Transaction) error
	// It forces the rotation.
	Rotate() error
	// It releases resources.
	Close() error
}

// A database that can update the Root object.
// On every Write(), it updates the Root object and then writes the Writer into
// a BurstDispatcher.
type WriteDatabase interface {
	Database
	// It applies the Writer to the Root object and returns its result and its
	// error as the first one.
	// If succcessful, it writes the Writer to the BurstDispatcher and returns
	// its error as the second one.
	Write(Writer) (interface{}, error, error)
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
