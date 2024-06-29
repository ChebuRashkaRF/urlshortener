package main

import (
	"github.com/ChebuRashkaRF/urlshortener/internal/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ChebuRashkaRF/urlshortener/cmd/config"
	"github.com/ChebuRashkaRF/urlshortener/internal/handler"
	"github.com/ChebuRashkaRF/urlshortener/internal/router"
)

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

func TestRun(t *testing.T) {
	tempFile, err := os.CreateTemp("", "urlstorage_test_*.json")
	require.NoError(t, err)

	r := router.NewRouter()
	ts := httptest.NewServer(r)

	// Извлечение порта из URL
	parts := strings.Split(ts.URL, ":")
	port := parts[len(parts)-1]

	config.Cnf = &config.Config{
		ServerAddress: ":" + port,
		BaseURL:       ts.URL,
	}

	urlStore, err := storage.NewURLStorage(tempFile.Name())
	require.NoError(t, err)
	handler.URLStore = urlStore

	t.Cleanup(func() {
		ts.Close()
		os.Remove(tempFile.Name())
		handler.URLStore.Close()
	})

	type wantPost struct {
		contentType string
		statusCode  int
	}

	wp := wantPost{
		contentType: "text/plain",
		statusCode:  http.StatusCreated,
	}

	resp, body := testRequest(t, ts, http.MethodPost, "/", "https://example.com")
	defer resp.Body.Close()

	assert.Equal(t, wp.statusCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
	assert.Equal(t, wp.contentType, resp.Header.Get("Content-Type"), "Content-Type не совпадает с ожидаемым")

	assert.Contains(t, body, config.Cnf.BaseURL)

	type wantGet struct {
		statusCode int
	}

	wg := wantGet{
		statusCode: http.StatusOK,
	}

	for k := range handler.URLStore.GetURLMap() {
		resp, _ := testRequest(t, ts, http.MethodGet, "/"+k, "")
		assert.Equal(t, wg.statusCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
		defer resp.Body.Close()
	}
}
