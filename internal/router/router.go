package router

import (
	"github.com/ChebuRashkaRF/urlshortener/internal/handler"
	"github.com/ChebuRashkaRF/urlshortener/internal/logger"
	"github.com/go-chi/chi/v5"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()

	r.With(logger.WithLogging).Route("/", func(r chi.Router) {
		r.Post("/", handler.ShortenURLHandler)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.RedirectHandler)
		})
	})

	return r
}
