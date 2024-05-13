package config

import (
	"flag"
	"os"
)

var (
	serverAddress string
	baseURL       string
)

func parseFlags() {
	flag.StringVar(&serverAddress, "a", ":8080", "server address")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "base URL")
	flag.Parse()

	if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
		serverAddress = envServerAddress
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		baseURL = envBaseURL
	}
}
