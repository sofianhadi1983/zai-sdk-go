package transport

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sofianhadi1983/zai-sdk-go/internal/constants"
)

func TestDefaultHTTPClientConfig(t *testing.T) {
	t.Parallel()

	config := DefaultHTTPClientConfig()

	if config == nil {
		t.Fatal("DefaultHTTPClientConfig returned nil")
	}

	if config.BaseURL != constants.ZaiBaseURL {
		t.Errorf("BaseURL = %s, want %s", config.BaseURL, constants.ZaiBaseURL)
	}

	if config.Timeout != constants.DefaultTimeout {
		t.Errorf("Timeout = %v, want %v", config.Timeout, constants.DefaultTimeout)
	}

	if config.ConnectTimeout != constants.DefaultConnectTimeout {
		t.Errorf("ConnectTimeout = %v, want %v", config.ConnectTimeout, constants.DefaultConnectTimeout)
	}
}

func TestNewHTTPClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		config *HTTPClientConfig
	}{
		{
			name:   "nil config uses defaults",
			config: nil,
		},
		{
			name:   "custom config",
			config: DefaultHTTPClientConfig(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := NewHTTPClient(tt.config)

			if client == nil {
				t.Fatal("NewHTTPClient returned nil")
			}

			if client.client == nil {
				t.Fatal("HTTPClient.client is nil")
			}

			if client.config == nil {
				t.Fatal("HTTPClient.config is nil")
			}
		})
	}
}

func TestHTTPClient_BuildURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		baseURL string
		path    string
		want    string
		wantErr bool
	}{
		{
			name:    "relative path",
			baseURL: "https://api.example.com/v1",
			path:    "chat/completions",
			want:    "https://api.example.com/v1/chat/completions",
			wantErr: false,
		},
		{
			name:    "relative path with leading slash",
			baseURL: "https://api.example.com/v1",
			path:    "/chat/completions",
			want:    "https://api.example.com/v1/chat/completions",
			wantErr: false,
		},
		{
			name:    "base URL without trailing slash",
			baseURL: "https://api.example.com/v1",
			path:    "embeddings",
			want:    "https://api.example.com/v1/embeddings",
			wantErr: false,
		},
		{
			name:    "base URL with trailing slash",
			baseURL: "https://api.example.com/v1/",
			path:    "images",
			want:    "https://api.example.com/v1/images",
			wantErr: false,
		},
		{
			name:    "absolute URL path",
			baseURL: "https://api.example.com/v1",
			path:    "https://api.other.com/test",
			want:    "https://api.other.com/test",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			config := &HTTPClientConfig{BaseURL: tt.baseURL}
			client := NewHTTPClient(config)

			got, err := client.buildURL(tt.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("buildURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("buildURL() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestHTTPClient_NewRequest(t *testing.T) {
	t.Parallel()

	client := NewHTTPClient(nil)
	ctx := context.Background()

	tests := []struct {
		name    string
		method  string
		path    string
		body    io.Reader
		wantErr bool
	}{
		{
			name:    "GET request without body",
			method:  http.MethodGet,
			path:    "/test",
			body:    nil,
			wantErr: false,
		},
		{
			name:    "POST request with body",
			method:  http.MethodPost,
			path:    "/create",
			body:    bytes.NewReader([]byte(`{"test": "data"}`)),
			wantErr: false,
		},
		{
			name:    "PUT request",
			method:  http.MethodPut,
			path:    "/update",
			body:    strings.NewReader("test"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req, err := client.NewRequest(ctx, tt.method, tt.path, tt.body)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if req == nil {
					t.Fatal("NewRequest() returned nil request")
				}

				if req.Method != tt.method {
					t.Errorf("Request.Method = %s, want %s", req.Method, tt.method)
				}

				// Check default headers are set
				if req.Header.Get(constants.HeaderContentType) == "" {
					t.Error("Content-Type header not set")
				}

				if req.Header.Get(constants.HeaderAccept) == "" {
					t.Error("Accept header not set")
				}

				if req.Header.Get(constants.HeaderUserAgent) == "" {
					t.Error("User-Agent header not set")
				}
			}
		})
	}
}

func TestHTTPClient_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		wantStatusCode int
		wantErr        bool
	}{
		{
			name: "successful request",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status": "ok"}`))
			},
			wantStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "server error",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "internal error"}`))
			},
			wantStatusCode: http.StatusInternalServerError,
			wantErr:        false,
		},
		{
			name: "bad request",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error": "bad request"}`))
			},
			wantStatusCode: http.StatusBadRequest,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create test server
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			// Create client
			config := &HTTPClientConfig{
				BaseURL: server.URL,
				Timeout: 5 * time.Second,
			}
			client := NewHTTPClient(config)

			// Create request
			ctx := context.Background()
			req, err := client.NewRequest(ctx, http.MethodGet, "/test", nil)
			if err != nil {
				t.Fatalf("NewRequest() failed: %v", err)
			}

			// Execute request
			resp, err := client.Do(ctx, req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				defer resp.Body.Close()

				if resp.StatusCode != tt.wantStatusCode {
					t.Errorf("StatusCode = %d, want %d", resp.StatusCode, tt.wantStatusCode)
				}
			}
		})
	}
}

func TestHTTPClient_RequestMiddleware(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if custom header was added by middleware
		if r.Header.Get("X-Custom-Header") != "test-value" {
			t.Error("Custom header not set by middleware")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := &HTTPClientConfig{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	}
	client := NewHTTPClient(config)

	// Add request middleware
	client.AddRequestMiddleware(func(req *http.Request) error {
		req.Header.Set("X-Custom-Header", "test-value")
		return nil
	})

	// Create and execute request
	ctx := context.Background()
	req, err := client.NewRequest(ctx, http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("NewRequest() failed: %v", err)
	}

	resp, err := client.Do(ctx, req)
	if err != nil {
		t.Fatalf("Do() failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestHTTPClient_ResponseMiddleware(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Response-Header", "server-value")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := &HTTPClientConfig{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	}
	client := NewHTTPClient(config)

	middlewareCalled := false

	// Add response middleware
	client.AddResponseMiddleware(func(resp *http.Response) error {
		middlewareCalled = true
		if resp.Header.Get("X-Response-Header") != "server-value" {
			t.Error("Response header not found")
		}
		return nil
	})

	// Create and execute request
	ctx := context.Background()
	req, err := client.NewRequest(ctx, http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("NewRequest() failed: %v", err)
	}

	resp, err := client.Do(ctx, req)
	if err != nil {
		t.Fatalf("Do() failed: %v", err)
	}
	defer resp.Body.Close()

	if !middlewareCalled {
		t.Error("Response middleware was not called")
	}
}

func TestHTTPClient_ContextCancellation(t *testing.T) {
	t.Parallel()

	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := &HTTPClientConfig{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	}
	client := NewHTTPClient(config)

	// Create context that will be cancelled
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	req, err := client.NewRequest(ctx, http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("NewRequest() failed: %v", err)
	}

	_, err = client.Do(ctx, req)
	if err == nil {
		t.Error("Expected error due to context cancellation, got nil")
	}
}

func TestHTTPClient_Close(t *testing.T) {
	t.Parallel()

	client := NewHTTPClient(nil)

	// Should not panic
	client.Close()
}

func TestHTTPClient_SetDefaultHeaders(t *testing.T) {
	t.Parallel()

	client := NewHTTPClient(nil)
	ctx := context.Background()

	req, err := client.NewRequest(ctx, http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("NewRequest() failed: %v", err)
	}

	// Check all default headers are set
	headers := map[string]string{
		constants.HeaderContentType:   constants.ContentTypeJSON,
		constants.HeaderAccept:         constants.AcceptJSON,
		constants.HeaderUserAgent:      constants.GetUserAgent(),
		constants.HeaderSourceChannel:  constants.DefaultSourceChannel,
	}

	for key, want := range headers {
		got := req.Header.Get(key)
		if got != want {
			t.Errorf("Header %s = %s, want %s", key, got, want)
		}
	}
}
