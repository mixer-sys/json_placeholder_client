package main

import (
	"github.com/mixer-sys/json_placeholder_client/internal/app/logger"
	"github.com/mixer-sys/json_placeholder_client/internal/app/server"
)

func main() {

	log := logger.GetLogger()
	log.Info.Println("Starting JSON Placeholder Client...")
	client, err := client.GetClient()
	if err != nil {
		log.Error.Println("Failed to create HTTP client:", err)
		return
	}
	posts := handlers.GetPosts(client)
	for _, post := range posts {
		log.Info.Printf("Post ID: %d, Title: %s", post.ID, post.Title)
	}
	log.Info.Println("Fetching posts completed successfully.")
	server.Server()
	log.Info.Println("Server is running. Press Ctrl+C to stop.")

}
