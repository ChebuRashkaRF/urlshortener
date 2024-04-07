package server

import (
	"log"
	"net/http"

	"github.com/ChebuRashkaRF/urlshortener/internal/handler"
)

func Run() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handler.ShortenURLHandler)
	mux.HandleFunc(`/{id}`, handler.RedirectHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
