package client

import (
	"context"
	"fmt"
	"net/http"

	"json_placeholder_client/internal/app/config"
	handlers "json_placeholder_client/internal/app/handlers"
	"json_placeholder_client/internal/app/logger"
)

func GetClient(cfg *config.Config) (client *http.Client, err error) {
	transport := &http.Transport{
		Proxy: http.ProxyURL(cfg.ProxyURL),
	}

	client = &http.Client{Transport: transport, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return nil
	},
	}
	return client, nil
}

func RunTestClient(ctx context.Context) error {
	cfg, err := config.GetConfig()
	log := logger.GetLogger(cfg)
	if err != nil {
		return fmt.Errorf("error getting config: %v", err)
	}
	client, err := GetClient(cfg)
	if err != nil {
		return fmt.Errorf("error creating HTTP client: %v", err)

	}

	posts, err := handlers.GetPosts(cfg, ctx, client)
	if err != nil {
		return fmt.Errorf("error GetPosts: %v", err)
	}
	log.Info("Retrieved posts successfully. Total posts:", len(posts))

	postId := 1
	post, err := handlers.GetPostByID(cfg, ctx, client, postId)
	if err != nil {
		return fmt.Errorf("error GetPostByID: %v", err)
	}
	log.Info("Got ID: %d, Title: %s, Body: %s", post.ID, post.Title, post.Body)

	_, err = handlers.CreatePost(cfg, ctx, client, post)
	if err != nil {
		return fmt.Errorf("error CreatePost: %v", err)
	}

	log.Info("Created Post ID: %d, Title: %s, Body: %s", post.ID, post.Title, post.Body)

	updated_post, err := handlers.UpdatePost(cfg, ctx, client, postId, post)
	if err != nil {
		return fmt.Errorf("error UpdatePost: %v", err)
	}
	log.Info("Updated Post ID: %d, Title: %s, Body: %s", updated_post.ID, updated_post.Title, updated_post.Body)

	handlers.DeletePost(cfg, ctx, client, postId)
	log.Info("Test client operations completed successfully.")
	return nil
}
