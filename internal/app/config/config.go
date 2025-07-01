package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func GetPort() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("error load .env file: %v", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("PORT is not set")
	}

	return port, nil
}

func GetRetries() (int, error) {
	err := godotenv.Load()
	if err != nil {
		return 0, fmt.Errorf("error load .env file: %v", err)
	}
	retriesStr := os.Getenv("RETRIES")
	if retriesStr == "" {
		return 0, fmt.Errorf("RETRIES is not set")
	}

	retries, err := strconv.Atoi(retriesStr)
	if err != nil {
		return 0, fmt.Errorf("error converting RETRIES to int: %v", err)
	}

	return retries, nil
}

func GetRetryDelay() (int, error) {
	err := godotenv.Load()
	if err != nil {
		return 0, fmt.Errorf("error load .env file: %v", err)
	}
	retryDelayStr := os.Getenv("RETRY_DELAY")
	if retryDelayStr == "" {
		return 0, fmt.Errorf("RETRY_DELAY is not set")
	}

	retryDelay, err := strconv.Atoi(retryDelayStr)
	if err != nil {
		return 0, fmt.Errorf("error converting RETRY_DELAY to int: %v", err)
	}

	return retryDelay, nil
}
