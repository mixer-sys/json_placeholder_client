package server

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"json_placeholder_client/internal/app/config"
	"json_placeholder_client/internal/app/logger"
)

func doRequestWithRetries(req *http.Request, retries int,
	delay time.Duration) (resp *http.Response, err error) {

	client := &http.Client{}

	for i := 0; i < retries; i++ {
		resp, err = client.Do(req)
		if err == nil {
			return resp, nil
		}
		slog.Warn("Request failed",
			slog.Int("attempt", i+1),
			slog.String("error", err.Error()),
		)
		time.Sleep(delay)
	}

	return nil, err
}

func handler(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequestWithContext(r.Context(),
		r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, "Failed to create request: ",
			http.StatusInternalServerError)
		return
	}

	for key, value := range r.Header {
		req.Header[key] = value
	}
	cfg, err := config.Load()
	if err != nil {
		http.Error(w, "Error getting config: "+err.Error(), http.StatusInternalServerError)
		return
	}

	retries := cfg.Retries
	delay := cfg.RetryDelay

	resp, err := doRequestWithRetries(req, retries,
		time.Duration(delay)*time.Second)
	if err != nil {
		http.Error(w, "Request failed after retries: ",
			http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, value := range resp.Header {
		w.Header()[key] = value
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func Run(cfg *config.Config, ctx context.Context) error {
	log := logger.GetLogger(cfg)
	http.HandleFunc("/", handler)

	address := ":" + cfg.Port
	log.Info("Server is running on ",
		slog.String("address", address),
	)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Error("Server failed to start: %w", err)
		return fmt.Errorf("server failed to start: %w", err)
	}
	return nil
}
