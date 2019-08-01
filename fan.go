package fan

import (
	internalTypes "github.com/theTardigrade/fan/internal/types"
	internalWork "github.com/theTardigrade/fan/internal/work"
)

type HandlerFunc internalTypes.HandlerFunc

func Handle(handlers ...HandlerFunc) (err error) {
	if l := len(handlers); l > 0 {
		datum := internalWork.NewDatum(l)

		go internalWork.Manage(datum)

		<-datum.WaitChan

		for i := 0; i < l; i++ {
			go internalWork.Add(internalTypes.HandlerFunc(handlers[i]), i, datum)
		}

		err = <-datum.ErrChan
	}

	return
}

func HandleRepeated(handler HandlerFunc, repeats int) (err error) {
	if repeats > 0 {
		datum := internalWork.NewDatum(repeats)

		go internalWork.Manage(datum)

		<-datum.WaitChan

		for i := 0; i < repeats; i++ {
			go internalWork.Add(internalTypes.HandlerFunc(handler), i, datum)
		}

		err = <-datum.ErrChan
	}

	return
}
