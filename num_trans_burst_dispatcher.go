package gobdb

// BurstDispatcher that writes to another and is able to rotate
// by number of transactions.
// No thread-safe.
type NumTransactionsBurstDispatcher struct {
	count, max int
	dispatcher BurstDispatcher
}

// New instance. It rotates after a number of Transactions has been written.
func NewNumTransactionsBurstDispatcher(max int, dispatcher BurstDispatcher) *NumTransactionsBurstDispatcher {
	return &NumTransactionsBurstDispatcher{0, max, dispatcher}
}

// Implements BurstDispatcher.Write().
func (bd *NumTransactionsBurstDispatcher) Write(transaction Transaction) (err error) {
	if err = bd.dispatcher.Write(transaction); err != nil {
		return
	}
	if bd.count++; bd.count >= bd.max {
		return bd.Rotate()
	}
	return
}

// Implements BurstDispatcher.Rotate().
func (bd *NumTransactionsBurstDispatcher) Rotate() (err error) {
	bd.count = 0
	return bd.dispatcher.Rotate()
}

// Implements BurstDispatcher.Close().
func (bd *NumTransactionsBurstDispatcher) Close() (err error) {
	return bd.dispatcher.Close()
}
