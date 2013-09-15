package gobdb

// It writes to only one Burst. There is no automatic rotation.
// No thread-safe.
type DefaultBurstDispatcher struct {
	burst      BurstWriter
	repository WriteBurstRepository
}

// New instance.
func NewDefaultBurstDispatcher(repository WriteBurstRepository) *DefaultBurstDispatcher {
	return &DefaultBurstDispatcher{nil, repository}
}

// Implements BurstDispatcher.Write().
func (bd *DefaultBurstDispatcher) Write(transaction Transaction) (err error) {
	if bd.burst == nil {
		bd.burst, err = bd.repository.WriteBurst()
		if err != nil {
			return
		}
	}
	return bd.burst.Write(transaction)
}

// Implements BurstDispatcher.Rotate().
func (bd *DefaultBurstDispatcher) Rotate() (err error) {
	if bd.burst == nil {
		return
	}
	if err = bd.burst.Close(); err != nil {
		return
	}
	bd.burst = nil
	return
}

// Implements BurstDispatcher.Close().
func (bd *DefaultBurstDispatcher) Close() (err error) {
	if bd.burst == nil {
		return
	}
	err = bd.burst.Close()
	bd.burst = nil
	return
}
