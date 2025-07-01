package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	logger "json_placeholder_client/internal/app/logger"

	"github.com/joho/godotenv"
)

type Post struct {
	UserId int    `json:"user_id"`
	ID     int    `json:"id,omitempty"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func GetUrl() string {
	log := logger.GetLogger()
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file: %v", err)
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		log.Error("BASE_URL is not set in environment variables")
		os.Exit(1)
	}

	return baseURL
}

func GetPosts(ctx context.Context, client *http.Client) ([]Post, error) {
	baseURL := GetUrl()
	url := baseURL + "/posts"
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var posts []Post
	if err := json.Unmarshal(body, &posts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return posts, nil
}

func GetPostByID(ctx context.Context, client *http.Client, id int) (post Post, err error) {
	baseURL := GetUrl()
	url := baseURL + "/posts/" + strconv.Itoa(id)
	resp, err := client.Get(url)
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

func CreatePost(ctx context.Context, client *http.Client, post Post) (created bool, err error) {
	baseURL := GetUrl()
	url := baseURL + "/posts"
	postData, err := json.Marshal(post)
	if err != nil {
		return false, fmt.Errorf("error serializing post to JSON: %v", err)
	}

	resp, err := client.Post(url, "application/json", io.NopCloser(bytes.NewBuffer(postData)))
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

func UpdatePost(ctx context.Context, client *http.Client, id int, post Post) (Post, error) {
	baseURL := GetUrl()
	url := baseURL + "/posts/" + strconv.Itoa(id)
	postData, err := json.Marshal(post)
	if err != nil {
		return post, fmt.Errorf("error serializing post to JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, io.NopCloser(bytes.NewBuffer(postData)))
	if err != nil {
		return post, fmt.Errorf("error creating PUT request for post ID %d: %v", id, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
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

func DeletePost(ctx context.Context, client *http.Client, id int) error {
	baseURL := GetUrl()
	url := baseURL + "/posts/" + strconv.Itoa(id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("error creating DELETE request for post ID %d: %v", id, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending DELETE request to %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		log := logger.GetLogger()
		log.Info("Failed to delete post with ID %d: %s\n", id, resp.Status)
		return fmt.Errorf("failed to delete post with ID %d: %s", id, resp.Status)
	}
	return nil
}
