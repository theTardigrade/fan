package fan

import (
	"sync"
)

type jobsheet struct {
	result            error
	resultIndex       int
	resultMutex       sync.RWMutex
	workCount         int
	addedWorkCount    chan int
	completedWorkload chan *worksheet
	stopWork          chan struct{}
}

func newJobsheet(workCount int, completedWorkload chan *worksheet) *jobsheet {
	return &jobsheet{
		workCount:         workCount,
		completedWorkload: completedWorkload,
		stopWork:          make(chan struct{}),
		addedWorkCount:    make(chan int, 1),
	}
}

func (j *jobsheet) addOne(handler Handler, index int) (added bool) {
	if j.isWorthStarting(index) {
		pendingWorkload <- newWorksheet(handler, index, j)
		added = true
	}

	return
}

func (jobsheet *jobsheet) Add(handlers []Handler) {
	var addedWorkCount int

	for i, l := 0, jobsheet.workCount; i < l; i++ {
		if jobsheet.addOne(handlers[i], i) {
			addedWorkCount++
		}
	}

	jobsheet.addedWorkCount <- addedWorkCount
}

func (jobsheet *jobsheet) AddRepeated(handler Handler, count int) {
	var addedWorkCount int

	for i, l := 0, jobsheet.workCount; i < l; i++ {
		if jobsheet.addOne(handler, i) {
			addedWorkCount++
		}
	}

	jobsheet.addedWorkCount <- addedWorkCount
}

func (j *jobsheet) isWorthStarting(index int) bool {
	defer j.resultMutex.RUnlock()
	j.resultMutex.RLock()

	return j.result == nil || j.resultIndex > index
}

func (j *jobsheet) SetResult(result error, resultIndex int) {
	defer j.resultMutex.Unlock()
	j.resultMutex.Lock()

	if i := j.resultIndex; resultIndex < i || i == 0 {
		j.result, j.resultIndex = result, resultIndex
	}
}

// runs in own goroutine
func (j *jobsheet) startWork() {
	for {
		select {
		case worksheet := <-pendingWorkload:
			worksheet.Work()
		case <-j.stopWork:
			return
		}
	}
}

func (j *jobsheet) Work() {
	for i := 0; i < numCPU; i++ {
		go j.startWork()
	}
}

func (j *jobsheet) WaitForSome(count int) {
}

func (j *jobsheet) Wait() {
	var workCount int

outerLoop:
	for {
		select {
		case <-j.completedWorkload:
			workCount--
		case wc := <-j.addedWorkCount:
			workCount += wc
			break outerLoop
		}
	}

	workSubcount := workCount / numCPU
	var wg sync.WaitGroup

	wg.Add(numCPU)

	for i := 0; i < numCPU; i++ {
		go func() {
			defer wg.Done()

			j.WaitForSome(workSubcount)
		}()
	}

	wg.Wait()

	if workCount -= workSubcount * numCPU; workCount > 0 {
		j.WaitForSome(workCount)
	}

	for i := 0; i < numCPU; i++ {
		j.stopWork <- struct{}{}
	}
}
