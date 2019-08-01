package work

type datum struct {
	ErrChan            chan error
	WaitChan           chan struct{}
	workChan           chan *worksheet
	cancellationRecord *cancellationRecord
	expectedCount      int
}

func NewDatum(expectedCount int) *datum {
	return &datum{
		ErrChan:            make(chan error),
		WaitChan:           make(chan struct{}),
		workChan:           make(chan *worksheet, expectedCount),
		cancellationRecord: newCancellationRecord(),
		expectedCount:      expectedCount,
	}
}
