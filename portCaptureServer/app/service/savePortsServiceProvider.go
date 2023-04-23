package service

import (
	"context"
	"fmt"
	"portCaptureServer/app/entity"
	"portCaptureServer/app/repository"
	"sync"

	"github.com/sirupsen/logrus"
)

// SavePortToDBParam is pased to the main worker threads
// to write to the db
type SavePortToDBParam struct {
	ctx         context.Context
	wg          *sync.WaitGroup
	db          repository.SavePortsRepository
	port        *entity.Port
	resultChann chan<- error
}

func NewSavePortToDBParam(
	ctx context.Context,
	wg *sync.WaitGroup,
	db repository.SavePortsRepository,
	resultChann chan<- error,
) SavePortToDBParam {

	return SavePortToDBParam{
		ctx:         ctx,
		wg:          wg,
		db:          db,
		resultChann: resultChann,
	}

}

type savePortsServiceProvider struct {
	savePortsToDBChann                 chan<- *SavePortToDBParam
	numberOfWorkerThreads              int
	savePortsServiceInstanceFactoryMap map[SavePortsInstanceType]SavePortsServiceInstanceFactory
	log                                *logrus.Logger
}

func NewSavePortsServiceProvider(
	savePortsServiceInstanceFactoryMap map[SavePortsInstanceType]SavePortsServiceInstanceFactory,
	numberOfWorkerThreads int,
	log *logrus.Logger) *savePortsServiceProvider {
	savePortsToDBChann := make(chan *SavePortToDBParam)

	// spawn the main Worker Threads, these are the threads that save the ports to the database
	for i := 0; i < numberOfWorkerThreads; i++ {
		// WORKER THREAD
		log.Infof("Spwaning Worker Thread #%d", i+1)
		go func() {
			for savePortsDBParams := range savePortsToDBChann {
				// using an anonymous function here, so that portsDBParams.wg.Done()
				// gets called no matter what and there are no hanging go routines
				// (in the unlikley event of a panic)
				func() {
					defer savePortsDBParams.wg.Done()

					ctx := savePortsDBParams.ctx
					port := savePortsDBParams.port
					resultChann := savePortsDBParams.resultChann

					db := savePortsDBParams.db

					log.Infof("Saving Port #%s To The Database", port.PrimaryUnloc)
					err := db.SavePort(ctx, port)

					if err != nil {
						resultChann <- err
					}
				}()
			}
		}()
	}

	return &savePortsServiceProvider{
		savePortsToDBChann:                 savePortsToDBChann,
		numberOfWorkerThreads:              numberOfWorkerThreads,
		savePortsServiceInstanceFactoryMap: savePortsServiceInstanceFactoryMap,
		log:                                log,
	}
}

func (spsp *savePortsServiceProvider) NewSavePortsInstance(
	ctx context.Context,
	savePortsInstanceType SavePortsInstanceType) (SavePortsServiceInstance, error) {
	savePortsInstanceFactory := spsp.savePortsServiceInstanceFactoryMap[savePortsInstanceType]
	if savePortsInstanceFactory == nil {
		return nil, fmt.Errorf(canNotFindSavePortsInstanceFactoryError, savePortsInstanceType)
	}
	return savePortsInstanceFactory.NewSavePortsInstance(ctx, spsp.savePortsToDBChann)
}
