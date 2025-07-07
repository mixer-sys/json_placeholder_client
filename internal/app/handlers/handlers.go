package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"json_placeholder_client/internal/app/config"
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

func NewPostClient(cfg *config.Config, client *http.Client) *PostClient {
	return &PostClient{
		BaseURL:  cfg.BaseURL,
		Client:   client,
		ProxyURL: cfg.ProxyURL,
	}
}

func (pc *PostClient) SetClient() {
	transport := &http.Transport{
		Proxy: http.ProxyURL(pc.ProxyURL),
	}

	pc.Client = &http.Client{Transport: transport, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return nil
	},
	}
}

func (pc *PostClient) GetPosts(ctx context.Context) (posts []Post, err error) {
	url := pc.BaseURL.String() + "/posts"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return posts, fmt.Errorf("failed to create request to %s: %w", url, err)
	}

	resp, err := pc.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, &posts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return posts, nil
}

func (pc *PostClient) GetPostByID(ctx context.Context, id int) (post Post, err error) {
	url := pc.BaseURL.String() + "/posts/" + strconv.Itoa(id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return post, fmt.Errorf("failed to create request to %s: %w", url, err)
	}

	resp, err := pc.Client.Do(req)
	if err != nil {
		return post, fmt.Errorf("failed to make GET request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return post, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, &post); err != nil {
		return post, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return post, nil
}

func (pc *PostClient) CreatePost(ctx context.Context, post Post) (created bool, err error) {

	url := pc.BaseURL.String() + "/posts"
	postData, err := json.Marshal(post)
	if err != nil {
		return false, fmt.Errorf("error serializing post to JSON: %v", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, io.NopCloser(bytes.NewBuffer(postData)))
	if err != nil {
		return false, fmt.Errorf("error creating request to %s: %v", url, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := pc.Client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error sending POST request to %s: %v", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error reading response from server: %v", err)
	}

	var createdPost Post
	if err := json.Unmarshal(body, &createdPost); err != nil {
		return false, fmt.Errorf("error deserializing response into Post struct: %v", err)
	}

	return true, nil
}

func (pc *PostClient) UpdatePost(ctx context.Context, id int, post Post) (Post, error) {

	url := pc.BaseURL.String() + "/posts/" + strconv.Itoa(id)
	postData, err := json.Marshal(post)
	if err != nil {
		return post, fmt.Errorf("error serializing post to JSON: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, io.NopCloser(bytes.NewBuffer(postData)))
	if err != nil {
		return post, fmt.Errorf("error creating PUT request for post ID %d: %v", id, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := pc.Client.Do(req)
	if err != nil {
		return post, fmt.Errorf("error sending PUT request to %s: %v", req.URL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return post, fmt.Errorf("error reading response from server for post ID %d: %v", id, err)
	}

	var updatedPost Post
	if err := json.Unmarshal(body, &updatedPost); err != nil {
		return post, fmt.Errorf("error deserializing response into Post struct for post ID %d: %v", id, err)
	}

	return updatedPost, nil
}

func (pc *PostClient) DeletePost(ctx context.Context, id int) error {
	url := pc.BaseURL.String() + "/posts/" + strconv.Itoa(id)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("error creating DELETE request for post ID %d: %v", id, err)
	}

	resp, err := pc.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending DELETE request to %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete post with ID %d: %s", id, resp.Status)
	}
	return nil
}
