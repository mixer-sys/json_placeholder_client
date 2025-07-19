package client

import (
	"net/http"
	"net/url"
)

type Post struct {
	UserId int    `json:"user_id"`
	ID     int    `json:"id,omitempty"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type PostClient struct {
	BaseURL  *url.URL
	Client   *http.Client
	ProxyURL *url.URL
}
