package service

import (
	"context"
	"portCaptureServer/app/entity"
	"portCaptureServer/app/repository"
	"sync"

	"github.com/sirupsen/logrus"
)

// savePortToDBParam is pased to the main worker threads
// to write to the db
type savePortToDBParam struct {
	ctx        context.Context
	wg         *sync.WaitGroup
	db         repository.SavePortsRepository
	port       *entity.Port
	errorChann chan<- error
}

type SavePortsService struct {
	savePortsToDBChann              chan<- *savePortToDBParam
	savePortsServiceInstanceFactory SavePortsServiceInstanceFactory
	log                             *logrus.Logger
}

func NewSavePortsService(
	savePortsServiceInstanceFactory SavePortsServiceInstanceFactory,
	numberOfWorkerThreads int,
	log *logrus.Logger) *SavePortsService {
	savePortsToDBChann := make(chan *savePortToDBParam)

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

	return &SavePortsService{
		savePortsToDBChann:              savePortsToDBChann,
		savePortsServiceInstanceFactory: savePortsServiceInstanceFactory,
		log:                             log,
	}
}

func (sps *SavePortsService) NewSavePortsInstance(ctx context.Context) (SavePortsServiceInstance, error) {
	return sps.savePortsServiceInstanceFactory.NewSavePortsInstance(ctx, sps.savePortsToDBChann)
}
