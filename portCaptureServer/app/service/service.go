package service

import (
	"context"
	"portCaptureServer/app/entity"
)

type SavePortsServiceInstanceFactory interface {
	NewSavePortsInstance(ctx context.Context, savePortsToDBChann chan<- *savePortToDBParam) (SavePortsServiceInstance, error)
}

type SavePortsServiceInstance interface {
	SavePort(port *entity.Port) error
	Finalize() error
}
