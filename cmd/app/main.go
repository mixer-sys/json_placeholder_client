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
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error getting config: %w", err)
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	log := logger.GetLogger(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("server.Run end")
				return
			default:
				err := server.Run(ctx, cfg)
				if err != nil {
					log.Error("Error starting server: %w", err)
					return
				}
			}
		}
	}()
	log.Info("Server is running. Press Ctrl+C to stop.")
	log.Info("Starting JSON Placeholder Client...")

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("RunTestClient end")
				return
			default:
				err := client.RunTestClient(ctx)
				if err != nil {
					log.Error("Error", err)
					return
				}
			}
		}
	}()
	<-sigChan
	log.Info("Received interrupt signal, shutting down...")
	cancel()
	<-ctx.Done()

	log.Info("Shutdown complete.")
}
