package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ShortenerRouter() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", ShortenURLHandler)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", RedirectHandler)
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
	URLMap = map[string]string{}
	ts := httptest.NewServer(ShortenerRouter())
	defer ts.Close()

	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name       string
		reqBody    string
		method     string
		want       want
		wantErr    string
		wantURLMap map[string]string
	}{
		{
			name:    "POST request method",
			reqBody: "https://example.com",
			method:  http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
			},
			wantURLMap: URLMap,
		},
		{
			name:    "Invalid reqBody",
			reqBody: "",
			method:  http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
			},
			wantErr: "Error empty body\n",
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

			if resp.StatusCode != http.StatusCreated {
				assert.Equal(t, tt.want.statusCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
				assert.Equal(t, body, tt.wantErr, "Не совпадает ошибка с ожидаемой")
				assert.Empty(t, tt.wantURLMap)
				return
			}

			assert.Equal(t, tt.want.statusCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"), "Content-Type не совпадает с ожидаемым")

			assert.Contains(t, body, "http://localhost:8080/")
			assert.NotEmpty(t, tt.wantURLMap)
			URLMap = map[string]string{}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	ts := httptest.NewServer(ShortenerRouter())
	defer ts.Close()

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
			URLMap["abc123"] = "https://example.com"
			defer delete(URLMap, "abc123")
			resp, body := testRequest(t, ts, tt.method, tt.request, "")

			assert.Equal(t, tt.wantStatusCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")

			if tt.wantErr != "" {
				assert.Equal(t, body, tt.wantErr, "Не совпадает ошибка с ожидаемой")
			}
		})
	}
}
