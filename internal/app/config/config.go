package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetPort() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("PORT is not set")
	}

	return port, nil
}
