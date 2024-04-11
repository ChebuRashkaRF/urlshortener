package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

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

func main() {
	log.Fatal(http.ListenAndServe(":8080", ShortenerRouter()))
}
