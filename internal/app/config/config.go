package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL    string
	Port       string
	Retries    int
	RetryDelay int
	LogLevel   string
	ProxyURL   *url.URL
}

func GetConfig() (*Config, error) {
	config := &Config{}
	err := godotenv.Load()
	if err != nil {
		return config, fmt.Errorf("error load .env file: %v", err)
	}
	config.Port = os.Getenv("PORT")
	if config.Port == "" {
		return config, fmt.Errorf("PORT is not set")
	}
	config.BaseURL = os.Getenv("BASE_URL")
	if config.BaseURL == "" {
		return config, fmt.Errorf("BASE_URL is not set")
	}
	retriesStr := os.Getenv("RETRIES")
	if retriesStr == "" {
		return config, fmt.Errorf("RETRIES is not set")
	}
	config.Retries, err = strconv.Atoi(retriesStr)
	if err != nil {
		return config, fmt.Errorf("error converting RETRIES to int: %v", err)
	}
	retryDelayStr := os.Getenv("RETRY_DELAY")
	if retryDelayStr == "" {
		return config, fmt.Errorf("RETRY_DELAY is not set")
	}
	config.RetryDelay, err = strconv.Atoi(retryDelayStr)
	if err != nil {
		return config, fmt.Errorf("error converting RETRY_DELAY to int: %v", err)
	}
	config.LogLevel = os.Getenv("LOG_LEVEL")
	if config.LogLevel == "" {
		return config, fmt.Errorf("LOG_LEVEL is not set")
	}
	ProxyURLStr := os.Getenv("PROXY_URL")
	if ProxyURLStr == "" {
		return config, fmt.Errorf("PROXY_URL is not set")
	}
	config.ProxyURL, err = url.Parse(ProxyURLStr)
	if err != nil {
		return config, fmt.Errorf("error parse proxyURL: %v", err)
	}
	return config, nil

}
