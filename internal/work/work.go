package fan

import (
	"math"
	"runtime"

	internalTypes "github.com/theTardigrade/fan/internal/types"
)

var (
	numCPU          = runtime.NumCPU()
	pendingWorkload = make(chan *worksheet, numCPU)
)

type worksheet struct {
	handler           internalTypes.HandlerFunc
	ordinal           int
	result            error
	completedWorkload chan<- *worksheet
}

func Chans(bufferSize int) (errChan chan error, waitChan chan struct{}, workChan chan *worksheet) {
	errChan = make(chan error)
	waitChan = make(chan struct{})
	workChan = make(chan *worksheet, bufferSize)

	return
}

func Add(
	handler internalTypes.HandlerFunc,
	ordinal int,
	completedWorkload chan<- *worksheet,
) {
	worksheet := &worksheet{
		handler:           handler,
		ordinal:           ordinal,
		completedWorkload: completedWorkload,
	}

	pendingWorkload <- worksheet
}

func Manage(expectedCount int, errChan chan<- error, runningChan chan<- struct{}, workChan <-chan *worksheet) {
	var seenCount int
	var errOrdinal int = math.MaxInt32
	var err error

	for runningChan <- struct{}{}; ; {
		select {
		case worksheet := <-workChan:
			if r := worksheet.result; r != nil {
				if o := worksheet.ordinal; o < errOrdinal {
					errOrdinal = o
					err = r
				}
			}

			if seenCount++; seenCount == expectedCount {
				errChan <- err
				return
			}
		}
	}
}

func runWorker() {
	for {
		select {
		case worksheet := <-pendingWorkload:
			if err := worksheet.handler(worksheet.ordinal); err != nil {
				worksheet.result = err
			}
			worksheet.completedWorkload <- worksheet
		}
	}
}

func init() {
	for i := 0; i < numCPU; i++ {
		go runWorker()
	}
}
