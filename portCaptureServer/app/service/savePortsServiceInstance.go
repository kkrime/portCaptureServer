package service

import (
	"context"
	"portCaptureServer/app/entity"
	"sync"

	"github.com/sirupsen/logrus"
)

type savePortsServiceInstance struct {
	savePortsToDBChann chan<- *savePortToDBParam
	// NOTE savePortToDBParam is added to keep things tidy/organized
	// and for convience when sending to savePortsToDBChann
	// (we can just copy savePortToDBParam and pass in the address of the copy)
	savePortToDBParam savePortToDBParam
	cancelCtx         context.CancelFunc
	wg                *sync.WaitGroup
	errorChann        chan error
	finalizeDB        func() error
	cancelDB          func() error
	log               *logrus.Logger
}

// func newSavePortsServiceInstance(ctx context.Context,
// 	savePortsToDBChann chan<- *savePortToDBParam,
// 	savePortToDBParam savePortToDBParam,
// 	cancelCtx context.CancelFunc,
// 	wg *sync.WaitGroup,
// 	db repository.SavePortsRepository,
// 	errorChann chan error,
// 	finalizeDB func() error,
// 	cancelDB func() error,
// 	log *logrus.Logger,
// ) SavePortsServiceInstance {
// 	return &savePortsServiceInstance{
// 		savePortsToDBChann: savePortsToDBChann,
// 		savePortToDBParam:  savePortToDBParam,
// 		cancelCtx:          cancelCtx,
// 		wg:                 wg,
// 		errorChann:         errorChann,
// 		finalizeDB:         finalizeDB,
// 		cancelDB:           cancelDB,
// 		log:                log,
// 	}
// }

func (spsi *savePortsServiceInstance) SavePort(port *entity.Port) error {
	select {
	// check if any errors occured
	case err := <-spsi.errorChann:
		spsi.log.Errorf("Error Occured, No Ports Were Saved To The Database: %s", err.Error())

		// cancel all database io
		spsi.cancelCtx()

		spsi.wg.Wait()
		close(spsi.errorChann)

		// drain the errorChann
		for len(spsi.errorChann) > 0 {
			<-spsi.errorChann
		}

		// roll back the transaction
		if spsi.cancelDB != nil {
			spsi.cancelDB()
		}

		return err
	default:
		// copy spsi.savePortToDBParam
		savePortToDBParam := spsi.savePortToDBParam
		savePortToDBParam.port = port

		spsi.wg.Add(1)
		spsi.savePortsToDBChann <- &savePortToDBParam
		return nil
	}
}

func (spsi *savePortsServiceInstance) Finalize() error {
	spsi.wg.Wait()
	close(spsi.errorChann)

	for err := range spsi.errorChann {

		if err != nil {
			spsi.log.Errorf("Error Occured, No Ports Were Saved To The Database: %s", err.Error())

			// drain the errorChann
			for len(spsi.errorChann) > 0 {
				<-spsi.errorChann
			}

			// any cancel procedures
			if spsi.cancelDB != nil {
				spsi.cancelDB()
			}

			return err
		}
	}

	// any finalization procedures
	if spsi.finalizeDB != nil {
		spsi.finalizeDB()
	}

	spsi.log.Infof("All Ports Successfully Saved To The Database!!")

	return nil
}
