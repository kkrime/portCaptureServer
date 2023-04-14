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
		fmt.Printf("err = %+v\n", err)
		return err
	}

	for _, port := range ports {
		err = stream.Send(port)
		if err != nil {
			return err
		}
	}

	return nil
}
