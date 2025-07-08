package client

import (
	"context"
	"fmt"
	"log/slog"

	"json_placeholder_client/internal/app/config"
	handlers "json_placeholder_client/internal/app/handlers"
	"json_placeholder_client/internal/app/logger"
)

func RunTestClient(ctx context.Context) error {
	cfg, err := config.Load()
	log := logger.GetLogger(cfg)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}
	pc := handlers.NewPostClient(cfg, nil)

	pc.SetClient()

	posts, err := pc.GetPosts(ctx)
	if err != nil {
		return fmt.Errorf("error GetPosts: %w", err)
	}
	log.Info("Retrieved posts successfully. Total posts:",
		len(posts),
	)

	postId := 1
	post, err := pc.GetPostByID(ctx, postId)
	if err != nil {
		return fmt.Errorf("error GetPostByID: %w", err)
	}
	log.Info("Got post",
		slog.Int("id", post.ID),
		slog.String("title", post.Title),
		slog.String("body", post.Body),
	)
	_, err = pc.CreatePost(ctx, post)
	if err != nil {
		return fmt.Errorf("error CreatePost: %w", err)
	}

	log.Info("Created post",
		slog.Int("id", post.ID),
		slog.String("title", post.Title),
		slog.String("body", post.Body),
	)
	updatedPost, err := pc.UpdatePost(ctx, postId, post)
	if err != nil {
		return fmt.Errorf("error UpdatePost: %w", err)
	}
	log.Info("Updated post",
		slog.Int("id", updatedPost.ID),
		slog.String("title", updatedPost.Title),
		slog.String("body", updatedPost.Body),
	)
	pc.DeletePost(ctx, postId)
	log.Info("Test client operations completed successfully.",
		slog.Int("postId", postId),
	)
	return nil
}
