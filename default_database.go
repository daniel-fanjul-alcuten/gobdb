package gobdb

// A Database that does not use any Repository.
type DefaultDatabase struct {
	root       Root
	lastId     TransactionId
	dispatcher BurstDispatcher
}

// New instance. The BurstDispatcher is optional.
func NewDefaultDatabase(root Root, dispatcher BurstDispatcher) *DefaultDatabase {
	return &DefaultDatabase{root, 0, dispatcher}
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
