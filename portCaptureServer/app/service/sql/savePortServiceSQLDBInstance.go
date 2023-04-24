package sql

import (
	"context"
	"portCaptureServer/app/repository"
	"portCaptureServer/app/service"
	"sync"

	"github.com/sirupsen/logrus"
)

type savePortsServiceSQLInstanceFactory struct {
	savePortsRepository repository.SavePortsRepository
	log                 *logrus.Logger
}

func NewSavePortsServiceSQLInstanceFactory(
	savePortsRepository repository.SavePortsRepository,
	log *logrus.Logger,
) service.SavePortsServiceInstanceFactory {
	return &savePortsServiceSQLInstanceFactory{
		savePortsRepository: savePortsRepository,
		log:                 log,
	}
}

func (spsif *savePortsServiceSQLInstanceFactory) NewSavePortsInstance(
	ctx context.Context,
	savePortsToDBChann chan<- *service.SavePortToDBParam,
) (service.SavePortsServiceInstance, error) {
	var wg sync.WaitGroup
	errorChann := make(chan error, 1)
	resultChann := make(chan error)

	ctx, cancelCtx := context.WithCancel(ctx)

	go func() {
		// only send the first error
		firstErrorOccured := false
		for err := range resultChann {
			if !firstErrorOccured && err != nil {
				firstErrorOccured = true
				// cancel context to stop any db io
				cancelCtx()
				// if any errors on saving the ports to the db, send it to
				// errorChann
				errorChann <- err
			}
		}
		close(errorChann)
	}()

	return service.NewSavePortsServiceInstance(
		savePortsToDBChann,
		service.NewSavePortToDBParam(
			ctx,
			&wg,
			spsif.savePortsRepository,
			resultChann,
		),
		&wg,
		errorChann,
		resultChann,
		nil,
		nil,
		spsif.log,
	), nil
}
