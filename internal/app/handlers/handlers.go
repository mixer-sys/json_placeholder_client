package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	logger "json_placeholder_client/internal/app/logger"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

type Post struct {
	UserId int    `json:"user_id"`
	ID     int    `json:"id,omitempty"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func GetUrl() string {
	log := logger.GetLogger()
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file: %v", err)
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		log.Error("BASE_URL is not set in environment variables")
		os.Exit(1)
	}

	return baseURL
}

func GetPosts(client *resty.Client) (posts []Post, err error) {
	baseURL := GetUrl()

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get(baseURL + "/posts/")

	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error response from server: %s", resp.Status())
	}

	if err := json.Unmarshal(resp.Body(), &posts); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return posts, nil
}

func GetPostByID(client *resty.Client, id int) (post Post, err error) {
	baseURL := GetUrl()

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get(baseURL + "/posts/" + strconv.Itoa(id))

	if err != nil {
		return post, fmt.Errorf("%v", err)
	}

	if err := json.Unmarshal(resp.Body(), &post); err != nil {
		return post, fmt.Errorf("%v", err)
	}

	return post, nil
}

func CreatePost(client *resty.Client, post Post) (created bool, err error) {
	baseURL := GetUrl()

	postData, err := json.Marshal(post)
	if err != nil {
		return false, fmt.Errorf("%v", err)
	}

	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(postData).
		Post(baseURL + "/posts")

	if err != nil {
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
}

func UpdatePost(client *resty.Client, id int, post Post) (updated bool, err error) {
	baseURL := GetUrl()

	putData, err := json.Marshal(post)
	if err != nil {
		return false, fmt.Errorf("%v", err)
	}

	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(putData).
		Put(baseURL + "/posts/" + strconv.Itoa(id))

	if err != nil {
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
}

func DeletePost(client *resty.Client, id int) (err error) {
	baseURL := GetUrl()

	_, err = client.R().
		Delete(baseURL + "/posts/" + strconv.Itoa(id))

	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}
