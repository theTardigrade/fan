package fan

import (
	internalTypes "github.com/theTardigrade/fan/internal/types"
	internalWork "github.com/theTardigrade/fan/internal/work"
)

type HandlerFunc internalTypes.HandlerFunc

func Handle(handlers ...HandlerFunc) (err error) {
	if l := len(handlers); l > 0 {
		errChan, waitChan, workChan, cancellationRecord := internalWork.Data(l)

		go internalWork.Manage(l, errChan, waitChan, workChan, cancellationRecord)

		<-waitChan

		for i := 0; i < l; i++ {
			go internalWork.Add(internalTypes.HandlerFunc(handlers[i]), i, workChan, cancellationRecord)
		}

		err = <-errChan
	}

	return
}

func HandleRepeated(handler HandlerFunc, repeats int) (err error) {
	if repeats > 0 {
		errChan, waitChan, workChan, cancellationRecord := internalWork.Data(repeats)

		go internalWork.Manage(repeats, errChan, waitChan, workChan, cancellationRecord)

		<-waitChan

		for i := 0; i < repeats; i++ {
			go internalWork.Add(internalTypes.HandlerFunc(handler), i, workChan, cancellationRecord)
		}

		err = <-errChan
	}

	return
}
