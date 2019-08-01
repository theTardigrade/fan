package work

import (
	"sync"
)

type seenRecord struct {
	data      map[int]bool
	mutex     sync.RWMutex
	upToValue int
}

func newSeenRecord() *seenRecord {
	return &seenRecord{
		data:      make(map[int]bool),
		upToValue: -1,
	}
}

func (r *seenRecord) HasSeenUpTo(n int) bool {
	defer r.mutex.Unlock()
	r.mutex.Lock()

	for i, j := n-1, r.upToValue; i > j; i-- {
		if seen := r.data[i]; !seen {
			return false
		}
	}

	if r.upToValue < n {
		r.upToValue = n
	}

	return true
}

func (r *seenRecord) SetSeen(n int) {
	r.mutex.Lock()
	r.data[n] = true
	r.mutex.Unlock()
}

func (r *seenRecord) Count() (l int) {
	r.mutex.RLock()
	l = len(r.data)
	r.mutex.RUnlock()

	return
}
