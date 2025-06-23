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
		return nil, err
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
	} else {
		log.Info.Println("Retrieved posts successfully. Total posts:", len(posts))
	}

	postID := 1
	post, err := handlers.GetPostByID(client, postID)
	if err != nil {
		return fmt.Errorf("error GetPostByID: %v", err)
	} else {
		log.Info.Printf("Got ID: %d, Title: %s, Body: %s", post.ID, post.Title, post.Body)
	}

	created_post, err := handlers.CreatePost(client, post)
	if err != nil {
		return fmt.Errorf("error CreatePost: %v", err)
	} else {
		log.Info.Printf("Created Post ID: %d, Title: %s, Body: %s", created_post.ID, created_post.Title, created_post.Body)
	}

	updated_post, err := handlers.UpdatePost(client, postID, post)
	if err != nil {
		return fmt.Errorf("error UpdatePost: %v", err)
	} else {
		log.Info.Printf("Updated Post ID: %d, Title: %s, Body: %s", updated_post.ID, updated_post.Title, updated_post.Body)
	}

	handlers.DeletePost(client, postID)
	log.Info.Println("Test client operations completed successfully.")
	return nil
}
