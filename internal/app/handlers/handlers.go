package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	logger "github.com/mixer-sys/json_placeholder_client/internal/app/logger"
)

type Post struct {
	UserID int    `json:"userId"`
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

func GetPosts(client *http.Client) []Post {
	baseURL := GetUrl()

	resp, err := client.Get(baseURL + "/posts")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var posts []Post
	if err := json.Unmarshal(body, &posts); err != nil {
		panic(err)
	}

	return posts
}

func GetPostByID(client *http.Client, id int) Post {
	baseURL := GetUrl()

	resp, err := client.Get(baseURL + "/posts/" + string(id))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var post Post
	if err := json.Unmarshal(body, &post); err != nil {
		panic(err)
	}

	return post
}

func CreatePost(client *http.Client, post Post) Post {
	baseURL := GetUrl()

	postData, err := json.Marshal(post)
	if err != nil {
		panic(err)
	}

	resp, err := client.Post(baseURL+"/posts", "application/json", io.NopCloser(bytes.NewBuffer(postData)))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var createdPost Post
	if err := json.Unmarshal(body, &createdPost); err != nil {
		panic(err)
	}

	return createdPost
}

func UpdatePost(client *http.Client, id int, post Post) Post {
	baseURL := GetUrl()

	postData, err := json.Marshal(post)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPut, baseURL+"/posts/"+string(id), io.NopCloser(bytes.NewBuffer(postData)))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var updatedPost Post
	if err := json.Unmarshal(body, &updatedPost); err != nil {
		panic(err)
	}

	return updatedPost
}

func DeletePost(client *http.Client, id int) {
	baseURL := GetUrl()

	req, err := http.NewRequest(http.MethodDelete, baseURL+"/posts/"+string(id), nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		log := logger.GetLogger()
		log.Error.Printf("Failed to delete post with ID %d: %s\n", id, resp.Status)
	}
}
