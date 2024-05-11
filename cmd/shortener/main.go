package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ChebuRashkaRF/urlshortener/cmd/config"
	"github.com/ChebuRashkaRF/urlshortener/cmd/router"
)

func run(cnf *config.Config) {
	fmt.Println("Running server on", cnf.ServerAddress)
	log.Fatal(http.ListenAndServe(cnf.ServerAddress, router.NewRouter()))
}

func main() {
	var serverAddress string
	var baseURL string

	flag.StringVar(&serverAddress, "a", ":8080", "server address")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "base URL")
	flag.Parse()

	if envServerAddress := os.Getenv("SERVER_ADDRESS"); envServerAddress != "" {
		serverAddress = envServerAddress
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		baseURL = envBaseURL
	}

	config.Cnf = config.NewConfig(serverAddress, baseURL)

	run(config.Cnf)
}
