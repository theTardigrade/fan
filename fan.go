package fan

type Handler (func(int) error)

func Handle(handlers []Handler) error {
	count := len(handlers)
	jobsheet := newJobsheet(count)

	jobsheet.Start()

	go jobsheet.Add(handlers)

	jobsheet.Wait()

	return jobsheet.result
}

func HandleRepeated(handler Handler, count int) error {
	jobsheet := newJobsheet(count)

	jobsheet.Start()

	go jobsheet.AddRepeated(handler)

	jobsheet.Wait()

	return jobsheet.result
}

func HandleWithAccumulation(handlers []Handler) *AccumulationResult {
	count := len(handlers)
	accumulationResult := NewAccumulationData(count)
	accumulationHandler := accumulationResult.handler(handlers)
	jobsheet := newJobsheet(count)

	jobsheet.Start()

	go jobsheet.AddRepeated(accumulationHandler)

	jobsheet.Wait()
	accumulationResult.shrinkData()

	return accumulationResult
}

func HandleRepeatedWithAccumulation(handler Handler, count int) *AccumulationResult {
	accumulationResult := NewAccumulationData(count)
	accumulationHandler := accumulationResult.handlerRepeated(handler)
	jobsheet := newJobsheet(count)

	jobsheet.Start()

	go jobsheet.AddRepeated(accumulationHandler)

	jobsheet.Wait()
	accumulationResult.shrinkData()

	return accumulationResult
}
