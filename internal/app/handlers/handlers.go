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
	UserID int    `json:"user_id"`
	ID     int    `json:"id,omitempty"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func GetUrl() string {
	log := logger.GetLogger()
	err := godotenv.Load()
	if err != nil {
		log.Error.Println("Warning: .env file not found or couldn't be loaded")
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		log.Error.Println("BASE_URL is not set in environment variables")
	}

	return baseURL
}

func GetPosts(client *http.Client) ([]Post, error) {
	baseURL := GetUrl()

	resp, err := client.Get(baseURL + "/posts")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var posts []Post
	if err := json.Unmarshal(body, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetPostByID(client *http.Client, id int) (Post, error) {
	baseURL := GetUrl()

	resp, err := client.Get(baseURL + "/posts/" + strconv.Itoa(id))
	if err != nil {
		return Post{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Post{}, err
	}

	var post Post
	if err := json.Unmarshal(body, &post); err != nil {
		return Post{}, err
	}

	return post, nil
}

func CreatePost(client *http.Client, post Post) (Post, error) {
	baseURL := GetUrl()

	postData, err := json.Marshal(post)
	if err != nil {
		return Post{}, err
	}

	resp, err := client.Post(baseURL+"/posts", "application/json", io.NopCloser(bytes.NewBuffer(postData)))
	if err != nil {
		return Post{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Post{}, err
	}

	var createdPost Post
	if err := json.Unmarshal(body, &createdPost); err != nil {
		return Post{}, err
	}

	return createdPost, nil
}

func UpdatePost(client *http.Client, id int, post Post) (Post, error) {
	baseURL := GetUrl()

	postData, err := json.Marshal(post)
	if err != nil {
		return Post{}, err
	}

	req, err := http.NewRequest(http.MethodPut, baseURL+"/posts/"+strconv.Itoa(id), io.NopCloser(bytes.NewBuffer(postData)))
	if err != nil {
		return Post{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return Post{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Post{}, err
	}

	var updatedPost Post
	if err := json.Unmarshal(body, &updatedPost); err != nil {
		return Post{}, err
	}

	return updatedPost, nil
}

func DeletePost(client *http.Client, id int) error {
	baseURL := GetUrl()

	req, err := http.NewRequest(http.MethodDelete, baseURL+"/posts/"+strconv.Itoa(id), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		log := logger.GetLogger()
		log.Info.Printf("Failed to delete post with ID %d: %s\n", id, resp.Status)
		return fmt.Errorf("failed to delete post with ID %d: %s", id, resp.Status)
	}
	return nil
}
