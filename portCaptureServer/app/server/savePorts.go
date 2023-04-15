package server

import (
	"context"
	"portCaptureServer/app/api/pb"
)

func (s *Server) SavePorts(portsStream pb.PortCaptureService_SavePortsServer) error {
	response := pb.PortCaptureServiceResponse{}

	err := s.savePortsService.SavePorts(context.Background(), portsStream)
	if err != nil {
		response.Error = err.Error()
		portsStream.SendAndClose(&response)
		return err
	}

	response.Success = true
	portsStream.SendAndClose(&response)

	return nil
}
