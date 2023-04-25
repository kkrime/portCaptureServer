package adapter

import "portCaptureServer/app/entity"

type PortsStream interface {
	Recv() (*entity.Port, error)
}
