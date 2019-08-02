package fan

type Handler (func(int) error)

func Handle(handlers []Handler) error {
	count := len(handlers)
	completedWorkload := make(chan *worksheet, numCPU)
	jobsheet := newJobsheet(count, completedWorkload)

	jobsheet.Work()

	go jobsheet.Add(handlers)

	jobsheet.Wait()

	return jobsheet.result
}

func HandleRepeated(handler Handler, count int) error {
	completedWorkload := make(chan *worksheet, numCPU)
	jobsheet := newJobsheet(count, completedWorkload)

	jobsheet.Work()

	go jobsheet.AddRepeated(handler, count)

	jobsheet.Wait()

	err := jobsheet.result

	return err
}
