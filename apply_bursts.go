package gobdb

import (
	"io"
)

type applyBurstsReader struct {
	Transaction
	BurstReader
}

// Applies bursts in order to a Root object. It receives and returns the last
// TransactionId applied to the Root. It sorts the []BurstId with SortBursts().
func ApplyBursts(root Root, transactionId TransactionId, burstIds []BurstId) (last TransactionId, err error) {

	last = transactionId
	next := last + 1
	SortBursts(burstIds)

	readers := []applyBurstsReader{}
	defer func() {
		for _, r := range readers {
			r.Close()
		}
	}()

	for {

		var (
			index       int
			transaction *Transaction
		)
		for index < len(readers) {
			transaction = &readers[index].Transaction
			for transaction.Id <= last {
				if *transaction, err = readers[index].Read(); err != nil {
					if err != io.EOF {
						return
					}
					err = nil
					lastIndex := len(readers) - 1
					if index < lastIndex {
						readers[index], readers[lastIndex] = readers[lastIndex], applyBurstsReader{}
					}
					readers = readers[:lastIndex]
					transaction = nil
					break
				}
			}
			if transaction != nil {
				if transaction.Id == next {
					break
				}
				transaction = nil
				index++
			}
		}

		if transaction != nil {
			if _, err = transaction.Write(root); err != nil {
				return
			}
			last, next = next, next+1
			continue
		}

		open := false
		for i := 0; !open && i < len(burstIds); {
			id := burstIds[i]
			skip := next > id.Last()
			open = id.First() <= next && !skip
			if open {
				var reader applyBurstsReader
				if reader.BurstReader, err = id.Repository().ReadBurst(id); err != nil {
					return
				}
				readers = append(readers, reader)
				skip = true
			}
			if skip {
				lastIndex := len(burstIds) - 1
				if index < lastIndex {
					burstIds[i], burstIds[lastIndex] = burstIds[lastIndex], nil
				}
				burstIds = burstIds[:lastIndex]
			} else {
				i++
			}
		}

		if !open {
			break
		}
	}

	for len(readers) > 0 {
		if err, readers = readers[0].Close(), readers[1:]; err != nil {
			return
		}
	}

	return
}
