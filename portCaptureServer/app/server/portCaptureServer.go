package server

import (
	"portCaptureServer/app/api/pb"
	"portCaptureServer/app/service"
)

type PortCaptureServer struct {
	pb.UnimplementedPortCaptureServiceServer
	savePortsService *service.SavePortsService
}

func NewPortCaptureServer(savePortsService *service.SavePortsService) *PortCaptureServer {
	return &PortCaptureServer{
		savePortsService: savePortsService,
	}
}
