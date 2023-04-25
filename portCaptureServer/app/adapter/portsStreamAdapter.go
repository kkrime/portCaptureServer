package adapter

import (
	"portCaptureServer/app/api/pb"
	"portCaptureServer/app/entity"
)

type portsStreamAdapter struct {
	portsStream pb.PortCaptureService_SavePortsServer
}

func NewPortsStreamAdapter(
	portsStream pb.PortCaptureService_SavePortsServer,
) PortsStream {
	return &portsStreamAdapter{
		portsStream: portsStream,
	}
}

func (ps *portsStreamAdapter) Recv() (*entity.Port, error) {

	port, err := ps.portsStream.Recv()
	if err != nil {
		return nil, err
	}

	return convertPBPortToEntityPort(port), nil
}
