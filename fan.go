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

func HandleWithAccumulation(handlers []Handler) *accumulationResult {
	count := len(handlers)
	accumulationResult := NewAccumulationData(count)
	accumulationHandler := func(i int) error {
		if err := handlers[i](i); err != nil {
			accumulationResult.AddError(err, i)
		}

		return nil
	}

	jobsheet := newJobsheet(count)

	jobsheet.Start()

	go jobsheet.AddRepeated(accumulationHandler)

	jobsheet.Wait()
	accumulationResult.ShrinkData()

	return accumulationResult
}

func HandleRepeatedWithAccumulation(handler Handler, count int) *accumulationResult {
	accumulationResult := NewAccumulationData(count)
	accumulationHandler := func(i int) error {
		if err := handler(i); err != nil {
			accumulationResult.AddError(err, i)
		}
		return nil
	}

	jobsheet := newJobsheet(count)

	jobsheet.Start()

	go jobsheet.AddRepeated(accumulationHandler)

	jobsheet.Wait()
	accumulationResult.ShrinkData()

	return accumulationResult
}
