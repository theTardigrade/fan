package work

import "sync"

type seenRecord struct {
	data  map[int]bool
	mutex sync.RWMutex
}

func newSeenRecord() *seenRecord {
	return &seenRecord{
		data: make(map[int]bool),
	}
}

func (r *seenRecord) HasSeenUpTo(n int) bool {
	defer r.mutex.RUnlock()
	r.mutex.RLock()

	for i := len(r.data) - 1; i >= 0; i-- {
		if seen := r.data[n]; !seen {
			return false
		}
	}

	return true
}

func (r *seenRecord) SetSeen(n int) {
	r.mutex.Lock()
	r.data[n] = true
	r.mutex.Unlock()
}

func (r *seenRecord) Count() int {
	return len(r.data)
}
