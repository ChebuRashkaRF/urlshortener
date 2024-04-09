package handler

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenURLHandler(t *testing.T) {
	urlMap = map[string]string{}
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
		wantUrlMap map[string]string
	}{
		{
			name:    "POST request method",
			reqBody: "https://example.com",
			method:  http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
			},
			wantUrlMap: urlMap,
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
				statusCode:  http.StatusBadRequest,
			},
			wantErr: "Only POST requests are allowed!\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, `/`, bytes.NewBufferString(tt.reqBody))
			w := httptest.NewRecorder()

			ShortenURLHandler(w, request)

			result := w.Result()

			if tt.wantErr != "" {
				assert.Equal(t, tt.want.statusCode, result.StatusCode, "Код ответа не совпадает с ожидаемым")
				assert.Equal(t, w.Body.String(), tt.wantErr, "Не совпадает ошибка с ожидаемой")
				return
			}

			assert.Equal(t, tt.want.statusCode, result.StatusCode, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"), "Content-Type не совпадает с ожидаемым")

			url, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Contains(t, string(url), "http://localhost:8080/")
			assert.NotEmpty(t, tt.wantUrlMap)
		})
	}
}

func TestRedirectHandler(t *testing.T) {
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
			wantStatusCode: http.StatusTemporaryRedirect,
		},
		{
			name:           "Invalid Method",
			method:         http.MethodPost,
			request:        "/abc123",
			wantStatusCode: http.StatusBadRequest,
			wantErr:        "Only GET requests are allowed!\n",
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
			urlMap["abc123"] = "https://example.com"
			if tt.wantStatusCode != http.StatusTemporaryRedirect {
				urlMap = make(map[string]string)
			}
			request := httptest.NewRequest(tt.method, tt.request, nil)
			w := httptest.NewRecorder()

			RedirectHandler(w, request)

			result := w.Result()

			assert.Equal(t, tt.wantStatusCode, result.StatusCode, "Код ответа не совпадает с ожидаемым")

			if tt.wantStatusCode == http.StatusTemporaryRedirect {
				assert.Equal(t, urlMap["abc123"], result.Header.Get("Location"), "Location не совпадает с ожидаемым")
			} else {
				assert.Equal(t, w.Body.String(), tt.wantErr, "Не совпадает ошибка с ожидаемой")
			}
		})
	}
}
