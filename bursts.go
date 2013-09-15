package gobdb

// A Burst is a gob stream of Transactions.
// It is required to contain them in order, but not to be consecutive.
type BurstId interface {
	First() TransactionId
	Last() TransactionId
	Repository() BurstRepository
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
	Bursts() ([]BurstId, error)
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
