package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	portCaptureServerPb "portCaptureServer/app/api/pb"
	"portCaptureServerTranslator/app/controller"
	"portCaptureServerTranslator/app/service"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	requestid "github.com/sumit-tembe/gin-requestid"
	"gitlab.com/avarf/getenvs"
	"google.golang.org/grpc"
)

type App interface {
	Run() error
}

type app struct {
}

func NewApp() (App, error) {
	app := &app{}

	return app, nil
}

func (a *app) createPortCatpureServerClient() (portCaptureServerPb.PortCaptureServiceClient, error) {
	portCaptureServerAdddress := getenvs.GetEnvString("PORT_CAPTURE_SERVER_ADDRESS", "localhost:20000")
	conn, err := grpc.Dial(portCaptureServerAdddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
		return nil, err
	}

	return portCaptureServerPb.NewPortCaptureServiceClient(conn), nil

}

func (a *app) Run() error {

	return a.startWebServer()
}

func (a *app) startWebServer() error {
	portCaptureServerClient, err := a.createPortCatpureServerClient()
	if err != nil {
		return err
	}
	sendPortService := service.NewSendPortService(portCaptureServerClient)
	sendPortsController := controller.NewSendPortController(sendPortService)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.LoggerWithConfig(requestid.GetLoggerConfig(nil, nil, nil)))

	router.GET("healthcheck", sendPortsController.HealthCheck)

	v1 := router.Group("v1")

	sendPorts := v1.Group("sendports")

	sendPorts.POST("", sendPortsController.SendPorts)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// graceful exit
	shutdownErrChann := make(chan error, 1)
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			shutdownErrChann <- fmt.Errorf("Server Shutdown: %w", err)
			return
		}
		// catching ctx.Done(). timeout of 5 seconds.
		select {
		case <-ctx.Done():
			log.Println("timeout of 5 seconds.")
		}
		close(shutdownErrChann)
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return <-shutdownErrChann
}
