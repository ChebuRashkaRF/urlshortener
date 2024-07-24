package handler

import (
	"encoding/json"
	"fmt"
	"github.com/ChebuRashkaRF/urlshortener/internal/storage"
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/ChebuRashkaRF/urlshortener/cmd/config"
	"github.com/ChebuRashkaRF/urlshortener/internal/logger"
	"github.com/ChebuRashkaRF/urlshortener/internal/models"
	"github.com/ChebuRashkaRF/urlshortener/internal/util"
)

var URLStore *storage.URLStorage
var DB storage.Database

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

	shortURL := util.GenerateShortID(inputURL)

	URLStore.Set(shortURL, inputURL)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", config.Cnf.BaseURL, shortURL)
}

func ShortenURLJSONHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("decoding request")
	var req models.ShortenURLRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputURL := req.URL
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	if parsedURL.Scheme == "" && parsedURL.Host == "" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	shortURL := util.GenerateShortID(inputURL)

	URLStore.Set(shortURL, inputURL)

	res := models.ShortenURLResponse{
		Result: fmt.Sprintf("%s/%s", config.Cnf.BaseURL, shortURL),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	if err := enc.Encode(res); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "id")

	originalURL, ok := URLStore.Get(shortURL)
	if !ok {
		http.Error(w, "URL not found", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	if err := DB.Ping(); err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
