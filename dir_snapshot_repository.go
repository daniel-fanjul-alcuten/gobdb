package gobdb

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const dirSnapshotRepositoryFileNameFormat = "snapshot-%d.gobdb"
const dirSnapshotRepositoryFileNameScanFormat = dirSnapshotRepositoryFileNameFormat + "\n"

// A SnapshotRepository and WriteSnapshotRepository that uses one file per Snapshot.
// Thread-safe, but SnapshotReaders and SnapshotWriters are not.
type DirSnapshotRepository struct {
	dir string
}

func NewDirSnapshotRepository(dir string) *DirSnapshotRepository {
	return &DirSnapshotRepository{dir}
}

func (r *DirSnapshotRepository) Snapshots() ([]SnapshotId, error) {
	file, err := os.Open(r.dir)
	if err != nil {
		return nil, err
	}
	infos, err1 := file.Readdir(-1)
	err2 := file.Close()
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}
	ids := make([]SnapshotId, 0, len(infos))
	for _, info := range infos {
		name := info.Name()
		var id int
		if n, err := fmt.Sscanf(name, dirSnapshotRepositoryFileNameScanFormat, &id); n == 1 && err == nil {
			ids = append(ids, &dirSnapshotId{TransactionId(id), r})
		}
	}
	return ids, nil
}

func (r *DirSnapshotRepository) WriteSnapshot(id TransactionId) (SnapshotWriter, error) {
	file, err := ioutil.TempFile(r.dir, "tmp-snapshot-")
	if err != nil {
		return nil, err
	}
	writer := bufio.NewWriter(file)
	encoder := gob.NewEncoder(writer)
	return &dirSnapshotWriter{file, writer, encoder, id, r}, nil
}

type dirSnapshotId struct {
	id         TransactionId
	repository *DirSnapshotRepository
}

func (id *dirSnapshotId) Id() TransactionId {
	return id.id
}

func (id *dirSnapshotId) Repository() SnapshotRepository {
	return id.repository
}

func (id *dirSnapshotId) Read() (SnapshotReader, error) {
	name := fmt.Sprintf(dirSnapshotRepositoryFileNameFormat, id.id)
	file, err := os.Open(filepath.Join(id.repository.dir, name))
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	decoder := gob.NewDecoder(reader)
	return &dirSnapshotReader{file, decoder, id}, nil
}

type dirSnapshotReader struct {
	file    *os.File
	decoder *gob.Decoder
	mid     *dirSnapshotId
}

func (br *dirSnapshotReader) Id() SnapshotId {
	return br.mid
}

func (br *dirSnapshotReader) Read() (Writer, error) {
	var writer Writer
	err := br.decoder.Decode(&writer)
	return writer, err
}

func (br *dirSnapshotReader) Close() error {
	return br.file.Close()
}

type dirSnapshotWriter struct {
	file       *os.File
	writer     *bufio.Writer
	encoder    *gob.Encoder
	id         TransactionId
	repository *DirSnapshotRepository
}

func (bw *dirSnapshotWriter) Id() TransactionId {
	return bw.id
}

func (bw *dirSnapshotWriter) Write(writer Writer) error {
	if bw.encoder == nil {
		return errors.New("gobdb: write() on closed SnapshotWriter")
	}
	return bw.encoder.Encode(&writer)
}

func (bw *dirSnapshotWriter) Close() error {
	if bw.encoder == nil {
		return errors.New("gobdb: close() on closed SnapshotWriter")
	}
	oldname := bw.file.Name()
	newname := fmt.Sprintf(dirSnapshotRepositoryFileNameFormat, bw.id)
	err1 := bw.writer.Flush()
	err2 := bw.file.Close()
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return os.Rename(oldname, filepath.Join(bw.repository.dir, newname))
}
