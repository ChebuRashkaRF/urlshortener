package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ChebuRashkaRF/urlshortener/cmd/config"
	"github.com/ChebuRashkaRF/urlshortener/internal/handler"
	"github.com/ChebuRashkaRF/urlshortener/internal/router"
	"github.com/ChebuRashkaRF/urlshortener/internal/storage"
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

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string, body string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestShortenURLHandler(t *testing.T) {
	ts := httptest.NewServer(router.NewRouter())
	defer ts.Close()

	// Извлечение порта из URL
	parts := strings.Split(ts.URL, ":")
	port := parts[len(parts)-1]

	config.Cnf = &config.Config{
		ServerAddress: ":" + port,
		BaseURL:       ts.URL,
	}

	handler.URLStore = storage.NewURLStorage()

	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name         string
		reqBody      string
		method       string
		want         want
		wantErr      string
		wantURLStore *storage.URLStorage
	}{
		{
			name:    "POST request method",
			reqBody: "https://example.com",
			method:  http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
			},
			wantURLStore: handler.URLStore,
		},
		{
			name:    "Invalid reqBody",
			reqBody: "",
			method:  http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
			},
			wantErr: "Invalid URL\n",
		},
		{
			name:    "Invalid URL",
			reqBody: "yandex.ru",
			method:  http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
			},
			wantErr: "Invalid URL\n",
		},
		{
			name:    "Invalid Method",
			reqBody: "https://example.com",
			method:  http.MethodGet,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusMethodNotAllowed,
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := testRequest(t, ts, tt.method, "/", tt.reqBody)
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				assert.Equal(t, tt.want.statusCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
				assert.Equal(t, body, tt.wantErr, "Не совпадает ошибка с ожидаемой")
				assert.Empty(t, handler.URLStore.GetURLMap())
				return
			}

			assert.Equal(t, tt.want.statusCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"), "Content-Type не совпадает с ожидаемым")

			assert.Contains(t, body, config.Cnf.BaseURL)
			assert.NotEmpty(t, tt.wantURLStore.GetURLMap())
			handler.URLStore = &storage.URLStorage{URLMap: make(map[string]string)}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	ts := httptest.NewServer(router.NewRouter())
	defer ts.Close()

	// Извлечение порта из URL
	parts := strings.Split(ts.URL, ":")
	port := parts[len(parts)-1]

	config.Cnf = &config.Config{
		ServerAddress: ":" + port,
		BaseURL:       ts.URL,
	}

	handler.URLStore = storage.NewURLStorage()

	tests := []struct {
		name           string
		method         string
		request        string
		wantStatusCode int
		wantErr        string
	}{
		{
			name:           "GET request method",
			method:         http.MethodGet,
			request:        "/abc123",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "Invalid Method",
			method:         http.MethodPost,
			request:        "/abc123",
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "URL Not Found",
			method:         http.MethodGet,
			request:        "/invalidid",
			wantStatusCode: http.StatusBadRequest,
			wantErr:        "URL not found\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.URLStore.Set("abc123", "https://example.com")
			defer delete(handler.URLStore.URLMap, "abc123")
			resp, body := testRequest(t, ts, tt.method, tt.request, "")
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatusCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")

			if tt.wantErr != "" {
				assert.Equal(t, body, tt.wantErr, "Не совпадает ошибка с ожидаемой")
			}
		})
	}
}
