package service

import (
	"context"
	"portCaptureServer/app/api/pb"
)

// compile time check to make sure pb.PortCaptureService_SavePortsServer can be passed as PortsStream
var _ PortsStream = (pb.PortCaptureService_SavePortsServer)(nil)

// PortsStream exists becuase:
// 1. I do not want to explicitly pass/use any protobuff (pb) specific code in the service layer
// because pb code should only be used in the server layer, but I need to call Recv()
// 2. It will simply testing
type PortsStream interface {
	Recv() (*pb.Port, error)
}

type SavePortsService interface {
	SavePorts(ctx context.Context, portStream PortsStream) error
}
