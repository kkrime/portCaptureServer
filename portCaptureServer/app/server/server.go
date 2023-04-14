package server

import "portCaptureServer/app/api/pb"

type Server struct {
	pb.UnimplementedPortCaptureServiceServer
}

func NewServer() *Server {
	return &Server{}
}
