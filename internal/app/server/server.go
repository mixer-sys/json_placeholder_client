package server

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"json_placeholder_client/internal/app/config"
	"json_placeholder_client/internal/app/logger"
)

func handler(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for key, value := range r.Header {
		req.Header[key] = value
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer resp.Body.Close()

	for key, value := range resp.Header {
		w.Header()[key] = value
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func Run(ctx *context.Context) error {
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
