package service

import (
	"context"
	"portCaptureServer/app/adapter"
)

type SavePortsInstanceType string

type SavePortsServiceProvider interface {
	NewSavePortsInstance(ctx context.Context, savePortsInstanceType SavePortsInstanceType) (SavePortsServiceInstance, error)
}

type SavePortsServiceInstanceFactory interface {
	NewSavePortsInstance(
		ctx context.Context,
		savePortsToDBChann chan<- *SavePortToDBParam) (SavePortsServiceInstance, error)
}

type SavePortsServiceInstance interface {
	SavePort(portsStream adapter.PortsStream) error
}
