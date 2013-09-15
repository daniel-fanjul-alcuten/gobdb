package gobdb

// A Database, WriteDatabase and SnapshotDatabase. No thread-safe.
type DefaultDatabase struct {
	root       Root
	lastId     TransactionId
	dispatcher BurstDispatcher
}

// New instance. The TransactionId is the last one that has been applied to the
// Root. The BurstDispatcher is optional.
func NewDefaultDatabase(root Root, lastId TransactionId, dispatcher BurstDispatcher) *DefaultDatabase {
	return &DefaultDatabase{root, lastId, dispatcher}
}

// Implements Database.Read().
func (db *DefaultDatabase) Read(reader Reader) interface{} {
	return reader.Read(db.root)
}

// Implements WriteDatabase.Write().
func (db *DefaultDatabase) Write(writer Writer) (value interface{}, err1 error, err2 error) {
	value, err1 = writer.Write(db.root)
	if err1 != nil {
		return
	}
	db.lastId++
	if db.dispatcher != nil {
		err2 = db.dispatcher.Write(Transaction{db.lastId, writer})
	}
	return
}

// Implements SnapshotDatabase.TakeSnapshot().
func (db *DefaultDatabase) TakeSnapshot(snapshooter Snapshooter, repository WriteSnapshotRepository) error {

	writer, err := repository.WriteSnapshot(db.lastId)
	if err != nil {
		return err
	}
	defer writer.Close()

	write := func(writers ...Writer) error {
		for _, w := range writers {
			if err := writer.Write(w); err != nil {
				return err
			}
		}
		return nil
	}

	if err := snapshooter(db.root, write); err != nil {
		return err
	}

	return writer.Close()
}
