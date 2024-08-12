package config

import (
	"flag"
	"os"
)

var (
	serverAddress   string
	baseURL         string
	fileStoragePath string
	databaseDSN     string
)

func parseFlags() {
	flag.StringVar(&serverAddress, "a", ":8080", "server address")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "base URL")
	flag.StringVar(&fileStoragePath, "f", "", "file storage path")
	flag.StringVar(&databaseDSN, "d", "", "database DSN")
	flag.Parse()

	if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
		serverAddress = envServerAddress
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		baseURL = envBaseURL
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		fileStoragePath = envFileStoragePath
	}
	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		databaseDSN = envDatabaseDSN
	}
}
