package router

import (
	"github.com/ChebuRashkaRF/urlshortener/internal/handler"
	"github.com/ChebuRashkaRF/urlshortener/internal/logger"
	"github.com/ChebuRashkaRF/urlshortener/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(logger.WithLogging)
	r.Use(middleware.GzipMiddleware)

	r.Route("/", func(r chi.Router) {
		r.Post("/", handler.ShortenURLHandler)
		r.Post("/api/shorten", handler.ShortenURLJSONHandler)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.RedirectHandler)
		})
	})

	return r
}
