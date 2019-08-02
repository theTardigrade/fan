package fan

import (
	"runtime"
)

var (
	numCPU = runtime.NumCPU()
)

type worksheet struct {
	handler  Handler
	index    int
	result   error
	jobsheet *jobsheet
}

func newWorksheet(
	handler Handler,
	index int,
	jobsheet *jobsheet,
) *worksheet {
	return &worksheet{
		handler:  handler,
		index:    index,
		jobsheet: jobsheet,
	}
}

func (w *worksheet) Work() {
	j := w.jobsheet

	if i := w.index; j.isWorthStarting(i) {
		if err := w.handler(i); err != nil {
			w.result = err
			j.SetResult(err, i)
		}
	}

	j.completedWorkload <- w
}
