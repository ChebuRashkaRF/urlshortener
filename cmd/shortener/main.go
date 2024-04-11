package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ChebuRashkaRF/urlshortener/cmd/config"
	"github.com/ChebuRashkaRF/urlshortener/internal/handler"
)

func ShortenerRouter() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", handler.ShortenURLHandler)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.RedirectHandler)
		})
	})

	return r
}

func run(cnf *config.Config) {
	fmt.Println("Running server on", cnf.ServerAddress)
	log.Fatal(http.ListenAndServe(cnf.ServerAddress, ShortenerRouter()))
}

func main() {
	serverAddress := flag.String("a", ":8080", "server address")
	baseURL := flag.String("b", "http://localhost:8080", "base URL")
	flag.Parse()

	config.Cnf = config.NewConfig(*serverAddress, *baseURL)

	run(config.Cnf)
}
