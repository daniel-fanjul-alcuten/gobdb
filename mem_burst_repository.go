package gobdb

import (
	"bytes"
	"encoding/gob"
	"errors"
	"sync"
)

// A BurstRepository and WriteBurstRepository that keeps the data in memory.
// Thread-safe, but BurstReaders and BurstWriters are not.
type MemBurstRepository struct {
	mutex  sync.Mutex
	count  int
	bursts map[TransactionId]map[TransactionId]map[*memBurstId][]byte
}

// New instance.
func NewMemBurstRepository() *MemBurstRepository {
	bursts := make(map[TransactionId]map[TransactionId]map[*memBurstId][]byte)
	return &MemBurstRepository{sync.Mutex{}, 0, bursts}
}

// Implements BurstRepository.Bursts().
func (r *MemBurstRepository) Bursts() ([]BurstId, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	ids := make([]BurstId, 0, r.count)
	for _, m2 := range r.bursts {
		for _, m3 := range m2 {
			for id := range m3 {
				ids = append(ids, id)
			}
		}
	}
	return ids, nil
}

// Implements WriteBurstRepository.WriteBurst().
func (r *MemBurstRepository) WriteBurst() (BurstWriter, error) {
	buffer := &bytes.Buffer{}
	encoder := gob.NewEncoder(buffer)
	return &memBurstWriter{encoder, buffer, 0, 0, r}, nil
}

type memBurstId struct {
	first, last TransactionId
	repository  *MemBurstRepository
}

func (id *memBurstId) First() TransactionId {
	return id.first
}

func (id *memBurstId) Last() TransactionId {
	return id.last
}

func (id *memBurstId) Repository() BurstRepository {
	return id.repository
}

func (id *memBurstId) Read() (BurstReader, error) {
	id.repository.mutex.Lock()
	defer id.repository.mutex.Unlock()
	m2, ok := id.repository.bursts[id.first]
	if ok {
		m3, ok := m2[id.last]
		if ok {
			data, ok := m3[id]
			if ok {
				decoder := gob.NewDecoder(bytes.NewReader(data))
				return &memBurstReader{decoder, id}, nil
			}
		}
	}
	return nil, errors.New("gobdb: BurstId not found on MemBurstRepository")
}

type memBurstReader struct {
	decoder *gob.Decoder
	mid     *memBurstId
}

func (br *memBurstReader) Id() BurstId {
	return br.mid
}

func (br *memBurstReader) Read() (Transaction, error) {
	var transaction Transaction
	err := br.decoder.Decode(&transaction)
	return transaction, err
}

func (br *memBurstReader) Close() error {
	return nil
}

type memBurstWriter struct {
	encoder     *gob.Encoder
	buffer      *bytes.Buffer
	first, last TransactionId
	repository  *MemBurstRepository
}

func (bw *memBurstWriter) First() TransactionId {
	return bw.first
}

func (bw *memBurstWriter) Last() TransactionId {
	return bw.last
}

func (bw *memBurstWriter) Write(transaction Transaction) error {
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

func (bw *memBurstWriter) Close() error {
	if bw.encoder == nil {
		return errors.New("gobdb: close() on closed BurstWriter")
	}
	if bw.last == 0 {
		return nil
	}
	bw.repository.mutex.Lock()
	defer bw.repository.mutex.Unlock()
	bw.repository.count++
	m2, ok := bw.repository.bursts[bw.first]
	if !ok {
		m2 = make(map[TransactionId]map[*memBurstId][]byte)
		bw.repository.bursts[bw.first] = m2
	}
	m3, ok := m2[bw.last]
	if !ok {
		m3 = make(map[*memBurstId][]byte)
		m2[bw.last] = m3
	}
	mid := &memBurstId{bw.first, bw.last, bw.repository}
	m3[mid] = bw.buffer.Bytes()
	bw.buffer = nil
	bw.encoder = nil
	return nil
}
