package client

import (
	"net/url"

	"github.com/go-resty/resty/v2"
)

type Post struct {
	UserId int    `json:"user_id"`
	ID     int    `json:"id,omitempty"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type PostClient struct {
	BaseURL *url.URL
	Client  *resty.Client
}
