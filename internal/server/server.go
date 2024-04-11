package server

import (
	"github.com/ChebuRashkaRF/urlshortener/internal/handler"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
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

func Run() {
	log.Fatal(http.ListenAndServe(":8080", ShortenerRouter()))
}
