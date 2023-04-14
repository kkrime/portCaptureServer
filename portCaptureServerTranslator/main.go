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
	/*
		jsonFile, err := os.ReadFile("/home/aa/Desktop/90poe/golang/ports.json")
		if err != nil {
			fmt.Println(err)
		}

		var users1 map[string]portCaptureServerPb.Port
		// var users1 portCaptureServerPb.Ports
		err = json.Unmarshal([]byte(jsonFile), &users1)
		fmt.Printf("err = %+v\n", err)
		fmt.Printf(" &users = %+v\n dddd", users1)
	*/

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
