package server

import (
	"portCaptureServer/app/api/pb"
	"portCaptureServer/app/service"
)

type PortCaptureServer struct {
	pb.UnimplementedPortCaptureServiceServer
	savePortsServiceProvider service.SavePortsServiceProvider
}

func NewPortCaptureServer(savePortsServiceProvider service.SavePortsServiceProvider) *PortCaptureServer {
	return &PortCaptureServer{
		savePortsServiceProvider: savePortsServiceProvider,
	}
}
