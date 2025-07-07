package client

import (
	"context"
	"fmt"

	"json_placeholder_client/internal/app/config"
	handlers "json_placeholder_client/internal/app/handlers"
	"json_placeholder_client/internal/app/logger"
)

func RunTestClient(ctx context.Context) error {
	cfg, err := config.Load()
	log := logger.GetLogger(cfg)
	if err != nil {
		return fmt.Errorf("error getting config: %v", err)
	}
	pc := handlers.NewPostClient(cfg, nil)

	pc.SetClient()

	posts, err := pc.GetPosts(ctx)
	if err != nil {
		return fmt.Errorf("error GetPosts: %v", err)
	}
	log.Info("Retrieved posts successfully. Total posts:", len(posts))

	postId := 1
	post, err := pc.GetPostByID(ctx, postId)
	if err != nil {
		return fmt.Errorf("error GetPostByID: %v", err)
	}
	log.Info("Got ID: %d, Title: %s, Body: %s", post.ID, post.Title, post.Body)

	_, err = pc.CreatePost(ctx, post)
	if err != nil {
		return fmt.Errorf("error CreatePost: %v", err)
	}

	log.Info("Created Post ID: %d, Title: %s, Body: %s", post.ID, post.Title, post.Body)

	updated_post, err := pc.UpdatePost(ctx, postId, post)
	if err != nil {
		return fmt.Errorf("error UpdatePost: %v", err)
	}
	log.Info("Updated Post ID: %d, Title: %s, Body: %s", updated_post.ID, updated_post.Title, updated_post.Body)

	pc.DeletePost(ctx, postId)
	log.Info("Test client operations completed successfully.")
	return nil
}
