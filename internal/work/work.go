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

func Add(
	handler internalTypes.HandlerFunc,
	ordinal int,
	datum *datum,
) {
	worksheet := &worksheet{
		handler:            handler,
		ordinal:            ordinal,
		completedWorkload:  datum.workChan,
		cancellationRecord: datum.cancellationRecord,
	}

	pendingWorkload <- worksheet
}

func Manage(
	datum *datum,
) {
	seenRecord := newSeenRecord()
	var errOrdinal int = math.MaxInt32
	var err error

	for datum.WaitChan <- struct{}{}; ; {
		select {
		case worksheet := <-datum.workChan:
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

			if seenRecord.SetSeen(o); seenRecord.Count() == datum.expectedCount {
				goto EndFor
			}
		}
	}
EndFor:

	datum.cancellationRecord.Set()
	datum.ErrChan <- err
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
