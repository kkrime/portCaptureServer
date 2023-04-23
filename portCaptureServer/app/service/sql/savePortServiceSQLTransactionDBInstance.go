package sql

import (
	"context"
	"portCaptureServer/app/repository"
	"portCaptureServer/app/repository/sql"
	"portCaptureServer/app/service"
	"sync"

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
	var wg sync.WaitGroup
	errorChann := make(chan error, 1)

	// TODO compile time check
	db := spsif.savePortsRepository.(sql.DB)
	tx, err := db.StartTransaction()
	if err != nil {
		return nil, err
	}

	ctx, cancelCtx := context.WithCancel(ctx)

	// return &service.savePortsServiceInstance{
	return service.NewSavePortsServiceInstance(
		savePortsToDBChann,
		service.NewSavePortToDBParam(
			ctx,
			&wg,
			db,
			errorChann,
		),
		&wg,
		cancelCtx,
		errorChann,
		tx.Commit,
		tx.Rollback,
		spsif.log,
	), nil
}
