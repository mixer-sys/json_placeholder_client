package client

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	handlers "json_placeholder_client/internal/app/handlers"
	"json_placeholder_client/internal/app/logger"

	"github.com/joho/godotenv"
)

func GetProxyURL() (*url.URL, error) {
	log := logger.GetLogger()
	err := godotenv.Load()
	if err != nil {
		log.Info.Println("Error loading .env file")
	}
	proxyStr := os.Getenv("PROXY_URL")
	if proxyStr == "" {
		log.Error.Println("PROXY_URL is not set in environment variables")
		return nil, fmt.Errorf("PROXY_URL is not set")
	}

	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		log.Error.Printf("Invalid PROXY_URL: %s, error: %v", proxyStr, err)
		return nil, err
	}
	return proxyURL, nil
}

func GetClient() (*http.Client, error) {
	log := logger.GetLogger()

	proxyURL, err := GetProxyURL()
	if err != nil {
		log.Error.Println("Error getting proxy URL:", err)
		return nil, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := &http.Client{Transport: transport, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return nil
	},
	}
	return client, nil
}

func RunTestClient() {
	log := logger.GetLogger()

	client, err := GetClient()
	if err != nil {
		log.Error.Println("Error creating HTTP client:", err)
		return
	}

	posts := handlers.GetPosts(client)
	log.Info.Println("Retrieved posts successfully. Total posts:", len(posts))

	postID := 1
	post := handlers.GetPostByID(client, postID)
	log.Info.Printf("Got ID: %d, Title: %s, Body: %s", post.ID, post.Title, post.Body)

	created_post := handlers.CreatePost(client, post)
	log.Info.Printf("Created Post ID: %d, Title: %s, Body: %s", created_post.ID, created_post.Title, created_post.Body)

	updated_post := handlers.UpdatePost(client, postID, post)
	log.Info.Printf("Updated Post ID: %d, Title: %s, Body: %s", updated_post.ID, updated_post.Title, updated_post.Body)

	handlers.DeletePost(client, postID)
	log.Info.Println("Test client operations completed successfully.")
}
