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
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}
	proxyStr := os.Getenv("PROXY_URL")
	if proxyStr == "" {
		return nil, fmt.Errorf("PROXY_URL is not set")

	}

	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		return nil, err
	}
	return proxyURL, nil
}

func GetClient() (*http.Client, error) {
	proxyURL, err := GetProxyURL()
	if err != nil {
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

func RunTestClient() error {
	log := logger.GetLogger()

	client, err := GetClient()
	if err != nil {
		return fmt.Errorf("error creating HTTP client: %v", err)

	}

	posts, err := handlers.GetPosts(client)
	if err != nil {
		return fmt.Errorf("error GetPosts: %v", err)
	}
	log.Info("Retrieved posts successfully. Total posts:", len(posts))

	postId := 1
	post, err := handlers.GetPostByID(client, postId)
	if err != nil {
		return fmt.Errorf("error GetPostByID: %v", err)
	}
	log.Info("Got ID: %d, Title: %s, Body: %s", post.ID, post.Title, post.Body)

	_, err = handlers.CreatePost(client, post)
	if err != nil {
		return fmt.Errorf("error CreatePost: %v", err)
	}

	log.Info("Created Post ID: %d, Title: %s, Body: %s", post.ID, post.Title, post.Body)

	updated_post, err := handlers.UpdatePost(client, postId, post)
	if err != nil {
		return fmt.Errorf("error UpdatePost: %v", err)
	}
	log.Info("Updated Post ID: %d, Title: %s, Body: %s", updated_post.ID, updated_post.Title, updated_post.Body)

	handlers.DeletePost(client, postId)
	log.Info("Test client operations completed successfully.")
	return nil
}
