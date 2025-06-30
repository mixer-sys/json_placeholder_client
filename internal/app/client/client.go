package client

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	handlers "json_placeholder_client/internal/app/handlers"
	"json_placeholder_client/internal/app/logger"

	"github.com/go-resty/resty/v2"
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

func GetClient() (*resty.Client, error) {
	proxyURL, err := GetProxyURL()
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := resty.New().
		SetTransport(transport).SetRedirectPolicy(resty.FlexibleRedirectPolicy(302))
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

	_, err = handlers.UpdatePost(client, postId, post)
	if err != nil {
		return fmt.Errorf("error UpdatePost: %v", err)
	}
	log.Info("Updated Post ID: %d, Title: %s, Body: %s", post.ID, post.Title, post.Body)

	handlers.DeletePost(client, postId)
	log.Info("Test client operations completed successfully.")
	return nil
}
