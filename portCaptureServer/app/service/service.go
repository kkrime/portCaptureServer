package service

import (
	"context"
	"portCaptureServer/app/entity"
)

type SavePortsInstanceType string

type SavePortsService interface {
	NewSavePortsInstance(ctx context.Context, savePortsInstanceType SavePortsInstanceType) (SavePortsServiceInstance, error)
}

type SavePortsServiceInstanceFactory interface {
	NewSavePortsInstance(ctx context.Context, savePortsToDBChann chan<- *SavePortToDBParam) (SavePortsServiceInstance, error)
}

type SavePortsServiceInstance interface {
	SavePort(port *entity.Port) error
	Finalize() error
}
