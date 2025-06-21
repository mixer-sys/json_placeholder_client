package main

import (
	"os"
	"os/signal"
	"time"

	"json_placeholder_client/internal/app/client"
	"json_placeholder_client/internal/app/logger"

	"json_placeholder_client/internal/app/server"
)

func main() {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	log := logger.GetLogger()
	go server.Server()
	time.Sleep(2 * time.Second)
	log.Info.Println("Server is running. Press Ctrl+C to stop.")
	log.Info.Println("Starting JSON Placeholder Client...")

	client.RunTestClient()

	<-sigChan
	log.Info.Println("Received interrupt signal, shutting down...")

}
