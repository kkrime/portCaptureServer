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
	resultChann := make(chan error)

	// TODO compile time check
	db := spsif.savePortsRepository.(sql.DB)
	tx, err := db.StartTransaction()
	if err != nil {
		return nil, err
	}

	ctx, cancelCtx := context.WithCancel(ctx)

	go func() {
		// only send the first error
		firstErrorOccured := false
		for err := range resultChann {
			if !firstErrorOccured {
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
			tx,
			resultChann,
		),
		&wg,
		errorChann,
		resultChann,
		tx.Commit,
		tx.Rollback,
		spsif.log,
	), nil

}
