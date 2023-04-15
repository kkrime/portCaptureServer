package app

import (
	"fmt"
	"log"
	"net"
	"portCaptureServer/app/api/pb"
	"portCaptureServer/app/config"
	"portCaptureServer/app/repository"
	"portCaptureServer/app/server"
	"portCaptureServer/app/service"

	"github.com/jmoiron/sqlx"
	"gitlab.com/avarf/getenvs"
	"google.golang.org/grpc"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

type App interface {
	Run() error
}

type app struct {
	server *server.Server
}

func NewApp() (App, error) {
	app := &app{}
	configFilePath := getenvs.GetEnvString("CONFIG_FILE_PATH", "./config/local_config.toml")

	config, err := config.ReadConfig(configFilePath)
	fmt.Printf("config = %+v\n", config)
	if err != nil {
		return nil, err
	}

	db, err := app.ConnectToDB(config.DBConfig)
	if err != nil {
		return nil, err
	}

	SavePortsRepository := repository.NewSavePortsRepository(db)
	SavePortService := service.NewSavePortsService(SavePortsRepository)
	app.server = server.NewServer(SavePortService)
	return app, nil
}

func (a *app) ConnectToDB(config config.DBConfig) (*sqlx.DB, error) {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Dbname)

	database, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		return nil, err
	}

	return database, nil
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
