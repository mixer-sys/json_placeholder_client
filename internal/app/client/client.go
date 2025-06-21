package client

import (
	"encoding/json"
	"fmt"
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
		log.Println("Warning: .env file not found or couldn't be loaded")
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		log.Fatal("BASE_URL is not set in environment variables")
	}

	return baseURL
}

func GetPosts() []Post {
	log := logger.GetLogger()
	baseURL := GetUrl()

	resp, err := http.Get(baseURL + "/posts")
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
	log.Infof("Was %d gotten\n", len(posts))

	for _, post := range posts {
		fmt.Printf("Post ID: %d, Title: %s\n", post.ID, post.Title)
	}
	return posts
}
