package fan

type Handler (func(int) error)

func Handle(handlers []Handler) error {
	count := len(handlers)
	jobsheet := newJobsheet(count)

	jobsheet.Work()

	go jobsheet.Add(handlers)

	jobsheet.Wait()

	return jobsheet.result
}

func HandleRepeated(handler Handler, count int) error {
	jobsheet := newJobsheet(count)

	jobsheet.Work()

	go jobsheet.AddRepeated(handler, count)

	jobsheet.Wait()

	return jobsheet.result
}
