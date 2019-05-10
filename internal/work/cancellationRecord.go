package work

import "sync"

type cancellationRecord struct {
	value bool
	mutex sync.RWMutex
}

func newCancellationRecord() *cancellationRecord {
	return &cancellationRecord{}
}

func (c *cancellationRecord) Set() {
	c.mutex.Lock()
	c.value = true
	c.mutex.Unlock()
}

func (c *cancellationRecord) IsSet() (value bool) {
	c.mutex.RLock()
	value = c.value
	c.mutex.RUnlock()
	return
}
