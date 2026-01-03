package models

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAPIResponse(t *testing.T) {
	t.Parallel()

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-ID", "test-request-123")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	// Make a request
	resp, err := http.Get(server.URL)
	require.NoError(t, err)

	elapsed := 100 * time.Millisecond
	apiResp := NewAPIResponse(resp, elapsed)

	assert.NotNil(t, apiResp)
	assert.Equal(t, resp, apiResp.HTTPResponse)
	assert.Equal(t, resp.Body, apiResp.Body)
	assert.Equal(t, resp.Header, apiResp.Headers)
	assert.Equal(t, http.StatusOK, apiResp.StatusCode)
	assert.Equal(t, "GET", apiResp.Method)
	assert.Equal(t, elapsed, apiResp.Elapsed)
	assert.Equal(t, "test-request-123", apiResp.RequestID)
	assert.False(t, apiResp.IsClosed)
}

func TestAPIResponse_Close(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	require.NoError(t, err)

	apiResp := NewAPIResponse(resp, 0)

	// Close the response
	err = apiResp.Close()
	assert.NoError(t, err)
	assert.True(t, apiResp.IsClosed)

	// Closing again should not error
	err = apiResp.Close()
	assert.NoError(t, err)
}

func TestAPIResponse_GetHeader(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Custom-Header", "custom-value")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	apiResp := NewAPIResponse(resp, 0)

	assert.Equal(t, "custom-value", apiResp.GetHeader("Custom-Header"))
	assert.Equal(t, "", apiResp.GetHeader("Nonexistent-Header"))
}

func TestAPIResponse_IsSuccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		statusCode int
		expected   bool
	}{
		{"200 OK", http.StatusOK, true},
		{"201 Created", http.StatusCreated, true},
		{"204 No Content", http.StatusNoContent, true},
		{"400 Bad Request", http.StatusBadRequest, false},
		{"500 Internal Server Error", http.StatusInternalServerError, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			resp, err := http.Get(server.URL)
			require.NoError(t, err)
			defer resp.Body.Close()

			apiResp := NewAPIResponse(resp, 0)
			assert.Equal(t, tt.expected, apiResp.IsSuccess())
		})
	}
}

func TestAPIResponse_IsClientError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		statusCode int
		expected   bool
	}{
		{"400 Bad Request", http.StatusBadRequest, true},
		{"404 Not Found", http.StatusNotFound, true},
		{"429 Too Many Requests", http.StatusTooManyRequests, true},
		{"200 OK", http.StatusOK, false},
		{"500 Internal Server Error", http.StatusInternalServerError, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			apiResp := &APIResponse{StatusCode: tt.statusCode}
			assert.Equal(t, tt.expected, apiResp.IsClientError())
		})
	}
}

func TestAPIResponse_IsServerError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		statusCode int
		expected   bool
	}{
		{"500 Internal Server Error", http.StatusInternalServerError, true},
		{"502 Bad Gateway", http.StatusBadGateway, true},
		{"503 Service Unavailable", http.StatusServiceUnavailable, true},
		{"200 OK", http.StatusOK, false},
		{"400 Bad Request", http.StatusBadRequest, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			apiResp := &APIResponse{StatusCode: tt.statusCode}
			assert.Equal(t, tt.expected, apiResp.IsServerError())
		})
	}
}

func TestAPIResponse_IsError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		statusCode int
		expected   bool
	}{
		{"400 Bad Request", http.StatusBadRequest, true},
		{"500 Internal Server Error", http.StatusInternalServerError, true},
		{"200 OK", http.StatusOK, false},
		{"201 Created", http.StatusCreated, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			apiResp := &APIResponse{StatusCode: tt.statusCode}
			assert.Equal(t, tt.expected, apiResp.IsError())
		})
	}
}

func TestNewStreamResponse(t *testing.T) {
	t.Parallel()

	apiResp := &APIResponse{
		Body:       io.NopCloser(strings.NewReader("test data")),
		StatusCode: http.StatusOK,
	}

	streamResp := NewStreamResponse(apiResp)

	assert.NotNil(t, streamResp)
	assert.Equal(t, apiResp, streamResp.APIResponse)
	assert.NotNil(t, streamResp.Reader)
	assert.NotNil(t, streamResp.Done)
}

func TestStreamResponse_Close(t *testing.T) {
	t.Parallel()

	apiResp := &APIResponse{
		Body:       io.NopCloser(strings.NewReader("test data")),
		StatusCode: http.StatusOK,
	}

	streamResp := NewStreamResponse(apiResp)

	// Close should close the Done channel
	err := streamResp.Close()
	assert.NoError(t, err)

	// Done channel should be closed
	select {
	case <-streamResp.Done:
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Error("Done channel was not closed")
	}

	assert.True(t, apiResp.IsClosed)
}

func TestAPIResponse_RequestIDFallback(t *testing.T) {
	t.Parallel()

	t.Run("from X-Request-ID", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Request-ID", "x-request-123")
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		resp, err := http.Get(server.URL)
		require.NoError(t, err)
		defer resp.Body.Close()

		apiResp := NewAPIResponse(resp, 0)
		assert.Equal(t, "x-request-123", apiResp.RequestID)
	})

	t.Run("from Request-ID fallback", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Request-ID", "request-456")
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		resp, err := http.Get(server.URL)
		require.NoError(t, err)
		defer resp.Body.Close()

		apiResp := NewAPIResponse(resp, 0)
		assert.Equal(t, "request-456", apiResp.RequestID)
	})

	t.Run("no request ID", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		resp, err := http.Get(server.URL)
		require.NoError(t, err)
		defer resp.Body.Close()

		apiResp := NewAPIResponse(resp, 0)
		assert.Equal(t, "", apiResp.RequestID)
	})
}
