package transport

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestDefaultRetryConfig(t *testing.T) {
	t.Parallel()

	config := DefaultRetryConfig()

	if config == nil {
		t.Fatal("DefaultRetryConfig returned nil")
	}

	if config.MaxRetries <= 0 {
		t.Error("MaxRetries should be positive")
	}

	if config.InitialBackoff <= 0 {
		t.Error("InitialBackoff should be positive")
	}

	if config.MaxBackoff <= config.InitialBackoff {
		t.Error("MaxBackoff should be greater than InitialBackoff")
	}

	if config.BackoffMultiplier <= 1.0 {
		t.Error("BackoffMultiplier should be greater than 1.0")
	}

	if len(config.RetryableStatusCodes) == 0 {
		t.Error("RetryableStatusCodes should not be empty")
	}
}

func TestNewRetryableHTTPClient(t *testing.T) {
	t.Parallel()

	httpClient := NewHTTPClient(nil)

	tests := []struct {
		name   string
		config *RetryConfig
	}{
		{
			name:   "nil config uses defaults",
			config: nil,
		},
		{
			name:   "custom config",
			config: DefaultRetryConfig(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := NewRetryableHTTPClient(httpClient, tt.config)

			if client == nil {
				t.Fatal("NewRetryableHTTPClient returned nil")
			}

			if client.client == nil {
				t.Fatal("RetryableHTTPClient.client is nil")
			}

			if client.config == nil {
				t.Fatal("RetryableHTTPClient.config is nil")
			}
		})
	}
}

func TestRetryableHTTPClient_SuccessfulRequest(t *testing.T) {
	t.Parallel()

	// Create test server that always succeeds
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	httpClient := NewHTTPClient(&HTTPClientConfig{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})
	retryClient := NewRetryableHTTPClient(httpClient, nil)

	ctx := context.Background()
	req, err := httpClient.NewRequest(ctx, http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}

	resp, err := retryClient.DoWithRetry(ctx, req)
	if err != nil {
		t.Fatalf("DoWithRetry failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestRetryableHTTPClient_RetryOn500(t *testing.T) {
	t.Parallel()

	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "server error"}`))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok"}`))
		}
	}))
	defer server.Close()

	httpClient := NewHTTPClient(&HTTPClientConfig{
		BaseURL: server.URL,
		Timeout: 10 * time.Second,
	})

	config := DefaultRetryConfig()
	config.InitialBackoff = 10 * time.Millisecond
	config.MaxBackoff = 100 * time.Millisecond
	retryClient := NewRetryableHTTPClient(httpClient, config)

	ctx := context.Background()
	req, err := httpClient.NewRequest(ctx, http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}

	resp, err := retryClient.DoWithRetry(ctx, req)
	if err != nil {
		t.Fatalf("DoWithRetry failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestRetryableHTTPClient_ExhaustRetries(t *testing.T) {
	t.Parallel()

	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"error": "service unavailable"}`))
	}))
	defer server.Close()

	httpClient := NewHTTPClient(&HTTPClientConfig{
		BaseURL: server.URL,
		Timeout: 10 * time.Second,
	})

	config := DefaultRetryConfig()
	config.MaxRetries = 2
	config.InitialBackoff = 10 * time.Millisecond
	config.MaxBackoff = 50 * time.Millisecond
	retryClient := NewRetryableHTTPClient(httpClient, config)

	ctx := context.Background()
	req, err := httpClient.NewRequest(ctx, http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}

	resp, err := retryClient.DoWithRetry(ctx, req)
	if err != nil {
		t.Fatalf("DoWithRetry failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusServiceUnavailable)
	}

	// Should have attempted initial + MaxRetries (1 + 2 = 3)
	expectedAttempts := 1 + config.MaxRetries
	if attempts != expectedAttempts {
		t.Errorf("Expected %d attempts, got %d", expectedAttempts, attempts)
	}
}

func TestRetryableHTTPClient_NoRetryOnNonRetryableStatusCode(t *testing.T) {
	t.Parallel()

	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "bad request"}`))
	}))
	defer server.Close()

	httpClient := NewHTTPClient(&HTTPClientConfig{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})
	retryClient := NewRetryableHTTPClient(httpClient, nil)

	ctx := context.Background()
	req, err := httpClient.NewRequest(ctx, http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}

	resp, err := retryClient.DoWithRetry(ctx, req)
	if err != nil {
		t.Fatalf("DoWithRetry failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}

	// Should only attempt once (no retry on 400)
	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}

func TestRetryableHTTPClient_NoRetryOnPOST(t *testing.T) {
	t.Parallel()

	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "server error"}`))
	}))
	defer server.Close()

	httpClient := NewHTTPClient(&HTTPClientConfig{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})
	retryClient := NewRetryableHTTPClient(httpClient, nil)

	ctx := context.Background()
	body := strings.NewReader(`{"data": "test"}`)
	req, err := httpClient.NewRequest(ctx, http.MethodPost, "/test", body)
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}

	resp, err := retryClient.DoWithRetry(ctx, req)
	if err != nil {
		t.Fatalf("DoWithRetry failed: %v", err)
	}
	defer resp.Body.Close()

	// Should only attempt once (POST is not idempotent)
	if attempts != 1 {
		t.Errorf("Expected 1 attempt (no retry for POST), got %d", attempts)
	}
}

func TestRetryableHTTPClient_RetryAfterHeader(t *testing.T) {
	t.Parallel()

	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			w.Header().Set("Retry-After", "1") // 1 second
			w.WriteHeader(http.StatusTooManyRequests)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	httpClient := NewHTTPClient(&HTTPClientConfig{
		BaseURL: server.URL,
		Timeout: 10 * time.Second,
	})
	retryClient := NewRetryableHTTPClient(httpClient, nil)

	ctx := context.Background()
	req, err := httpClient.NewRequest(ctx, http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}

	start := time.Now()
	resp, err := retryClient.DoWithRetry(ctx, req)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("DoWithRetry failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	// Should have waited at least 1 second due to Retry-After header
	if elapsed < 1*time.Second {
		t.Errorf("Expected to wait at least 1 second, but only waited %v", elapsed)
	}

	if attempts != 2 {
		t.Errorf("Expected 2 attempts, got %d", attempts)
	}
}

func TestRetryableHTTPClient_ContextCancellation(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	httpClient := NewHTTPClient(&HTTPClientConfig{
		BaseURL: server.URL,
		Timeout: 10 * time.Second,
	})

	config := DefaultRetryConfig()
	config.InitialBackoff = 1 * time.Second
	retryClient := NewRetryableHTTPClient(httpClient, config)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	req, err := httpClient.NewRequest(ctx, http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}

	_, err = retryClient.DoWithRetry(ctx, req)
	if err == nil {
		t.Error("Expected context cancellation error, got nil")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Expected context.DeadlineExceeded, got %v", err)
	}
}

func TestRetryableHTTPClient_IsIdempotentMethod(t *testing.T) {
	t.Parallel()

	client := NewRetryableHTTPClient(NewHTTPClient(nil), nil)

	tests := []struct {
		method     string
		idempotent bool
	}{
		{http.MethodGet, true},
		{http.MethodHead, true},
		{http.MethodOptions, true},
		{http.MethodTrace, true},
		{http.MethodPut, true},
		{http.MethodDelete, true},
		{http.MethodPost, false},
		{http.MethodPatch, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.method, func(t *testing.T) {
			t.Parallel()

			got := client.isIdempotentMethod(tt.method)
			if got != tt.idempotent {
				t.Errorf("isIdempotentMethod(%s) = %v, want %v", tt.method, got, tt.idempotent)
			}
		})
	}
}

func TestRetryableHTTPClient_IsRetryableStatusCode(t *testing.T) {
	t.Parallel()

	client := NewRetryableHTTPClient(NewHTTPClient(nil), nil)

	tests := []struct {
		statusCode int
		retryable  bool
	}{
		{429, true},  // Too Many Requests
		{500, true},  // Internal Server Error
		{502, true},  // Bad Gateway
		{503, true},  // Service Unavailable
		{504, true},  // Gateway Timeout
		{200, false}, // OK
		{400, false}, // Bad Request
		{401, false}, // Unauthorized
		{404, false}, // Not Found
	}

	for _, tt := range tests {
		tt := tt
		t.Run(string(rune(tt.statusCode)), func(t *testing.T) {
			t.Parallel()

			got := client.isRetryableStatusCode(tt.statusCode)
			if got != tt.retryable {
				t.Errorf("isRetryableStatusCode(%d) = %v, want %v", tt.statusCode, got, tt.retryable)
			}
		})
	}
}

func TestRetryableHTTPClient_CalculateBackoff(t *testing.T) {
	t.Parallel()

	config := &RetryConfig{
		InitialBackoff:    100 * time.Millisecond,
		MaxBackoff:        5 * time.Second,
		BackoffMultiplier: 2.0,
		EnableJitter:      false, // Disable for predictable testing
	}

	client := NewRetryableHTTPClient(NewHTTPClient(nil), config)

	tests := []struct {
		name       string
		attempt    int
		retryAfter time.Duration
		wantMin    time.Duration
		wantMax    time.Duration
	}{
		{
			name:       "first retry",
			attempt:    0,
			retryAfter: 0,
			wantMin:    100 * time.Millisecond,
			wantMax:    100 * time.Millisecond,
		},
		{
			name:       "second retry",
			attempt:    1,
			retryAfter: 0,
			wantMin:    200 * time.Millisecond,
			wantMax:    200 * time.Millisecond,
		},
		{
			name:       "third retry",
			attempt:    2,
			retryAfter: 0,
			wantMin:    400 * time.Millisecond,
			wantMax:    400 * time.Millisecond,
		},
		{
			name:       "with retry-after header",
			attempt:    0,
			retryAfter: 2 * time.Second,
			wantMin:    2 * time.Second,
			wantMax:    2 * time.Second,
		},
		{
			name:       "exceeds max backoff",
			attempt:    10,
			retryAfter: 0,
			wantMin:    5 * time.Second,
			wantMax:    5 * time.Second,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			backoff := client.calculateBackoff(tt.attempt, tt.retryAfter)

			if backoff < tt.wantMin || backoff > tt.wantMax {
				t.Errorf("calculateBackoff(%d, %v) = %v, want between %v and %v",
					tt.attempt, tt.retryAfter, backoff, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestRetryableHTTPClient_ParseRetryAfter(t *testing.T) {
	t.Parallel()

	client := NewRetryableHTTPClient(NewHTTPClient(nil), nil)

	tests := []struct {
		name        string
		headerValue string
		wantMin     time.Duration
		wantMax     time.Duration
	}{
		{
			name:        "seconds format",
			headerValue: "60",
			wantMin:     60 * time.Second,
			wantMax:     60 * time.Second,
		},
		{
			name:        "zero seconds",
			headerValue: "0",
			wantMin:     0,
			wantMax:     0,
		},
		{
			name:        "invalid format",
			headerValue: "invalid",
			wantMin:     0,
			wantMax:     0,
		},
		{
			name:        "too large value",
			headerValue: "9999",
			wantMin:     0,
			wantMax:     0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resp := &http.Response{
				Header: http.Header{
					"Retry-After": []string{tt.headerValue},
				},
			}

			duration := client.parseRetryAfter(resp)

			if duration < tt.wantMin || duration > tt.wantMax {
				t.Errorf("parseRetryAfter(%s) = %v, want between %v and %v",
					tt.headerValue, duration, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestRetryableHTTPClient_Close(t *testing.T) {
	t.Parallel()

	httpClient := NewHTTPClient(nil)
	retryClient := NewRetryableHTTPClient(httpClient, nil)

	// Should not panic
	retryClient.Close()
}

func TestRetryableHTTPClient_GetClient(t *testing.T) {
	t.Parallel()

	httpClient := NewHTTPClient(nil)
	retryClient := NewRetryableHTTPClient(httpClient, nil)

	got := retryClient.GetClient()
	if got != httpClient {
		t.Error("GetClient() did not return the correct client")
	}
}

func TestRetryableHTTPClient_CloneRequest(t *testing.T) {
	t.Parallel()

	client := NewRetryableHTTPClient(NewHTTPClient(nil), nil)
	ctx := context.Background()

	t.Run("nil body", func(t *testing.T) {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", nil)
		cloned, err := client.cloneRequest(ctx, req)

		if err != nil {
			t.Fatalf("cloneRequest failed: %v", err)
		}

		if cloned == nil {
			t.Fatal("cloneRequest returned nil")
		}

		if cloned.Method != req.Method {
			t.Errorf("Method = %s, want %s", cloned.Method, req.Method)
		}
	})

	t.Run("with GetBody", func(t *testing.T) {
		bodyContent := "test body"
		req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "http://example.com", strings.NewReader(bodyContent))

		cloned, err := client.cloneRequest(ctx, req)

		if err != nil {
			t.Fatalf("cloneRequest failed: %v", err)
		}

		if cloned == nil {
			t.Fatal("cloneRequest returned nil")
		}

		if cloned.Body == nil {
			t.Fatal("cloned request has nil body")
		}

		// Read the cloned body
		body, _ := io.ReadAll(cloned.Body)
		if string(body) != bodyContent {
			t.Errorf("Body = %s, want %s", string(body), bodyContent)
		}
	})
}
