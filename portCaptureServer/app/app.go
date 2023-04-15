package app

import (
	"log"
	"net"
	"portCaptureServer/app/api/pb"
	"portCaptureServer/app/server"
	"portCaptureServer/app/service"

	"google.golang.org/grpc"
)

type App interface {
	Run() error
}

type app struct {
	server *server.Server
}

func NewApp() App {
	app := &app{}
	SavePortService := service.NewSavePortsService(nil)
	app.server = server.NewServer(SavePortService)
	return app
}

func (a *app) Run() error {
	listner, err := net.Listen("tcp", ":20000")
	if err != nil {
		// log.Errorf("failed to listen: %v", err)
		return err
	}
	grpcServer := grpc.NewServer()
	pb.RegisterPortCaptureServiceServer(grpcServer, a.server)
	log.Printf("server listening on %v", listner.Addr())
	return grpcServer.Serve(listner)
}
