package server

import (
	"fmt"
	"portCaptureServer/app/api/pb"
)

func (s *Server) SavePorts(portsStream pb.PortCaptureService_SavePortsServer) error {

	for {
		ports, err := portsStream.Recv()
		if err != nil {
			fmt.Printf("err = %+v\n", err)
			return err
		}

		fmt.Printf("ports = %+v\n", ports)
	}
}
