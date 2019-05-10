package fan

import (
	internalTypes "github.com/theTardigrade/fan/internal/types"
	internalWork "github.com/theTardigrade/fan/internal/work"
)

type HandlerFunc internalTypes.HandlerFunc

func Handle(handlers ...HandlerFunc) (err error) {
	if l := len(handlers); l > 0 {
		errChan, waitChan, workChan := internalWork.Chans(l)

		go internalWork.Manage(l, errChan, waitChan, workChan)

		<-waitChan

		for i := 0; i < l; i++ {
			go internalWork.Add(internalTypes.HandlerFunc(handlers[i]), i, workChan)
		}

		err = <-errChan
	}

	return
}

func HandleRepeated(handler HandlerFunc, repeats int) (err error) {
	if repeats > 0 {
		errChan, waitChan, workChan := internalWork.Chans(repeats)

		go internalWork.Manage(repeats, errChan, waitChan, workChan)

		<-waitChan

		for i := 0; i < repeats; i++ {
			go internalWork.Add(internalTypes.HandlerFunc(handler), i, workChan)
		}

		err = <-errChan
	}

	return
}
