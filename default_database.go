package gobdb

// A Database, WriteDatabase and SnapshotDatabase.
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
func (db *DefaultDatabase) Write(writer Writer) (value interface{}, err error) {
	value, err = writer.Write(db.root)
	if err != nil {
		return
	}
	db.lastId++
	if db.dispatcher != nil {
		err = db.dispatcher.Write(Transaction{db.lastId, writer})
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

	return snapshooter(db.root, func(writers ...Writer) error {
		for _, w := range writers {
			if err := writer.Write(w); err != nil {
				return err
			}
		}
		return nil
	})
}
