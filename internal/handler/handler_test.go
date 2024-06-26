package handler_test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/ChebuRashkaRF/urlshortener/internal/middleware"
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
	"github.com/ChebuRashkaRF/urlshortener/internal/storage"
)

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string, body io.Reader, headers map[string]string) (*http.Response, []byte) {
	req, err := http.NewRequest(method, ts.URL+path, body)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	var respBody []byte

	if resp.Header.Get("Content-Encoding") == "gzip" {
		gr, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)
		defer gr.Close()
		respBody, err = io.ReadAll(gr)
		require.NoError(t, err)
	} else {
		respBody, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
	}

	return resp, respBody
}

func setupTestRouter(handlerFunc http.HandlerFunc, method, path string, middlewareList ...func(http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()
	for _, mw := range middlewareList {
		r.Use(mw)
	}
	r.MethodFunc(method, path, handlerFunc)
	return r
}

func TestShortenURLHandler(t *testing.T) {
	r := setupTestRouter(handler.ShortenURLHandler, http.MethodPost, "/")
	ts := httptest.NewServer(r)
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
		headers      map[string]string
		want         want
		wantErr      string
		wantURLStore *storage.URLStorage
	}{
		{
			name:    "POST request method",
			reqBody: "https://example.com",
			method:  http.MethodPost,
			headers: nil,
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
			headers: nil,
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
			headers: nil,
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
			headers: nil,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusMethodNotAllowed,
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := testRequest(t, ts, tt.method, "/", strings.NewReader(tt.reqBody), tt.headers)
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				assert.Equal(t, tt.want.statusCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
				assert.Equal(t, string(body), tt.wantErr, "Не совпадает ошибка с ожидаемой")
				assert.Empty(t, handler.URLStore.GetURLMap())
				return
			}

			assert.Equal(t, tt.want.statusCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"), "Content-Type не совпадает с ожидаемым")

			assert.Contains(t, string(body), config.Cnf.BaseURL)
			assert.NotEmpty(t, tt.wantURLStore.GetURLMap())
			handler.URLStore = &storage.URLStorage{URLMap: make(map[string]string)}
		})
	}
}

func TestShortenURLJSONHandler(t *testing.T) {
	r := setupTestRouter(handler.ShortenURLJSONHandler, http.MethodPost, "/api/shorten")
	ts := httptest.NewServer(r)
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
		headers      map[string]string
		want         want
		wantErr      string
		wantURLStore *storage.URLStorage
	}{
		{
			name:    "POST request method",
			reqBody: `{"url": "https://example.com"}`,
			method:  http.MethodPost,
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
			},
			wantURLStore: handler.URLStore,
		},
		{
			name:    "Invalid reqBody",
			reqBody: `{"url": ""}`,
			method:  http.MethodPost,
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusBadRequest,
			},
			wantErr: "Invalid URL\n",
		},
		{
			name:    "Invalid URL",
			reqBody: `{"url": "yandex.ru"}`,
			method:  http.MethodPost,
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusBadRequest,
			},
			wantErr: "Invalid URL\n",
		},
		{
			name:    "Invalid Method",
			reqBody: `{"url": "https://example.com"}`,
			method:  http.MethodGet,
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusMethodNotAllowed,
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := testRequest(t, ts, tt.method, "/api/shorten", strings.NewReader(tt.reqBody), tt.headers)
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				assert.Equal(t, tt.want.statusCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
				assert.Equal(t, string(body), tt.wantErr, "Не совпадает ошибка с ожидаемой")
				assert.Empty(t, handler.URLStore.GetURLMap())
				return
			}

			assert.Equal(t, tt.want.statusCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"), "Content-Type не совпадает с ожидаемым")

			for id := range handler.URLStore.GetURLMap() {
				successBody := fmt.Sprintf(`{"result": "%s/%s"}`, ts.URL, id)
				assert.JSONEq(t, successBody, string(body), "ответ не совпадает с ожидаемым")
			}
			handler.URLStore = &storage.URLStorage{URLMap: make(map[string]string)}
		})
	}
}

func TestGzipCompressionShortenURLHandler(t *testing.T) {
	r := setupTestRouter(handler.ShortenURLHandler, http.MethodPost, "/", middleware.GzipMiddleware)
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Извлечение порта из URL
	parts := strings.Split(ts.URL, ":")
	port := parts[len(parts)-1]

	config.Cnf = &config.Config{
		ServerAddress: ":" + port,
		BaseURL:       ts.URL,
	}

	requestBody := "https://example.com"

	t.Run("sends_gzip", func(t *testing.T) {
		handler.URLStore = storage.NewURLStorage()
		headers := map[string]string{
			"Content-Encoding": "gzip",
		}
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		resp, body := testRequest(t, ts, http.MethodPost, "/", buf, headers)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		for id := range handler.URLStore.GetURLMap() {
			successBody := fmt.Sprintf("%s/%s", ts.URL, id)
			assert.Equal(t, successBody, string(body), "ответ не совпадает с ожидаемым")
		}
	})

	t.Run("accepts_gzip", func(t *testing.T) {
		handler.URLStore = storage.NewURLStorage()
		headers := map[string]string{
			"Accept-Encoding": "gzip",
		}

		buf := bytes.NewBufferString(requestBody)

		resp, body := testRequest(t, ts, http.MethodPost, "/", buf, headers)

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		for id := range handler.URLStore.GetURLMap() {
			successBody := fmt.Sprintf("%s/%s", ts.URL, id)
			assert.Equal(t, successBody, string(body), "ответ не совпадает с ожидаемым")
		}
	})
}

func TestGzipCompressionShortenURLJSONHandler(t *testing.T) {
	r := setupTestRouter(handler.ShortenURLJSONHandler, http.MethodPost, "/api/shorten", middleware.GzipMiddleware)
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Извлечение порта из URL
	parts := strings.Split(ts.URL, ":")
	port := parts[len(parts)-1]

	config.Cnf = &config.Config{
		ServerAddress: ":" + port,
		BaseURL:       ts.URL,
	}

	requestBody := `{"url": "https://example.com"}`

	t.Run("sends_gzip", func(t *testing.T) {
		handler.URLStore = storage.NewURLStorage()
		headers := map[string]string{
			"Content-Encoding": "gzip",
			"Content-Type":     "application/json",
			"Accept-Encoding":  "",
		}
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		resp, body := testRequest(t, ts, http.MethodPost, "/api/shorten", buf, headers)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		for id := range handler.URLStore.GetURLMap() {
			successBody := fmt.Sprintf(`{"result":"%s/%s"}`, ts.URL, id)
			assert.JSONEq(t, successBody, string(body), "ответ не совпадает с ожидаемым")
		}
	})

	t.Run("accepts_gzip", func(t *testing.T) {
		handler.URLStore = storage.NewURLStorage()
		headers := map[string]string{
			"Content-Type":    "application/json",
			"Accept-Encoding": "gzip",
		}

		buf := bytes.NewBufferString(requestBody)

		resp, body := testRequest(t, ts, http.MethodPost, "/api/shorten", buf, headers)

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()
		//zr, err := gzip.NewReader(resp.Body)
		//require.NoError(t, err)
		//
		//b, err := io.ReadAll(zr)
		//require.NoError(t, err)

		for id := range handler.URLStore.GetURLMap() {
			successBody := fmt.Sprintf(`{"result":"%s/%s"}`, ts.URL, id)
			assert.JSONEq(t, successBody, string(body), "ответ не совпадает с ожидаемым")
		}
	})
}

func TestRedirectHandler(t *testing.T) {
	r := setupTestRouter(handler.RedirectHandler, http.MethodGet, "/{id}")
	ts := httptest.NewServer(r)
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
		headers        map[string]string
		request        string
		wantStatusCode int
		wantErr        string
	}{
		{
			name:           "GET request method",
			method:         http.MethodGet,
			headers:        nil,
			request:        "/abc123",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "Invalid Method",
			method:         http.MethodPost,
			headers:        nil,
			request:        "/abc123",
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "URL Not Found",
			method:         http.MethodGet,
			headers:        nil,
			request:        "/invalidid",
			wantStatusCode: http.StatusBadRequest,
			wantErr:        "URL not found\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.URLStore.Set("abc123", "https://example.com")
			defer delete(handler.URLStore.URLMap, "abc123")
			resp, body := testRequest(t, ts, tt.method, tt.request, strings.NewReader(""), tt.headers)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatusCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")

			if tt.wantErr != "" {
				assert.Equal(t, string(body), tt.wantErr, "Не совпадает ошибка с ожидаемой")
			}
		})
	}
}
