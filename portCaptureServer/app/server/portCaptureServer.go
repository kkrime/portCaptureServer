package server

import (
	"context"
	"portCaptureServer/app/api/pb"
	"portCaptureServer/app/service"
)

type PortCaptureServer struct {
	pb.UnimplementedPortCaptureServiceServer
	savePortsServiceProvider service.SavePortsServiceProvider
	masterCtx                context.Context
}

func NewPortCaptureServer(
	savePortsServiceProvider service.SavePortsServiceProvider,
	masterCtx context.Context) *PortCaptureServer {
	return &PortCaptureServer{
		savePortsServiceProvider: savePortsServiceProvider,
		masterCtx:                masterCtx,
	}
}
