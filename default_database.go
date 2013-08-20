package gobdb

// A Database that does not use any Repository.
type DefaultDatabase struct {
	root       Root
	lastId     TransactionId
	dispatcher BurstDispatcher
}

// New instance. The SnapshotId and the BurstDispatcher are optional.
func NewDefaultDatabase(root Root, snapshotId SnapshotId, dispatcher BurstDispatcher) (db *DefaultDatabase, err error) {
	var id TransactionId
	if snapshotId != nil {
		err = ApplySnapshot(root, snapshotId)
		if err != nil {
			return
		}
		id = snapshotId.Id()
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
