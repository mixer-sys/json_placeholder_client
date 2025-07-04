package config

import (
	"fmt"
	"net/url"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL    string   `env:"BASE_URL" envDefault:"http://jsonplaceholder.typicode.com"`
	Port       string   `env:"PORT" envDefault:"8080"`
	Retries    int      `env:"RETRIES" envDefault:"3"`
	RetryDelay int      `env:"RETRY_DELAY" envDefault:"1000"`
	LogLevel   string   `env:"LOG_LEVEL" envDefault:"info"`
	ProxyURL   *url.URL `env:"PROXY_URL" envDefault:"http://localhost:8080"`
}

func Load() (*Config, error) {
	var config Config
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}
	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("error parsing environment variables: %w", err)
	}
	return &config, nil
}
