package work

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
	handler            internalTypes.HandlerFunc
	ordinal            int
	result             error
	completedWorkload  chan<- *worksheet
	cancellationRecord *cancellationRecord
}

func Data(bufferSize int) (errChan chan error, waitChan chan struct{}, workChan chan *worksheet, cancellationRecord *cancellationRecord) {
	errChan = make(chan error)
	waitChan = make(chan struct{})
	workChan = make(chan *worksheet, bufferSize)
	cancellationRecord = newCancellationRecord()

	return
}

func Add(
	handler internalTypes.HandlerFunc,
	ordinal int,
	completedWorkload chan<- *worksheet,
	cancellationRecord *cancellationRecord,
) {
	worksheet := &worksheet{
		handler:            handler,
		ordinal:            ordinal,
		completedWorkload:  completedWorkload,
		cancellationRecord: cancellationRecord,
	}

	pendingWorkload <- worksheet
}

func Manage(
	expectedCount int,
	errChan chan<- error,
	runningChan chan<- struct{},
	workChan <-chan *worksheet,
	cancellationRecord *cancellationRecord,
) {
	seenRecord := newSeenRecord()
	var errOrdinal int = math.MaxInt32
	var err error

	for runningChan <- struct{}{}; ; {
		select {
		case worksheet := <-workChan:
			o := worksheet.ordinal

			if r := worksheet.result; r != nil {
				if o < errOrdinal {
					errOrdinal = o
					err = r

					if seenRecord.HasSeenUpTo(o) {
						goto EndFor
					}
				}
			}

			if seenRecord.SetSeen(o); seenRecord.Count() == expectedCount {
				goto EndFor
			}
		}
	}
EndFor:

	cancellationRecord.Set()
	errChan <- err
}

func runWorker() {
	for {
		select {
		case worksheet := <-pendingWorkload:
			if !worksheet.cancellationRecord.IsSet() {
				if err := worksheet.handler(worksheet.ordinal); err != nil {
					worksheet.result = err
				}

				worksheet.completedWorkload <- worksheet
			}
		}
	}
}

func init() {
	for i := 0; i < numCPU; i++ {
		go runWorker()
	}
}
