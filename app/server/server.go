package server

import "portCaptureServer/app/server/api/pb"

type Server struct {
	pb.UnimplementedPortCaptureServiceServer
}

func NewServer() *Server {
	return &Server{}
}
