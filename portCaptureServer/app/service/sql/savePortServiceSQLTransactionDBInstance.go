package sql

import (
	"context"
	"portCaptureServer/app/repository"
	"portCaptureServer/app/repository/sql"
	"portCaptureServer/app/service"

	"github.com/sirupsen/logrus"
)

type savePortsServiceSQLTransactionInstanceFactory struct {
	savePortsRepository repository.SavePortsRepository
	log                 *logrus.Logger
}

func NewSavePortsServiceSQLTransactionInstanceFactory(
	savePortsRepository repository.SavePortsRepository,
	log *logrus.Logger,
) service.SavePortsServiceInstanceFactory {
	return &savePortsServiceSQLTransactionInstanceFactory{
		savePortsRepository: savePortsRepository,
		log:                 log,
	}
}

func (spsif *savePortsServiceSQLTransactionInstanceFactory) NewSavePortsInstance(
	ctx context.Context,
	savePortsToDBChann chan<- *service.SavePortToDBParam,
) (service.SavePortsServiceInstance, error) {

	db := spsif.savePortsRepository.(sql.DB)
	tx, err := db.StartTransaction()
	if err != nil {
		return nil, err
	}

	return service.InitalizeNewSavePortsServiceInstance(
		ctx,
		savePortsToDBChann,
		tx,
		tx.Commit,
		tx.Rollback,
		spsif.log,
	), nil

}
