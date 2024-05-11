package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ChebuRashkaRF/urlshortener/cmd/config"
	"github.com/ChebuRashkaRF/urlshortener/cmd/router"
	"github.com/ChebuRashkaRF/urlshortener/internal/handler"
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
	ts := httptest.NewServer(router.NewRouter())
	defer ts.Close()

	// Извлечение порта из URL
	parts := strings.Split(ts.URL, ":")
	port := parts[len(parts)-1]

	config.Cnf = &config.Config{
		ServerAddress: ":" + port,
		BaseURL:       ts.URL,
	}

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
