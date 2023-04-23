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
	return &savePortsServiceSQLTransactionInstanceFactory{
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

	ctx, cancelCtx := context.WithCancel(ctx)

	return service.NewSavePortsServiceInstance(
		savePortsToDBChann,
		service.NewSavePortToDBParam(
			ctx,
			&wg,
			spsif.savePortsRepository,
			errorChann,
		),
		&wg,
		cancelCtx,
		errorChann,
		nil,
		nil,
		spsif.log,
	), nil
}
