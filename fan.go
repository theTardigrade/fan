package fan

import "sync"

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

func HandleWithAccumulation(handlers []Handler) (errors []error) {
	var errorsMutex sync.Mutex

	accumulationHandler := func(i int) error {
		if err := handlers[i](i); err != nil {
			defer errorsMutex.Unlock()
			errorsMutex.Lock()

			errors = append(errors, err)
		}

		return nil
	}

	count := len(handlers)
	jobsheet := newJobsheet(count)

	jobsheet.Start()

	go jobsheet.AddRepeated(accumulationHandler)

	jobsheet.Wait()

	return
}

func HandleRepeatedWithAccumulation(handler Handler, count int) (errors []error) {
	var errorsMutex sync.Mutex

	accumulationHandler := func(i int) error {
		if err := handler(i); err != nil {
			defer errorsMutex.Unlock()
			errorsMutex.Lock()

			errors = append(errors, err)
		}

		return nil
	}

	jobsheet := newJobsheet(count)

	jobsheet.Start()

	go jobsheet.AddRepeated(accumulationHandler)

	jobsheet.Wait()

	return
}
