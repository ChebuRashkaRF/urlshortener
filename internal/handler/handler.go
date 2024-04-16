package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"

	"github.com/ChebuRashkaRF/urlshortener/cmd/config"
	"github.com/ChebuRashkaRF/urlshortener/cmd/storage"
	"github.com/ChebuRashkaRF/urlshortener/internal/util"
)

var URLStore *storage.URLStorage

func init() {
	URLStore = storage.NewURLStorage()
}
func ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	inputURL := string(body)

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	if parsedURL.Scheme == "" && parsedURL.Host == "" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id := util.GenerateShortID(inputURL)

	URLStore.Set(id, inputURL)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", config.Cnf.BaseURL, id)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	originalURL, ok := URLStore.Get(id)
	if !ok {
		http.Error(w, "URL not found", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
