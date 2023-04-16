package service

import (
	"context"
	"io"
	"portCaptureServer/app/api/pb"
	"portCaptureServer/app/entity"
	"portCaptureServer/app/repository"
	"sync"
)

// savePortToDBParam is pased to the main worker threads
// to write to the db
type savePortToDBParam struct {
	ctx         context.Context
	wg          *sync.WaitGroup
	transaction repository.Transaction
	port        *entity.Port
	resultChann chan<- error
}

type savePortsService struct {
	savePortsToDBChann  chan<- *savePortToDBParam
	savePortsRepository repository.SavePortsRepository
}

func NewSavePortsService(savePortsRepository repository.SavePortsRepository, numberOfWorkerThreads int) SavePortsService {
	savePortsToDBChann := make(chan *savePortToDBParam)

	// spawn the main Worker Threads, these are the threads that save the ports to the database
	for i := 0; i < numberOfWorkerThreads; i++ {
		// WORKER THREAD
		go func() {
			for savePortsDBParams := range savePortsToDBChann {
				// using an anonymous function here, so that portsDBParams.wg.Done()
				// gets called no matter what and there are no hanging go routines
				// (in the unlikley event of a panic)
				func() {
					defer savePortsDBParams.wg.Done()

					ctx := savePortsDBParams.ctx
					transaction := savePortsDBParams.transaction
					port := savePortsDBParams.port
					resultChann := savePortsDBParams.resultChann

					resultChann <- savePortsRepository.SavePort(ctx, transaction, port)
				}()
			}
		}()
	}

	return &savePortsService{
		savePortsToDBChann:  savePortsToDBChann,
		savePortsRepository: savePortsRepository,
	}
}

// SavePort saves incoming ports (from gRPC stream) to the database via the worker threads
// it works on a all or nothing basis; if there are any errors saving any of the ports,
// then no ports will be saved to the db
func (spp *savePortsService) SavePorts(ctx context.Context, portStream PortsStream) error {
	var wg sync.WaitGroup
	resultChann := make(chan error)
	errorChann := make(chan error, 1)

	// start transaction
	transactoin, err := spp.savePortsRepository.StartTransaction()
	if err != nil {
		return err
	}

	ctx, cancelCtx := context.WithCancel(ctx)

	// this go routines manages the results from the worker threads above
	go func() {
		// only send the first error
		firstEroorOccured := false
		for err := range resultChann {
			if err != nil && !firstEroorOccured {
				firstEroorOccured = true
				// cancel context to stop any db io
				cancelCtx()
				// if any errors on saving the ports to the db, send it to
				// errorChann
				errorChann <- err
			}
		}
		close(errorChann)
	}()

	// this loop reads the ports from the netowrk (gRPC stream) and
	// the ports to the worker thread
	for {
		port, err := portStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			// cancel context to stop any db io
			cancelCtx()
			errorChann <- err
			break
		}

		portEntity := convertPBPortToEntityPort(port)

		wg.Add(1)
		spp.savePortsToDBChann <- &savePortToDBParam{
			ctx:         ctx,
			wg:          &wg,
			transaction: transactoin,
			port:        portEntity,
			resultChann: resultChann,
		}
	}

	// wait for all the ports to be written to the db
	wg.Wait()
	close(resultChann)

	for dbSaveError := range errorChann {

		if dbSaveError != nil {
			// roll back the transaction
			transactoin.Rollback()

			// drain the errorChann
			for len(errorChann) > 0 {
				<-errorChann
			}

			return dbSaveError
		}
	}

	// commit changes to the db
	return transactoin.Commit()
}

// convertPBPortToEntityPort converts the pb (protobuf) format for ports
// to the entity format
func convertPBPortToEntityPort(port *pb.Port) *entity.Port {
	return &entity.Port{
		Name:         port.Name,
		PrimaryUnloc: port.PrimaryUnloc,
		Code:         port.Code,
		City:         port.City,
		Country:      port.Country,
		Alias: func() *[]entity.Alias {
			alias := make([]entity.Alias, 0, len(port.Alias))
			for _, a := range port.Alias {
				alias = append(alias, entity.Alias{
					Name: a,
				})
			}
			return &alias
		}(),
		Regions: func() *[]entity.Region {
			regions := make([]entity.Region, 0, len(port.Regions))
			for _, r := range port.Regions {
				regions = append(regions, entity.Region{
					Name: r,
				})
			}
			return &regions
		}(),
		Coordinantes: func() [2]float32 {
			if len(port.Coordinates) == 2 {
				return [2]float32{port.Coordinates[0], port.Coordinates[1]}
			}
			return [2]float32{-1, -1}
		}(),
		Province: port.Province,
		Timezone: port.Timezone,
		Unlocs: func() *[]entity.Unloc {
			unlocs := make([]entity.Unloc, 0, len(port.Unlocs))
			for _, u := range port.Unlocs {
				unlocs = append(unlocs, entity.Unloc{
					Name: u,
				})
			}
			return &unlocs
		}(),
	}
}
