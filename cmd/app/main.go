package main

import (
	"os"
	"os/signal"

	"json_placeholder_client/internal/app/client"
	"json_placeholder_client/internal/app/logger"

	"json_placeholder_client/internal/app/server"
)

func main() {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	log := logger.GetLogger()
	go server.Server()
	log.Info.Println("Server is running. Press Ctrl+C to stop.")
	log.Info.Println("Starting JSON Placeholder Client...")

	err := client.RunTestClient()
	if err != nil {
		log.Error.Println(err)
	}

	<-sigChan
	log.Info.Println("Received interrupt signal, shutting down...")

}
