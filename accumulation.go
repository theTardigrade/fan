package fan

import (
	"sort"
	"sync"
)

type accumulationResultDatum struct {
	Err   error
	Index int
}

type AccumulationResult struct {
	Data  []accumulationResultDatum
	mutex sync.Mutex
}

func (a *AccumulationResult) Len() int           { return len(a.Data) }
func (a *AccumulationResult) Swap(i, j int)      { a.Data[i], a.Data[j] = a.Data[j], a.Data[i] }
func (a *AccumulationResult) Less(i, j int) bool { return a.Data[i].Index < a.Data[j].Index }

func NewAccumulationData(workCount int) *AccumulationResult {
	return &AccumulationResult{
		Data: make([]accumulationResultDatum, 0, workCount),
	}
}

// concurrent
func (a *AccumulationResult) addError(err error, Index int) {
	defer a.mutex.Unlock()
	a.mutex.Lock()

	a.Data = append(a.Data, accumulationResultDatum{
		Err:   err,
		Index: Index,
	})
}

// non-concurrent
func (a *AccumulationResult) handler(rawHandlers []Handler) Handler {
	return func(i int) error {
		if err := rawHandlers[i](i); err != nil {
			a.addError(err, i)
		}

		return nil
	}
}

// non-concurrent
func (a *AccumulationResult) handlerRepeated(rawHandler Handler) Handler {
	return func(i int) error {
		if err := rawHandler(i); err != nil {
			a.addError(err, i)
		}

		return nil
	}
}

// non-concurrent
func (a *AccumulationResult) shrinkData() {
	a.Data = a.Data[:len(a.Data)]
}

// non-concurrent
func (a *AccumulationResult) Errors() []error {
	sort.Sort(a)

	l := len(a.Data)
	e := make([]error, l)

	for i := 0; i < l; i++ {
		e[i] = a.Data[i].Err
	}

	return e
}
