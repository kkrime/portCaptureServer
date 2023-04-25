package server

import (
	"context"
	"portCaptureServer/app/adapter"
	"portCaptureServer/app/api/pb"
	sqlService "portCaptureServer/app/service/sql"
)

func (s *PortCaptureServer) SavePorts(portsStream pb.PortCaptureService_SavePortsServer) error {
	response := pb.PortCaptureServiceResponse{}

	savePortServiceInstance, err := s.savePortsServiceProvider.NewSavePortsInstance(
		context.Background(),
		sqlService.SQLTransactionDB)

	if err != nil {
		// TODO: forward err.Error() to Slack channel #HowTheHellCouldThisHavePossiblyHappend
		return err
	}

	portsStreamAdapter := adapter.NewPortsStreamAdapter(portsStream)
	if err != nil {
		response.Error = err.Error()
		return portsStream.SendAndClose(&response)
	}

	err = savePortServiceInstance.SavePort(portsStreamAdapter)
	if err != nil {
		response.Error = err.Error()
		return portsStream.SendAndClose(&response)
	}

	response.Success = true
	return portsStream.SendAndClose(&response)
}
