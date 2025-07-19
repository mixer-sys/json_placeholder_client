package client

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"golang.org/x/exp/slog"

	"json_placeholder_client/internal/app/config"

	"json_placeholder_client/internal/app/logger"
)

func New(cfg *config.Config) *PostClient {
	client := resty.New().
		SetBaseURL(cfg.BaseURL.String()).
		SetProxy(cfg.ProxyURL.String()).
		SetHeader("Content-Type", "application/json")

	return &PostClient{
		BaseURL: cfg.BaseURL,
		Client:  client,
	}
}

func (pc *PostClient) GetPosts(ctx context.Context) (
	posts []Post, err error) {
	resp, err := pc.Client.R().
		SetContext(ctx).
		SetResult(&posts).
		Get("/posts")
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("error response: %s", resp.Status())
	}
	return posts, nil
}

func (pc *PostClient) GetPostByID(ctx context.Context, id int) (
	post Post, err error) {
	resp, err := pc.Client.R().
		SetContext(ctx).
		SetResult(&post).
		Get(fmt.Sprintf("/posts/%d", id))
	if err != nil {
		return post, fmt.Errorf("failed to make GET request: %w", err)
	}
	if resp.IsError() {
		return post, fmt.Errorf("error response: %s", resp.Status())
	}
	return post, nil
}

func (pc *PostClient) CreatePost(ctx context.Context, post Post) (
	bool, error) {
	resp, err := pc.Client.R().
		SetContext(ctx).
		SetBody(post).
		Post("/posts")
	if err != nil {
		return false, fmt.Errorf("error sending POST request: %w", err)
	}
	if resp.IsError() {
		return false, fmt.Errorf("error response: %s", resp.Status())
	}
	return true, nil
}

func (pc *PostClient) UpdatePost(ctx context.Context, id int,
	post Post) (bool, error) {
	resp, err := pc.Client.R().
		SetContext(ctx).
		SetBody(post).
		Put(fmt.Sprintf("/posts/%d", id))
	if err != nil {
		return false, fmt.Errorf("error sending PUT request: %w", err)
	}
	if resp.IsError() {
		return false, fmt.Errorf("error response: %s", resp.Status())
	}
	return true, nil
}

func (pc *PostClient) DeletePost(ctx context.Context, id int) error {
	resp, err := pc.Client.R().
		SetContext(ctx).
		Delete(fmt.Sprintf("/posts/%d", id))
	if err != nil {
		return fmt.Errorf("error sending DELETE request: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("failed to delete post with ID %d: %s", id, resp.Status())
	}
	return nil
}

func RunTestClient(ctx context.Context) error {
	cfg, err := config.Load()
	log := logger.GetLogger(cfg)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}
	pc := New(cfg)

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
	_, err = pc.UpdatePost(ctx, postId, post)
	if err != nil {
		return fmt.Errorf("error UpdatePost: %w", err)
	}
	log.Info("Updated post",
		slog.Int("id", post.ID),
		slog.String("title", post.Title),
		slog.String("body", post.Body),
	)
	pc.DeletePost(ctx, postId)
	log.Info("Test client operations completed successfully.",
		slog.Int("postId", postId),
	)
	return nil
}
