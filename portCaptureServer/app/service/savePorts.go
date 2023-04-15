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
	wg          sync.WaitGroup
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
		go func() {
			for portsDBParams := range savePortsToDBChann {
				defer portsDBParams.wg.Done()

				ctx := portsDBParams.ctx
				transaction := portsDBParams.transaction
				port := portsDBParams.port
				resultChann := portsDBParams.resultChann

				resultChann <- savePortsRepository.SavePort(ctx, transaction, port)
			}
		}()
	}

	return &savePortsService{
		savePortsToDBChann: savePortsToDBChann,
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
		for {
			port, err := portStream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				// cancel context to stop any db io
				cancelCtx()
				errorChann <- err
				return
			}

			portEntity := convertPBPortToEntityPort(port)

			wg.Add(1)
			spp.savePortsToDBChann <- &savePortToDBParam{
				ctx:         ctx,
				wg:          wg,
				transaction: transactoin,
				port:        portEntity,
				resultChann: resultChann,
			}
		}
	}()

	go func() {
		for err := range resultChann {
			if err != nil {
				// cancel context to stop any db io
				cancelCtx()
				errorChann <- err
				return
			}
		}

		// closing error chan will result in a read of defulat value of error (nil)
		// close(errorChann)
		// not closing channel, line 68 might cause a race condition
		// sending nil instead
		errorChann <- nil
	}()

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
		Name:    port.Name,
		Code:    port.Code,
		City:    port.City,
		Country: port.Country,
		Alias: func() *[]entity.Alias {
			alias := make([]entity.Alias, len(port.Alias))
			for _, a := range port.Alias {
				alias = append(alias, entity.Alias{
					Name: a,
				})
			}
			return &alias
		}(),
		Regions: func() *[]entity.Region {
			regions := make([]entity.Region, len(port.Regions))
			for _, r := range port.Regions {
				regions = append(regions, entity.Region{
					Name: r,
				})
			}
			return &regions
		}(),
		Province: port.Province,
		Timezone: port.Timezone,
		Unlocs: func() *[]entity.Unloc {
			unlocs := make([]entity.Unloc, len(port.Unlocs))
			for _, u := range port.Unlocs {
				unlocs = append(unlocs, entity.Unloc{
					Name: u,
				})
			}
			return &unlocs
		}(),
	}
}
