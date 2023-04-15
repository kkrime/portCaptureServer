package app

import (
	"fmt"
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

	"github.com/jmoiron/sqlx"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/logrusadapter"
	"github.com/sirupsen/logrus"
	"gitlab.com/avarf/getenvs"
	"google.golang.org/grpc"

	// _ "github.com/bmizerany/pq"
	// _ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	// _ "github.com/lib/pq"
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

	dbLog := CreateNewLogger()

	database.DB = sqldblogger.OpenDriver(dsn, database.DB.Driver(), logrusadapter.New(dbLog),
		sqldblogger.WithTimeFormat(sqldblogger.TimeFormatRFC3339),
		sqldblogger.WithLogDriverErrorSkip(true),
		sqldblogger.WithSQLQueryAsMessage(true))

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

	// graceful shut down
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-sigCh
		log.Printf("got signal %v, attempting graceful shutdown", s)
		grpcServer.GracefulStop()
		log.Printf("Graceful shutdown complete")
	}()

	return grpcServer.Serve(listner)
}

func CreateNewLogger() *logrus.Logger {
	return &logrus.Logger{
		Out:   os.Stdout,
		Level: logrus.DebugLevel,
		Formatter: &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			PadLevelText:    true,
			ForceColors:     true,
		},
	}
}
