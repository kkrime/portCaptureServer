package service

import (
	"context"
	"portCaptureServer/app/repository"
	"portCaptureServer/app/repository/sql"
	"sync"

	"github.com/sirupsen/logrus"
)

type savePortsServiceInstanceFactory struct {
	savePortsRepository repository.SavePortsRepository
	log                 *logrus.Logger
}

func NewSavePortsServiceInstanceFactory(
	savePortsRepository repository.SavePortsRepository,
	log *logrus.Logger,
) SavePortsServiceInstanceFactory {
	return &savePortsServiceInstanceFactory{
		savePortsRepository: savePortsRepository,
		log:                 log,
	}
}

func (spsif *savePortsServiceInstanceFactory) NewSavePortsInstance(
	ctx context.Context,
	savePortsToDBChann chan<- *savePortToDBParam,
) (SavePortsServiceInstance, error) {
	var wg sync.WaitGroup
	errorChann := make(chan error, 1)

	// TODO compile time check
	db := spsif.savePortsRepository.(sql.DB)
	tx, err := db.StartTransaction()
	if err != nil {
		return nil, err
	}

	ctx, cancelCtx := context.WithCancel(ctx)

	return &savePortsServiceInstance{
		savePortsToDBChann: savePortsToDBChann,
		savePortToDBParam: savePortToDBParam{
			ctx:        ctx,
			wg:         &wg,
			db:         tx,
			errorChann: errorChann,
		},
		cancelCtx:  cancelCtx,
		wg:         &wg,
		errorChann: errorChann,
		finalizeDB: tx.Commit,
		cancelDB:   tx.Rollback,
		log:        spsif.log,
	}, nil
}
