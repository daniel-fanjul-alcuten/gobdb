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

const dirBurstRepositoryFileNameFormat = "burst-%d-%d.gobdb"
const dirBurstRepositoryFileNameScanFormat = dirBurstRepositoryFileNameFormat + "\n"

// A BurstRepository and WriteBurstRepository that uses one file per Burst.
// Thread-safe, but BurstReaders and BurstWriters are not.
type DirBurstRepository struct {
	dir string
}

func NewDirBurstRepository(dir string) *DirBurstRepository {
	return &DirBurstRepository{dir}
}

func (r *DirBurstRepository) Bursts() ([]BurstId, error) {
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
	ids := make([]BurstId, 0, len(infos))
	for _, info := range infos {
		name := info.Name()
		var first, last int
		if n, err := fmt.Sscanf(name, dirBurstRepositoryFileNameScanFormat, &first, &last); n == 2 && err == nil {
			ids = append(ids, &dirBurstId{TransactionId(first), TransactionId(last), r})
		}
	}
	return ids, nil
}

func (r *DirBurstRepository) WriteBurst() (BurstWriter, error) {
	file, err := ioutil.TempFile(r.dir, "tmp-burst-")
	if err != nil {
		return nil, err
	}
	writer := bufio.NewWriter(file)
	encoder := gob.NewEncoder(writer)
	return &dirBurstWriter{file, writer, encoder, 0, 0, r}, nil
}

type dirBurstId struct {
	first, last TransactionId
	repository  *DirBurstRepository
}

func (id *dirBurstId) First() TransactionId {
	return id.first
}

func (id *dirBurstId) Last() TransactionId {
	return id.last
}

func (id *dirBurstId) Repository() BurstRepository {
	return id.repository
}

func (id *dirBurstId) Read() (BurstReader, error) {
	name := fmt.Sprintf(dirBurstRepositoryFileNameFormat, id.first, id.last)
	file, err := os.Open(filepath.Join(id.repository.dir, name))
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	decoder := gob.NewDecoder(reader)
	return &dirBurstReader{file, decoder, id}, nil
}

type dirBurstReader struct {
	file    *os.File
	decoder *gob.Decoder
	mid     *dirBurstId
}

func (br *dirBurstReader) Id() BurstId {
	return br.mid
}

func (br *dirBurstReader) Read() (Transaction, error) {
	var transaction Transaction
	err := br.decoder.Decode(&transaction)
	return transaction, err
}

func (br *dirBurstReader) Close() error {
	return br.file.Close()
}

type dirBurstWriter struct {
	file        *os.File
	writer      *bufio.Writer
	encoder     *gob.Encoder
	first, last TransactionId
	repository  *DirBurstRepository
}

func (bw *dirBurstWriter) First() TransactionId {
	return bw.first
}

func (bw *dirBurstWriter) Last() TransactionId {
	return bw.last
}

func (bw *dirBurstWriter) Write(transaction Transaction) error {
	if transaction.Id <= bw.last {
		return errors.New("gobdb: write() of transaction with invalid id")
	}
	if bw.encoder == nil {
		return errors.New("gobdb: write() on closed BurstWriter")
	}
	err := bw.encoder.Encode(&transaction)
	if err == nil {
		if bw.first == 0 {
			bw.first = transaction.Id
		}
		bw.last = transaction.Id
	}
	return err
}

func (bw *dirBurstWriter) Close() error {
	if bw.encoder == nil {
		return errors.New("gobdb: close() on closed BurstWriter")
	}
	oldname := bw.file.Name()
	if bw.last == 0 {
		err1 := bw.file.Close()
		err2 := os.Remove(oldname)
		if err1 != nil {
			return err1
		}
		if err2 != nil {
			return err2
		}
		return nil
	}
	newname := fmt.Sprintf(dirBurstRepositoryFileNameFormat, bw.first, bw.last)
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
