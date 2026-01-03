package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBaseClient(t *testing.T) {
	t.Parallel()

	t.Run("with valid config", func(t *testing.T) {
		t.Parallel()

		config := &Config{
			APIKey: "test-key.test-secret",
		}

		client, err := NewBaseClient(config)
		require.NoError(t, err)
		require.NotNil(t, client)

		assert.NotNil(t, client.httpClient)
		assert.NotNil(t, client.tokenGenerator)
		assert.NotNil(t, client.logger)
	})

	t.Run("without API key", func(t *testing.T) {
		t.Parallel()

		config := &Config{}

		client, err := NewBaseClient(config)
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "API key is required")
	})

	t.Run("with defaults", func(t *testing.T) {
		t.Parallel()

		config := &Config{
			APIKey: "test-key.test-secret",
		}

		client, err := NewBaseClient(config)
		require.NoError(t, err)

		assert.NotEmpty(t, client.config.BaseURL)
		assert.NotZero(t, client.config.Timeout)
		assert.NotZero(t, client.config.MaxRetries)
	})

	t.Run("with custom values", func(t *testing.T) {
		t.Parallel()

		config := &Config{
			APIKey:     "test-key.test-secret",
			BaseURL:    "https://custom.api.com",
			Timeout:    30 * time.Second,
			MaxRetries: 5,
		}

		client, err := NewBaseClient(config)
		require.NoError(t, err)

		assert.Equal(t, "https://custom.api.com", client.config.BaseURL)
		assert.Equal(t, 30*time.Second, client.config.Timeout)
		assert.Equal(t, 5, client.config.MaxRetries)
	})

	t.Run("with token cache disabled", func(t *testing.T) {
		t.Parallel()

		config := &Config{
			APIKey:            "test-key.test-secret",
			DisableTokenCache: true,
		}

		client, err := NewBaseClient(config)
		require.NoError(t, err)
		assert.NotNil(t, client.tokenGenerator)
	})
}

func TestBaseClient_Get(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer")

		// Check query parameters
		assert.Equal(t, "value1", r.URL.Query().Get("param1"))
		assert.Equal(t, "value2", r.URL.Query().Get("param2"))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	config := &Config{
		APIKey:  "test-key.test-secret",
		BaseURL: server.URL,
	}

	client, err := NewBaseClient(config)
	require.NoError(t, err)
	defer client.Close()

	query := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}

	resp, err := client.Get(context.Background(), "/test", query)
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer resp.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestBaseClient_Post(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer")
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Read and verify body
		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "test-value", body["key"])

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	config := &Config{
		APIKey:  "test-key.test-secret",
		BaseURL: server.URL,
	}

	client, err := NewBaseClient(config)
	require.NoError(t, err)
	defer client.Close()

	body := map[string]string{
		"key": "test-value",
	}

	resp, err := client.Post(context.Background(), "/test", body)
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer resp.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestBaseClient_Put(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := &Config{
		APIKey:  "test-key.test-secret",
		BaseURL: server.URL,
	}

	client, err := NewBaseClient(config)
	require.NoError(t, err)
	defer client.Close()

	resp, err := client.Put(context.Background(), "/test", map[string]string{"key": "value"})
	require.NoError(t, err)
	defer resp.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestBaseClient_Delete(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	config := &Config{
		APIKey:  "test-key.test-secret",
		BaseURL: server.URL,
	}

	client, err := NewBaseClient(config)
	require.NoError(t, err)
	defer client.Close()

	resp, err := client.Delete(context.Background(), "/test")
	require.NoError(t, err)
	defer resp.Close()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestBaseClient_ParseJSON(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "hello",
			"status":  "success",
		})
	}))
	defer server.Close()

	config := &Config{
		APIKey:  "test-key.test-secret",
		BaseURL: server.URL,
	}

	client, err := NewBaseClient(config)
	require.NoError(t, err)
	defer client.Close()

	resp, err := client.Get(context.Background(), "/test", nil)
	require.NoError(t, err)

	var result map[string]string
	err = client.ParseJSON(resp, &result)
	require.NoError(t, err)

	assert.Equal(t, "hello", result["message"])
	assert.Equal(t, "success", result["status"])
}

func TestBaseClient_ErrorHandling(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		statusCode     int
		errorResponse  map[string]interface{}
		expectedErrMsg string
	}{
		{
			name:       "400 Bad Request",
			statusCode: http.StatusBadRequest,
			errorResponse: map[string]interface{}{
				"error": map[string]string{
					"message": "Invalid request",
					"code":    "invalid_request",
				},
			},
			expectedErrMsg: "Invalid request",
		},
		{
			name:       "401 Unauthorized",
			statusCode: http.StatusUnauthorized,
			errorResponse: map[string]interface{}{
				"error": map[string]string{
					"message": "Authentication failed",
				},
			},
			expectedErrMsg: "Authentication failed",
		},
		{
			name:       "429 Too Many Requests",
			statusCode: http.StatusTooManyRequests,
			errorResponse: map[string]interface{}{
				"message": "Rate limit exceeded",
			},
			expectedErrMsg: "Rate limit exceeded",
		},
		{
			name:       "500 Internal Server Error",
			statusCode: http.StatusInternalServerError,
			errorResponse: map[string]interface{}{
				"error": map[string]string{
					"message": "Internal error",
				},
			},
			expectedErrMsg: "Internal error",
		},
		{
			name:       "503 Service Unavailable",
			statusCode: http.StatusServiceUnavailable,
			errorResponse: map[string]interface{}{
				"message": "Service unavailable",
			},
			expectedErrMsg: "Service unavailable",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.errorResponse)
			}))
			defer server.Close()

			config := &Config{
				APIKey:  "test-key.test-secret",
				BaseURL: server.URL,
			}

			client, err := NewBaseClient(config)
			require.NoError(t, err)
			defer client.Close()

			_, err = client.Get(context.Background(), "/test", nil)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErrMsg)
		})
	}
}

func TestBaseClient_Authentication(t *testing.T) {
	t.Parallel()

	t.Run("with JWT token", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			assert.True(t, len(authHeader) > 7) // "Bearer " + token
			assert.Contains(t, authHeader, "Bearer ")

			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		config := &Config{
			APIKey:            "test-key.test-secret",
			BaseURL:           server.URL,
			DisableTokenCache: false,
		}

		client, err := NewBaseClient(config)
		require.NoError(t, err)
		defer client.Close()

		_, err = client.Get(context.Background(), "/test", nil)
		require.NoError(t, err)
	})

	t.Run("with raw API key", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			assert.Equal(t, "Bearer test-key.test-secret", authHeader)

			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		config := &Config{
			APIKey:            "test-key.test-secret",
			BaseURL:           server.URL,
			DisableTokenCache: true,
		}

		client, err := NewBaseClient(config)
		require.NoError(t, err)
		defer client.Close()

		_, err = client.Get(context.Background(), "/test", nil)
		require.NoError(t, err)
	})
}

func TestBaseClient_GetConfig(t *testing.T) {
	t.Parallel()

	config := &Config{
		APIKey: "test-key.test-secret",
	}

	client, err := NewBaseClient(config)
	require.NoError(t, err)

	retrievedConfig := client.GetConfig()
	assert.Equal(t, config, retrievedConfig)
}

func TestBaseClient_GetLogger(t *testing.T) {
	t.Parallel()

	config := &Config{
		APIKey: "test-key.test-secret",
	}

	client, err := NewBaseClient(config)
	require.NoError(t, err)

	logger := client.GetLogger()
	assert.NotNil(t, logger)
}

func TestBytesReader_Read(t *testing.T) {
	t.Parallel()

	data := []byte("hello world")
	reader := newBytesReader(data)

	// Read in chunks
	buf := make([]byte, 5)

	n, err := reader.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, "hello", string(buf))

	n, err = reader.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, " worl", string(buf))

	buf = make([]byte, 5)
	n, err = reader.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.Equal(t, "d", string(buf[:n]))

	// Should return EOF
	n, err = reader.Read(buf)
	assert.ErrorIs(t, err, io.EOF)
	assert.Equal(t, 0, n)
}
