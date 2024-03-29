package main

import (
	"log"
	"portCaptureServerTranslator/app"
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
		log.Fatal(err)
	}
}
