package service

import (
	"context"
	"encoding/json"
	"fmt"
	portCaptureServerPb "portCaptureServer/app/api/pb"
)

type sendPortService struct {
	portCaptureServerClient portCaptureServerPb.PortCaptureServiceClient
}

func NewSendPortService(portCaptureServerClient portCaptureServerPb.PortCaptureServiceClient) SendPortService {
	return &sendPortService{
		portCaptureServerClient: portCaptureServerClient,
	}
}

func (sps *sendPortService) SendPort(ctx context.Context, portData *[]byte) error {
	var ports map[string]*portCaptureServerPb.Port

	err := json.Unmarshal(*portData, &ports)
	if err != nil {
		return err
	}

	stream, err := sps.portCaptureServerClient.SavePorts(ctx)
	if err != nil {
		return err
	}

	for uloc, port := range ports {
		port.PrimaryUnloc = uloc
		err = stream.Send(port)
		if err != nil {
			break
		}

		// for really large files
		// this will help free memeory
		// via the garbage collector
		ports[uloc] = nil
	}

	response, err := stream.CloseAndRecv()

	if err != nil {
		return err
	}

	if response.Success == false {
		return fmt.Errorf(response.Error)
	}

	return nil
}
