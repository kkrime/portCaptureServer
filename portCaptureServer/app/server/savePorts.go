package server

import (
	"portCaptureServer/app/adapter"
	"portCaptureServer/app/api/pb"
	sqlService "portCaptureServer/app/service/sql"
)

func (s *PortCaptureServer) SavePorts(portsStream pb.PortCaptureService_SavePortsServer) (err error) {
	response := pb.PortCaptureServiceResponse{}
	defer func() {
		if err != nil {
			response.Error = err.Error()
		} else {
			response.Success = true
		}
		err = portsStream.SendAndClose(&response)
	}()

	savePortServiceInstance, err := s.savePortsServiceProvider.NewSavePortsInstance(
		s.masterCtx,
		sqlService.SQLTransactionDB)

	if err != nil {
		// TODO: forward err.Error() to Slack channel #HowTheHellCouldThisHavePossiblyHappend
		return
	}

	portsStreamAdapter := adapter.NewPortsStreamAdapter(portsStream)

	err = savePortServiceInstance.SavePorts(portsStreamAdapter)

	return
}
