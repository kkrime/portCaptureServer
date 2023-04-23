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
	ctx        context.Context
	wg         *sync.WaitGroup
	db         repository.SavePortsRepository
	port       *entity.Port
	errorChann chan<- error
}

func NewSavePortToDBParam(
	ctx context.Context,
	wg *sync.WaitGroup,
	db repository.SavePortsRepository,
	errorChann chan<- error,
) SavePortToDBParam {

	return SavePortToDBParam{
		ctx:        ctx,
		wg:         wg,
		db:         db,
		errorChann: errorChann,
	}

}

type savePortsService struct {
	savePortsToDBChann                 chan<- *SavePortToDBParam
	savePortsServiceInstanceFactoryMap map[SavePortsInstanceType]SavePortsServiceInstanceFactory
	log                                *logrus.Logger
}

func NewSavePortsService(
	savePortsServiceInstanceFactoryMap map[SavePortsInstanceType]SavePortsServiceInstanceFactory,
	numberOfWorkerThreads int,
	log *logrus.Logger) *savePortsService {
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
					resultChann := savePortsDBParams.errorChann

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

	return &savePortsService{
		savePortsToDBChann:                 savePortsToDBChann,
		savePortsServiceInstanceFactoryMap: savePortsServiceInstanceFactoryMap,
		log:                                log,
	}
}

// TODO factories should not be pointers
func (sps *savePortsService) NewSavePortsInstance(
	ctx context.Context,
	savePortsInstanceType SavePortsInstanceType) (SavePortsServiceInstance, error) {
	savePortsInstanceFactory := sps.savePortsServiceInstanceFactoryMap[savePortsInstanceType]
	if savePortsInstanceFactory == nil {
		return nil, fmt.Errorf(canNotFindSavePortsInstanceFactoryError, savePortsInstanceType)
	}
	return savePortsInstanceFactory.NewSavePortsInstance(ctx, sps.savePortsToDBChann)
}
