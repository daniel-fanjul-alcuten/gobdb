package gobdb

// A Burst is a gob stream of Transactions.
// It is required to contain them in order, but not to be consecutive.
type BurstId interface {
	// The first transaction in the Burst.
	First() TransactionId
	// The last transaction in the Burst.
	Last() TransactionId
	// The Repository it comes from.
	Repository() BurstRepository
	// Get a BurstReader of a Burst.
	Read() (BurstReader, error)
}

// It reads the Transactions of a Burst.
type BurstReader interface {
	// The id of the Burst.
	Id() BurstId
	// It reads a Transaction until io.EOF.
	Read() (Transaction, error)
	// It closes the stream.
	Close() error
}

// A container that can read Bursts.
type BurstRepository interface {
	// List of all Bursts.
	Bursts() ([]BurstId, error)
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
