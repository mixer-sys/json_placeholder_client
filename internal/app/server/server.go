package server

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"json_placeholder_client/internal/app/logger"

	"github.com/joho/godotenv"
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
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, value := range resp.Header {
		w.Header()[key] = value
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

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
func Server() {
	log := logger.GetLogger()
	http.HandleFunc("/", handler)
	log.Info("Server is running on port 8080...")
	port, err := GetPort()
	if err != nil {
		log.Error("Error get port: %v", err)
	}
	address := ":" + port
	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Error("Server failed to start: %v", err)
	}
}
