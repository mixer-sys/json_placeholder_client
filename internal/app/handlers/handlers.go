package handlers

import (
	"bytes"
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

func GetPosts(client *http.Client) ([]Post, error) {
	baseURL := GetUrl()

	resp, err := client.Get(baseURL + "/posts")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	var posts []Post
	if err := json.Unmarshal(body, &posts); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return posts, nil
}

func GetPostByID(client *http.Client, id int) (post Post, err error) {
	baseURL := GetUrl()

	resp, err := client.Get(baseURL + "/posts/" + strconv.Itoa(id))
	if err != nil {
		return post, fmt.Errorf("%v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return post, fmt.Errorf("%v", err)
	}

	if err := json.Unmarshal(body, &post); err != nil {
		return post, fmt.Errorf("%v", err)
	}

	return post, nil
}

func CreatePost(client *http.Client, post Post) (created bool, err error) {
	baseURL := GetUrl()

	postData, err := json.Marshal(post)
	if err != nil {
		return false, fmt.Errorf("%v", err)
	}

	resp, err := client.Post(baseURL+"/posts", "application/json", io.NopCloser(bytes.NewBuffer(postData)))
	if err != nil {
		return false, fmt.Errorf("%v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("%v", err)
	}

	var createdPost Post
	if err := json.Unmarshal(body, &createdPost); err != nil {
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
}

func UpdatePost(client *http.Client, id int, post Post) (Post, error) {
	baseURL := GetUrl()

	postData, err := json.Marshal(post)
	if err != nil {
		return post, fmt.Errorf("%v", err)
	}

	req, err := http.NewRequest(http.MethodPut, baseURL+"/posts/"+strconv.Itoa(id), io.NopCloser(bytes.NewBuffer(postData)))
	if err != nil {
		return post, fmt.Errorf("%v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return post, fmt.Errorf("%v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return post, fmt.Errorf("%v", err)
	}

	var updatedPost Post
	if err := json.Unmarshal(body, &updatedPost); err != nil {
		return post, fmt.Errorf("%v", err)
	}

	return updatedPost, nil
}

func DeletePost(client *http.Client, id int) error {
	baseURL := GetUrl()

	req, err := http.NewRequest(http.MethodDelete, baseURL+"/posts/"+strconv.Itoa(id), nil)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		log := logger.GetLogger()
		log.Info("Failed to delete post with ID %d: %s\n", id, resp.Status)
		return fmt.Errorf("failed to delete post with ID %d: %s", id, resp.Status)
	}
	return nil
}
