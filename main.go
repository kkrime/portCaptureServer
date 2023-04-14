package main

import (
	"log"
	"portCaptureServer/app"
)

func main() {
	if err := app.NewApp().Run(); err != nil {
		log.Fatal(err)
	}
}
