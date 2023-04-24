package sql

import (
	"context"
	"portCaptureServer/app/repository"
	"portCaptureServer/app/service"

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
	return service.InitalizeNewSavePortsServiceInstance(
		ctx,
		savePortsToDBChann,
		spsif.savePortsRepository,
		nil,
		nil,
		spsif.log,
	), nil
}
