package service

import (
	"context"
	"io"
	"portCaptureServer/app/api/pb"
	"portCaptureServer/app/entity"
	"portCaptureServer/app/repository"
	"sync"
)

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

func NewSavePortsService(savePortsRepository repository.SavePortsRepository) SavePortsService {
	savePortsToDBChann := make(chan *savePortToDBParam)

	for i := 0; i < 50; i++ {
		// for i := 0; i < 2; i++ {
		go func() {
			for portsDBParams := range savePortsToDBChann {

				ctx := portsDBParams.ctx
				transaction := portsDBParams.transaction
				port := portsDBParams.port
				resultChann := portsDBParams.resultChann

				resultChann <- savePortsRepository.SavePort(ctx, transaction, port)

				portsDBParams.wg.Done()
			}
		}()
	}

	return &savePortsService{
		savePortsToDBChann:  savePortsToDBChann,
		savePortsRepository: savePortsRepository,
	}
}

func (spp *savePortsService) SavePorts(ctx context.Context, portStream PortsStream) error {
	var wg sync.WaitGroup
	resultChann := make(chan error)
	errorChann := make(chan error)

	ctx, cancelCtx := context.WithCancel(ctx)

	transactoin, err := spp.savePortsRepository.StartTransaction()
	if err != nil {
		return err
	}

	go func() {
		for err := range resultChann {
			if err != nil {
				// cancel context to stop any db io
				cancelCtx()
				errorChann <- err
				return
			}
		}
		close(errorChann)
	}()

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

	wg.Wait()
	close(resultChann)

	for dbSaveError := range errorChann {

		if dbSaveError != nil {
			// roll back the transaction
			transactoin.Rollback()

			// empty resultChann
			for len(resultChann) > 0 {
				<-resultChann
			}

			// empty errorChann
			for len(errorChann) > 0 {
				<-errorChann
			}

			return dbSaveError
		}
	}

	return transactoin.Commit()
}

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
