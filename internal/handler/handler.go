package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/ChebuRashkaRF/urlshortener/internal/util"
)

var urlMap = make(map[string]string)

func ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}

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

	urlMap[id] = url

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "http://localhost:8080/%s", id)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	id := r.URL.Path[1:]

	originalURL, ok := urlMap[id]
	if !ok {
		http.Error(w, "URL not found", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
