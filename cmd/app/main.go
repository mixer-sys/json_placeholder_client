package main

import (
	"github.com/mixer-sys/json_placeholder_client/internal/app/logger"
	"github.com/mixer-sys/json_placeholder_client/internal/app/server"
)

func main() {

	log := logger.GetLogger()
	go server.Server()

	log.Info.Println("Server is running. Press Ctrl+C to stop.")

	log.Info.Println("Starting JSON Placeholder Client...")

}
