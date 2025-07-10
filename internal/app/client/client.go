package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/exp/slog"

	"json_placeholder_client/internal/app/config"

	"json_placeholder_client/internal/app/logger"
)

type Post struct {
	UserId int    `json:"user_id"`
	ID     int    `json:"id,omitempty"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type PostClient struct {
	BaseURL  *url.URL
	Client   *http.Client
	ProxyURL *url.URL
}

func NewPostClient(cfg *config.Config) *PostClient {
	transport := &http.Transport{
		Proxy: http.ProxyURL(cfg.ProxyURL),
	}

	client := &http.Client{
		Transport: transport,
		CheckRedirect: func(
			req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	return &PostClient{
		BaseURL:  cfg.BaseURL,
		Client:   client,
		ProxyURL: cfg.ProxyURL,
	}
}

func (pc *PostClient) GetPosts(ctx context.Context) (
	[]Post, error) {
	url := pc.BaseURL.String() + "/posts"
	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet, url, nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create request to %s: %w", url, err,
		)
	}

	resp, err := pc.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to make GET request to %s: %w", url, err,
		)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read response body: %w", err,
		)
	}
	var posts []Post
	if err := json.Unmarshal(body, &posts); err != nil {
		return nil, fmt.Errorf(
			"failed to unmarshal response body: %w", err,
		)
	}

	return posts, nil
}

func (pc *PostClient) GetPostByID(ctx context.Context, id int) (
	post Post, err error) {
	url := pc.BaseURL.String() + "/posts/" + strconv.Itoa(id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return post, fmt.Errorf(
			"failed to create request to %s: %w", url, err)
	}

	resp, err := pc.Client.Do(req)
	if err != nil {
		return post, fmt.Errorf(
			"failed to make GET request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return post, fmt.Errorf(
			"failed to read response body: %w", err)
	}
	if err := json.Unmarshal(body, &post); err != nil {
		return post, fmt.Errorf(
			"failed to unmarshal response body: %w", err)
	}

	return post, nil
}

func (pc *PostClient) CreatePost(ctx context.Context, post Post) (
	bool, error) {
	url := pc.BaseURL.String() + "/posts"
	postData, err := json.Marshal(post)
	if err != nil {
		return false, fmt.Errorf(
			"error serializing post to JSON: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, io.NopCloser(bytes.NewBuffer(postData)))
	if err != nil {
		return false, fmt.Errorf(
			"error creating request to %s: %w", url, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := pc.Client.Do(req)
	if err != nil {
		return false, fmt.Errorf(
			"error sending POST request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf(
			"error reading response from server: %w", err)
	}

	var createdPost Post
	if err := json.Unmarshal(body, &createdPost); err != nil {
		return false, fmt.Errorf(
			"error deserializing response into Post struct: %w", err,
		)
	}

	return true, nil
}

func (pc *PostClient) UpdatePost(ctx context.Context, id int,
	post Post) (bool, error) {
	url := pc.BaseURL.String() + "/posts/" + strconv.Itoa(id)
	postData, err := json.Marshal(post)
	if err != nil {
		return false, fmt.Errorf(
			"error serializing post to JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx, http.MethodPut, url,
		io.NopCloser(bytes.NewBuffer(postData)))
	if err != nil {
		return false, fmt.Errorf(
			"error creating PUT request for post ID %d: %w",
			id, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := pc.Client.Do(req)
	if err != nil {
		return false, fmt.Errorf(
			"error sending PUT request to %s: %w",
			req.URL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf(
			"error reading response from server for post ID %d: %w",
			id, err)
	}

	var updatedPost Post
	if err := json.Unmarshal(body, &updatedPost); err != nil {
		return false, fmt.Errorf(
			"error deserializing response into Post struct for post ID %d: %w",
			id, err)
	}

	return true, nil
}

func (pc *PostClient) DeletePost(ctx context.Context, id int) error {
	url := pc.BaseURL.String() + "/posts/" + strconv.Itoa(id)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("error creating DELETE request for post ID %d: %w",
			id, err)
	}

	resp, err := pc.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending DELETE request to %s: %w",
			url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete post with ID %d: %s",
			id, resp.Status)
	}
	return nil
}

func RunTestClient(ctx context.Context) error {
	cfg, err := config.Load()
	log := logger.GetLogger(cfg)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}
	pc := NewPostClient(cfg)

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
