package gobdb

import (
	"bytes"
	"encoding/gob"
	"errors"
	"sync"
)

// A container of Snapshots that keeps the data in memory.
type MemSnapshotRepository struct {
	mutex sync.Mutex
	count int
	snaps map[TransactionId]map[*memSnapshotId][]byte
}

// New instance.
func NewMemSnapshotRepository() *MemSnapshotRepository {
	snaps := make(map[TransactionId]map[*memSnapshotId][]byte)
	return &MemSnapshotRepository{sync.Mutex{}, 0, snaps}
}

// Implements SnapshotRepository.Snapshots().
func (r *MemSnapshotRepository) Snapshots() []SnapshotId {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	ids := make([]SnapshotId, 0, r.count)
	for _, m2 := range r.snaps {
		for id := range m2 {
			ids = append(ids, id)
		}
	}
	return ids
}

// Implements SnapshotRepository.ReadSnapshot().
func (r *MemSnapshotRepository) ReadSnapshot(id SnapshotId) (SnapshotReader, error) {
	mid, ok := id.(*memSnapshotId)
	if !ok {
		return nil, errors.New("gobdb: wrong type of SnapshotId on MemSnapshotRepository")
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	m2, ok := r.snaps[mid.id]
	if ok {
		data, ok := m2[mid]
		if ok {
			decoder := gob.NewDecoder(bytes.NewReader(data))
			return &memSnapshotReader{decoder, mid}, nil
		}
	}
	return nil, errors.New("gobdb: SnapshotId not found on MemSnapshotRepository")
}

// Implements WriteSnapshotRepository.WriteSnapshot().
func (r *MemSnapshotRepository) WriteSnapshot(id TransactionId) (SnapshotWriter, error) {
	buffer := &bytes.Buffer{}
	encoder := gob.NewEncoder(buffer)
	return &memSnapshotWriter{encoder, buffer, id, r}, nil
}

type memSnapshotId struct {
	id         TransactionId
	repository *MemSnapshotRepository
}

func (id *memSnapshotId) Id() TransactionId {
	return id.id
}

func (id *memSnapshotId) Repository() SnapshotRepository {
	return id.repository
}

type memSnapshotReader struct {
	decoder *gob.Decoder
	mid     *memSnapshotId
}

func (br *memSnapshotReader) Id() SnapshotId {
	return br.mid
}

func (br *memSnapshotReader) Read() (Writer, error) {
	var writer Writer
	err := br.decoder.Decode(&writer)
	return writer, err
}

func (br *memSnapshotReader) Close() error {
	return nil
}

type memSnapshotWriter struct {
	encoder    *gob.Encoder
	buffer     *bytes.Buffer
	id         TransactionId
	repository *MemSnapshotRepository
}

func (bw *memSnapshotWriter) Id() TransactionId {
	return bw.id
}

func (bw *memSnapshotWriter) Write(writer Writer) error {
	if bw.encoder == nil {
		return errors.New("gobdb: write() on closed SnapshotWriter")
	}
	return bw.encoder.Encode(&writer)
}

func (bw *memSnapshotWriter) Close() error {
	if bw.encoder == nil {
		return errors.New("gobdb: close() on closed SnapshotWriter")
	}
	bw.repository.mutex.Lock()
	defer bw.repository.mutex.Unlock()
	bw.repository.count++
	m2, ok := bw.repository.snaps[bw.id]
	if !ok {
		m2 = make(map[*memSnapshotId][]byte)
		bw.repository.snaps[bw.id] = m2
	}
	mid := &memSnapshotId{bw.id, bw.repository}
	m2[mid] = bw.buffer.Bytes()
	bw.buffer = nil
	bw.encoder = nil
	return nil
}
