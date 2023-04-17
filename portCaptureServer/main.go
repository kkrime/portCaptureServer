package main

import (
	"portCaptureServer/app"
	"portCaptureServer/app/logger"
)

func run() error {
	app, err := app.NewApp()
	if err != nil {
		return err
	}

	return app.Run()
}

func main() {
	if err := run(); err != nil {
		logger.CreateNewLogger().Fatal(err)
	}
}
