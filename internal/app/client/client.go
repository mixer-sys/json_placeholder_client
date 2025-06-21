package client

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/mixer-sys/json_placeholder_client/internal/app/logger"
)

func GetProxyURL() (*url.URL, error) {
	log := logger.GetLogger()
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

func client() {
	log := logger.GetLogger()

	proxyURL, err := GetProxyURL()
	if err != nil {
		log.Error.Println("Error getting proxy URL:", err)
	}

	baseURL := os.Getenv("BASE_URL")

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := &http.Client{Transport: transport}

	resp, err := client.Get(baseURL)
	if err != nil {
		log.Error.Println("Error making GET request:", err)
		return
	}
	defer resp.Body.Close()

	log.Info.Printf("Response status: %s", resp.Status)
}
