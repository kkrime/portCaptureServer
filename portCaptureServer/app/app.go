package app

import (
	"context"
	"net"
	"os"
	"os/signal"
	"portCaptureServer/app/api/pb"
	"portCaptureServer/app/config"
	"portCaptureServer/app/logger"
	sqlRepository "portCaptureServer/app/repository/sql"
	"portCaptureServer/app/server"
	"portCaptureServer/app/service"
	sqlService "portCaptureServer/app/service/sql"
	"syscall"

	"gitlab.com/avarf/getenvs"
	"google.golang.org/grpc"
)

type App interface {
	Run() error
}

type app struct {
	portCaptureServer *server.PortCaptureServer
	masterContext     context.Context
	ctxCancel         context.CancelFunc
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

	log := logger.CreateNewLogger()

	numberOfWorkerThreads := config.PortCapture.WorkerThreads

	savePortsRepository := sqlRepository.NewSQLDB(db)
	savePortsServiceInstanceFactory := sqlService.NewSavePortsServiceSQLTransactionInstanceFactory(savePortsRepository, log)

	savePortsServiceInstanceFactoryMap := map[service.SavePortsInstanceType]service.SavePortsServiceInstanceFactory{
		sqlService.SQLTransactionDB: savePortsServiceInstanceFactory,
		// sqlService.SQLDB:            sqlService.NewSavePortsServiceSQLInstanceFactory(savePortsRepository, log),
	}

	savePortServiceProvider := service.NewSavePortsServiceProvider(
		savePortsServiceInstanceFactoryMap,
		numberOfWorkerThreads,
		log)

	masterCtx, ctxCancel := context.WithCancel(context.Background())
	app.masterContext = masterCtx
	app.ctxCancel = ctxCancel

	app.portCaptureServer = server.NewPortCaptureServer(savePortServiceProvider, masterCtx)

	return app, nil
}

func (a *app) Run() error {
	log := logger.CreateNewLogger()
	listner, err := net.Listen("tcp", ":20000")
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPortCaptureServiceServer(grpcServer, a.portCaptureServer)
	log.Infof("server listening on %v", listner.Addr())

	// graceful shut down
	ctx, _ := signal.NotifyContext(a.masterContext, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		<-ctx.Done()
		a.ctxCancel()
		grpcServer.GracefulStop()
	}()

	// start the gRPC server
	return grpcServer.Serve(listner)
}
