package service

import (
	"fmt"
	"portCaptureServer/app/entity"
	"sync"

	"github.com/sirupsen/logrus"
)

type savePortsServiceInstance struct {
	savePortsToDBChann chan<- *SavePortToDBParam
	// NOTE savePortToDBParam is added to keep things tidy/organized
	// and for convience when sending to savePortsToDBChann
	// (we can just copy savePortToDBParam and pass in the address of the copy)
	savePortToDBParam SavePortToDBParam
	wg                *sync.WaitGroup
	errorChann        chan error
	resultChann       chan error
	finalizeDB        func() error
	cancelDB          func() error
	log               *logrus.Logger
}

func NewSavePortsServiceInstance(
	savePortsToDBChann chan<- *SavePortToDBParam,
	savePortToDBParam SavePortToDBParam,
	wg *sync.WaitGroup,
	errorChann chan error,
	resultChann chan error,
	finalizeDB func() error,
	cancelDB func() error,
	log *logrus.Logger,
) SavePortsServiceInstance {
	return &savePortsServiceInstance{
		savePortsToDBChann: savePortsToDBChann,
		savePortToDBParam:  savePortToDBParam,
		wg:                 wg,
		errorChann:         errorChann,
		resultChann:        resultChann,
		finalizeDB:         finalizeDB,
		cancelDB:           cancelDB,
		log:                log,
	}
}

func (spsi *savePortsServiceInstance) SavePort(port *entity.Port) error {
	select {
	// check if any errors occured
	case err := <-spsi.errorChann:
		spsi.log.Errorf("Error Occured, No Ports Were Saved To The Database: %s", err.Error())

		spsi.wg.Wait()
		close(spsi.resultChann)

		// any cancel procedures
		if spsi.cancelDB != nil {
			cancelDBErr := spsi.cancelDB()
			if cancelDBErr != nil {
				err = fmt.Errorf("%v: %w", err, cancelDBErr)
			}
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
	close(spsi.resultChann)

	for err := range spsi.errorChann {

		if err != nil {
			spsi.log.Errorf("Error Occured, No Ports Were Saved To The Database: %s", err.Error())

			// any cancel procedures
			if spsi.cancelDB != nil {
				cancelDBErr := spsi.cancelDB()
				if cancelDBErr != nil {
					err = fmt.Errorf("%v: %w", err, cancelDBErr)
				}
			}

			return err
		}
	}

	// any finalization procedures
	if spsi.finalizeDB != nil {
		err := spsi.finalizeDB()
		if err != nil {
			return err
		}
	}

	spsi.log.Infof("All Ports Successfully Saved To The Database!!")

	return nil
}
