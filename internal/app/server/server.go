package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"json_placeholder_client/internal/app/config"
	"json_placeholder_client/internal/app/logger"
)

func doRequestWithRetries(req *http.Request, retries int, delay time.Duration) (resp *http.Response, err error) {
	client := &http.Client{}

	for i := 0; i < retries; i++ {
		resp, err = client.Do(req)
		if err == nil {
			return resp, nil
		}
		time.Sleep(delay)
	}

	return nil, err
}

func handler(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequestWithContext(r.Context(), r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, "Failed to create request: ", http.StatusInternalServerError)
		return
	}

	for key, value := range r.Header {
		req.Header[key] = value
	}

	retries, err := config.GetRetries()
	if err != nil {
		http.Error(w, "Error getting retries from config: ", http.StatusInternalServerError)
		return
	}
	delay, err := config.GetRetryDelay()
	if err != nil {
		http.Error(w, "Error getting retry delay from config", http.StatusGatewayTimeout)
		return
	}

	resp, err := doRequestWithRetries(req, retries, time.Duration(delay)*time.Second)
	if err != nil {
		http.Error(w, "Request failed after retries: ", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, value := range resp.Header {
		w.Header()[key] = value
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func Run(ctx context.Context) error {
	log := logger.GetLogger()
	http.HandleFunc("/", handler)
	log.Info("Server is running on port 8080...")
	port, err := config.GetPort()
	if err != nil {
		log.Error("Error get port: %v", err)
		return fmt.Errorf("error getting port: %v", err)
	}
	address := ":" + port
	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Error("Server failed to start: %v", err)
		return fmt.Errorf("server failed to start: %v", err)
	}
	return nil
}
