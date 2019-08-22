package fan

import (
	"sort"
	"sync"
)

type accumulationResultDatum struct {
	Err   error
	Index int
}

type accumulationResult struct {
	Data  []accumulationResultDatum
	mutex sync.Mutex
}

func (a *accumulationResult) Len() int           { return len(a.Data) }
func (a *accumulationResult) Swap(i, j int)      { a.Data[i], a.Data[j] = a.Data[j], a.Data[i] }
func (a *accumulationResult) Less(i, j int) bool { return a.Data[i].Index < a.Data[j].Index }

func NewAccumulationData(workCount int) *accumulationResult {
	return &accumulationResult{
		Data: make([]accumulationResultDatum, 0, workCount),
	}
}

// concurrent
func (a *accumulationResult) AddError(err error, Index int) {
	defer a.mutex.Unlock()
	a.mutex.Lock()

	a.Data = append(a.Data, accumulationResultDatum{
		Err:   err,
		Index: Index,
	})
}

// non-concurrent
func (a *accumulationResult) ShrinkData() {
	a.Data = a.Data[:len(a.Data)]
}

// non-concurrent
func (a *accumulationResult) Errors() []error {
	sort.Sort(a)

	l := len(a.Data)
	e := make([]error, l)

	for i := 0; i < l; i++ {
		e[i] = a.Data[i].Err
	}

	return e
}
