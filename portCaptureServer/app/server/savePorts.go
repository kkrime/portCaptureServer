package server

import (
	"context"
	"io"
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

	// TODO double check this
	savePortServiceInstance, err := s.savePortsService.NewSavePortsInstance(
		context.Background(),
		sqlService.SQLTransactionDB)
	if err != nil {
		// TODO: forward err.Error() to Slack channel #HowTheHellCouldThisHavePossiblyHappend
		return
	}

	for {
		// read in ports from request
		var port *pb.Port
		port, err = portsStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}

		portEntity := convertPBPortToEntityPort(port)

		err = savePortServiceInstance.SavePort(portEntity)
		if err != nil {
			return
		}
	}

	err = savePortServiceInstance.Finalize()
	if err != nil {
		return
	}

	return
}
