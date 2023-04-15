package server

import (
	"portCaptureServer/app/api/pb"
	"portCaptureServer/app/service"
)

type Server struct {
	pb.UnimplementedPortCaptureServiceServer
	savePortsService service.SavePortsService
}

func NewServer(savePortsService service.SavePortsService) *Server {
	return &Server{
		savePortsService: savePortsService,
	}
}
