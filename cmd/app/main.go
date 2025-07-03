package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"json_placeholder_client/internal/app/client"
	"json_placeholder_client/internal/app/config"
	"json_placeholder_client/internal/app/logger"
	"json_placeholder_client/internal/app/server"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("Error getting config: %v", err)
		os.Exit(1)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	log := logger.GetLogger(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err := server.Run(cfg, ctx)
		if err != nil {
			log.Error("Error starting server: %v", err)
			os.Exit(1)
		}
	}()
	log.Info("Server is running. Press Ctrl+C to stop.")
	log.Info("Starting JSON Placeholder Client...")

	defer cancel()
	go func() {
		err := client.RunTestClient(ctx)
		if err != nil {
			log.Error("Error", err)
			os.Exit(1)
		}

	}()
	<-sigChan
	log.Info("Received interrupt signal, shutting down...")
	cancel()
	<-ctx.Done()
	log.Info("Shutdown complete.")
}
