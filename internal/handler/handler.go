package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ChebuRashkaRF/urlshortener/cmd/config"
	"github.com/ChebuRashkaRF/urlshortener/internal/util"
)

var URLMap = make(map[string]string)

func ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	url := string(body)

	if url == "" {
		http.Error(w, "Error empty body", http.StatusBadRequest)
		return
	}

	id := util.GenerateShortID(url)

	URLMap[id] = url

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", config.Cnf.BaseURL, id)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	originalURL, ok := URLMap[id]
	if !ok {
		http.Error(w, "URL not found", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
