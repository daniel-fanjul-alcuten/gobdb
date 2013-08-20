package gobdb

// A Database that does not use any Repository.
type DefaultDatabase struct {
	root Root
}

// New instance.
func NewDefaultDatabase(root Root) *DefaultDatabase {
	return &DefaultDatabase{root}
}

// Implements Database.Read().
func (db *DefaultDatabase) Read(reader Reader) interface{} {
	return reader.Read(db.root)
}

// Implements WriteDatabase.Write().
func (db *DefaultDatabase) Write(writer Writer) (interface{}, error) {
	return writer.Write(db.root), nil
}
