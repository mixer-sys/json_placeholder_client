package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"json_placeholder_client/internal/app/client"
	"json_placeholder_client/internal/app/config"
	"json_placeholder_client/internal/app/logger"
	"json_placeholder_client/internal/app/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error getting config: %w", err)
		os.Exit(1)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	log := logger.GetLogger(cfg)
	ctx, cancel := context.WithCancel(context.Background())

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
					os.Exit(1)
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
					os.Exit(1)
				}
			}
		}
	}()
	<-sigChan
	log.Info("Received interrupt signal, shutting down...")
	cancel()
	<-ctx.Done()
	time.Sleep(1 * time.Second)
	log.Info("Shutdown complete.")
}
