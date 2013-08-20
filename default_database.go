package gobdb

// A Database that does not use any Repository.
type DefaultDatabase struct {
	root       Root
	lastId     TransactionId
	dispatcher BurstDispatcher
}

// New instance. The SnapshotRepository and the BurstDispatcher are optional.
func NewDefaultDatabase(root Root, snapshots SnapshotRepository, dispatcher BurstDispatcher) (db *DefaultDatabase, err error) {
	var id TransactionId
	if snapshots != nil {
		id, err = ApplySnapshot(root, snapshots)
		if err != nil {
			return
		}
	}
	db = &DefaultDatabase{root, id, dispatcher}
	return
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
	if db.dispatcher != nil {
		db.lastId++
		err = db.dispatcher.Write(Transaction{db.lastId, writer})
	}
	return
}
