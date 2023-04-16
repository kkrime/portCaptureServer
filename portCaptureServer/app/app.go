package app

import (
	"log"
	"net"
	"os"
	"os/signal"
	"portCaptureServer/app/api/pb"
	"portCaptureServer/app/config"
	"portCaptureServer/app/repository"
	"portCaptureServer/app/server"
	"portCaptureServer/app/service"
	"syscall"
	"time"

	"gitlab.com/avarf/getenvs"
	"google.golang.org/grpc"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type App interface {
	Run() error
}

type app struct {
	portCaptureServer *server.PortCaptureServer
}

func NewApp() (App, error) {
	app := &app{}
	configFilePath := getenvs.GetEnvString("CONFIG_FILE_PATH", "./config/local_config.toml")

	config, err := config.ReadConfig(configFilePath)
	if err != nil {
		return nil, err
	}

	db, err := app.ConnectToDB(config.DBConfig)
	if err != nil {
		return nil, err
	}

	numberOfWorkerThreads := config.PortCapture.WorkerThreads

	SavePortsRepository := repository.NewSavePortsRepository(db)
	SavePortService := service.NewSavePortsService(SavePortsRepository, numberOfWorkerThreads)
	app.portCaptureServer = server.NewPortCaptureServer(SavePortService)
	return app, nil
}

func (a *app) Run() error {
	listner, err := net.Listen("tcp", ":20000")
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPortCaptureServiceServer(grpcServer, a.portCaptureServer)
	log.Printf("server listening on %v", listner.Addr())

	// graceful shut down
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-sigCh
		log.Printf("got signal %v, attempting to GRACEFULLY shutdown", s)
		gracefulShutdownChann := make(chan struct{})
		go func() {
			grpcServer.GracefulStop()
			gracefulShutdownChann <- struct{}{}
		}()

		// create a timeout of 5 seconds
		gracefulShutdownTimeoutChann := time.NewTimer(time.Second * 5)
		select {
		case <-gracefulShutdownChann:
			log.Printf("Graceful shutdown complete")
		case <-gracefulShutdownTimeoutChann.C:
			log.Printf("Graceful shutdown timed out, shutting down forefully")
			grpcServer.Stop()
		}
	}()

	// start the gRPC server
	return grpcServer.Serve(listner)
}
