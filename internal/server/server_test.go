package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRun(t *testing.T) {
	// Создаем тестовый сервер с заданным маршрутизатором
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, ShortenURLHandlerMock)
	mux.HandleFunc(`/{id}`, RedirectHandlerMock)
	server := httptest.NewServer(mux)
	defer server.Close()

	req1, err := http.NewRequest("GET", server.URL, nil)
	require.NoError(t, err)
	req2, err := http.NewRequest("GET", server.URL+"/abc123", nil)
	require.NoError(t, err)

	resp1, err := http.DefaultClient.Do(req1)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp1.StatusCode)
	defer resp1.Body.Close()

	resp2, err := http.DefaultClient.Do(req2)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)
	defer resp2.Body.Close()

}

func ShortenURLHandlerMock(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func RedirectHandlerMock(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
