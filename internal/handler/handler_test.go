package handler_test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/ChebuRashkaRF/urlshortener/internal/middleware"
	"github.com/ChebuRashkaRF/urlshortener/internal/storage/mock"
	"github.com/ChebuRashkaRF/urlshortener/internal/util"
	"github.com/golang/mock/gomock"
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

func setupTestEnvironment(t *testing.T, handlerFunc http.HandlerFunc, method, path string, middlewareList ...func(http.Handler) http.Handler) (*httptest.Server, string, *mock.MockStorage) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock.NewMockStorage(ctrl)

	r := setupTestRouter(handlerFunc, method, path, middlewareList...)
	ts := httptest.NewServer(r)

	// Извлечение порта из URL
	parts := strings.Split(ts.URL, ":")
	port := parts[len(parts)-1]

	config.Cnf = &config.Config{
		ServerAddress: ":" + port,
		BaseURL:       ts.URL,
	}

	handler.URLStore = mockStore

	t.Cleanup(func() {
		ts.Close()
	})

	return ts, ts.URL, mockStore
}

func TestShortenURLHandler(t *testing.T) {
	ts, baseURL, mockStore := setupTestEnvironment(t, handler.ShortenURLHandler, http.MethodPost, "/")

	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		reqBody string
		method  string
		headers map[string]string
		want    want
		wantErr string
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
			if tt.want.statusCode == http.StatusCreated {
				mockStore.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			}

			resp, body := testRequest(t, ts, tt.method, "/", strings.NewReader(tt.reqBody), tt.headers)
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				assert.Equal(t, tt.want.statusCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
				assert.Equal(t, string(body), tt.wantErr, "Не совпадает ошибка с ожидаемой")
				return
			}

			assert.Equal(t, tt.want.statusCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"), "Content-Type не совпадает с ожидаемым")

			assert.Contains(t, string(body), baseURL)
		})
	}
}

func TestShortenURLJSONHandler(t *testing.T) {
	ts, baseURL, mockStore := setupTestEnvironment(t, handler.ShortenURLJSONHandler, http.MethodPost, "/api/shorten")

	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		reqURL  string
		method  string
		headers map[string]string
		want    want
		wantErr string
	}{
		{
			name:   "POST request method",
			reqURL: "https://example.com",
			method: http.MethodPost,
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
			},
		},
		{
			name:   "Invalid reqBody",
			reqURL: "",
			method: http.MethodPost,
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
			name:   "Invalid URL",
			reqURL: "yandex.ru",
			method: http.MethodPost,
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
			name:   "Invalid Method",
			reqURL: "https://example.com",
			method: http.MethodGet,
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
			if tt.want.statusCode == http.StatusCreated {
				mockStore.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			}

			reqBody := fmt.Sprintf(`{"url": "%s"}`, tt.reqURL)
			resp, body := testRequest(t, ts, tt.method, "/api/shorten", strings.NewReader(reqBody), tt.headers)
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				assert.Equal(t, tt.want.statusCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
				assert.Equal(t, string(body), tt.wantErr, "Не совпадает ошибка с ожидаемой")
				return
			}

			assert.Equal(t, tt.want.statusCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"), "Content-Type не совпадает с ожидаемым")

			id := util.GenerateShortID(tt.reqURL)
			successBody := fmt.Sprintf(`{"result": "%s/%s"}`, baseURL, id)
			assert.JSONEq(t, successBody, string(body), "ответ не совпадает с ожидаемым")

		})
	}
}

func TestGzipCompressionShortenURLHandler(t *testing.T) {
	ts, baseURL, mockStore := setupTestEnvironment(t, handler.ShortenURLHandler, http.MethodPost, "/", middleware.GzipMiddleware)

	requestBody := "https://example.com"

	t.Run("sends_gzip", func(t *testing.T) {
		headers := map[string]string{
			"Content-Encoding": "gzip",
		}
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		mockStore.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		resp, body := testRequest(t, ts, http.MethodPost, "/", buf, headers)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		assert.Contains(t, string(body), baseURL, "ответ не совпадает с ожидаемым")
	})

	t.Run("accepts_gzip", func(t *testing.T) {
		headers := map[string]string{
			"Accept-Encoding": "gzip",
		}

		buf := bytes.NewBufferString(requestBody)

		mockStore.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		resp, body := testRequest(t, ts, http.MethodPost, "/", buf, headers)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		defer resp.Body.Close()

		assert.Contains(t, string(body), baseURL, "ответ не совпадает с ожидаемым")
	})
}

func TestGzipCompressionShortenURLJSONHandler(t *testing.T) {
	ts, baseURL, mockStore := setupTestEnvironment(t, handler.ShortenURLJSONHandler, http.MethodPost, "/api/shorten", middleware.GzipMiddleware)

	requestURL := "https://example.com"
	requestBody := `{"url": "https://example.com"}`

	t.Run("sends_gzip", func(t *testing.T) {
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

		mockStore.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		resp, body := testRequest(t, ts, http.MethodPost, "/api/shorten", buf, headers)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		defer resp.Body.Close()

		id := util.GenerateShortID(requestURL)
		successBody := fmt.Sprintf(`{"result": "%s/%s"}`, baseURL, id)
		assert.JSONEq(t, successBody, string(body), "ответ не совпадает с ожидаемым")
	})

	t.Run("accepts_gzip", func(t *testing.T) {
		headers := map[string]string{
			"Content-Type":    "application/json",
			"Accept-Encoding": "gzip",
		}

		buf := bytes.NewBufferString(requestBody)

		mockStore.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		resp, body := testRequest(t, ts, http.MethodPost, "/api/shorten", buf, headers)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		defer resp.Body.Close()

		id := util.GenerateShortID(requestURL)
		successBody := fmt.Sprintf(`{"result": "%s/%s"}`, baseURL, id)
		assert.JSONEq(t, successBody, string(body), "ответ не совпадает с ожидаемым")
	})
}

func TestRedirectHandler(t *testing.T) {
	ts, _, mockStore := setupTestEnvironment(t, handler.RedirectHandler, http.MethodGet, "/{id}")

	tests := []struct {
		name           string
		method         string
		headers        map[string]string
		request        string
		wantStatusCode int
		wantErr        string
		setupMock      func()
	}{
		{
			name:           "GET request method",
			method:         http.MethodGet,
			headers:        nil,
			request:        "/abc123",
			wantStatusCode: http.StatusOK,
			setupMock: func() {
				mockStore.EXPECT().Get("abc123").Return("https://example.com", true).Times(1)
			},
		},
		{
			name:           "Invalid Method",
			method:         http.MethodPost,
			headers:        nil,
			request:        "/abc123",
			wantStatusCode: http.StatusMethodNotAllowed,
			setupMock:      func() {},
		},
		{
			name:           "URL Not Found",
			method:         http.MethodGet,
			headers:        nil,
			request:        "/invalidid",
			wantStatusCode: http.StatusBadRequest,
			wantErr:        "URL not found\n",
			setupMock: func() {
				mockStore.EXPECT().Get("invalidid").Return("", false).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			resp, body := testRequest(t, ts, tt.method, tt.request, strings.NewReader(""), tt.headers)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatusCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")

			if tt.wantErr != "" {
				assert.Equal(t, string(body), tt.wantErr, "Не совпадает ошибка с ожидаемой")
			}
		})
	}
}
