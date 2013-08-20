package gobdb

// It writes to only one Burst. No rotation.
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

// Implements BurstDispatcher.Close().
func (bd *DefaultBurstDispatcher) Close() (err error) {
	if bd.burst == nil {
		return
	}
	if err = bd.burst.Close(); err != nil {
		return
	}
	bd.burst = nil
	return
}
