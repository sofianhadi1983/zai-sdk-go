package constants

import (
	"testing"
	"time"
)

func TestSDKConstants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{"SDK Version", SDKVersion, "0.1.0"},
		{"SDK Title", SDKTitle, "Z.ai Go SDK"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.got != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestBaseURLs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{"Zai Base URL", ZaiBaseURL, "https://api.z.ai/api/paas/v4"},
		{"Zhipu Base URL", ZhipuBaseURL, "https://open.bigmodel.cn/api/paas/v4"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.got != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestHTTPClientDefaults(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"Default Timeout", DefaultTimeout, 300 * time.Second},
		{"Default Connect Timeout", DefaultConnectTimeout, 8 * time.Second},
		{"Default Max Retries", DefaultMaxRetries, 3},
		{"Default Max Connections", DefaultMaxConnections, 50},
		{"Default Max Idle Conns", DefaultMaxIdleConns, 10},
		{"Default Max Idle Conns Per Host", DefaultMaxIdleConnsPerHost, 10},
		{"Default Idle Conn Timeout", DefaultIdleConnTimeout, 90 * time.Second},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.got != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestRetryConfiguration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"Initial Retry Delay", InitialRetryDelay, 500 * time.Millisecond},
		{"Max Retry Delay", MaxRetryDelay, 8 * time.Second},
		{"Retry Backoff Multiplier", RetryBackoffMultiplier, 2.0},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.got != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestHTTPHeaders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{"Authorization Header", HeaderAuthorization, "Authorization"},
		{"Content-Type Header", HeaderContentType, "Content-Type"},
		{"Accept Header", HeaderAccept, "Accept"},
		{"Accept-Language Header", HeaderAcceptLanguage, "Accept-Language"},
		{"User-Agent Header", HeaderUserAgent, "User-Agent"},
		{"Source Channel Header", HeaderSourceChannel, "x-source-channel"},
		{"Raw Response Header", HeaderRawResponse, "X-Stainless-Raw-Response"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.got != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestHeaderValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{"Default Source Channel", DefaultSourceChannel, "go-sdk"},
		{"Content-Type JSON", ContentTypeJSON, "application/json"},
		{"Content-Type Form URLEncoded", ContentTypeFormURLEncoded, "application/x-www-form-urlencoded"},
		{"Content-Type Multipart", ContentTypeMultipartFormData, "multipart/form-data"},
		{"Accept JSON", AcceptJSON, "application/json"},
		{"Accept-Language English", AcceptLanguageEnglish, "en-US,en"},
		{"Accept-Language Chinese", AcceptLanguageChinese, "zh-CN,zh"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.got != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestGetUserAgent(t *testing.T) {
	t.Parallel()

	userAgent := GetUserAgent()
	expected := "Z.ai Go SDK/0.1.0"

	if userAgent != expected {
		t.Errorf("GetUserAgent() = %q, want %q", userAgent, expected)
	}

	// Verify it contains version
	if userAgent == "" {
		t.Error("GetUserAgent() returned empty string")
	}

	// Verify format
	if len(userAgent) == 0 {
		t.Error("GetUserAgent() should not be empty")
	}
}

func TestStatusCodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		got      int
		expected int
	}{
		{"Status Too Many Requests", StatusTooManyRequests, 429},
		{"Status Internal Server Error", StatusInternalServerError, 500},
		{"Status Bad Gateway", StatusBadGateway, 502},
		{"Status Service Unavailable", StatusServiceUnavailable, 503},
		{"Status Gateway Timeout", StatusGatewayTimeout, 504},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.got != tt.expected {
				t.Errorf("%s = %d, want %d", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestRetryableStatusCodes(t *testing.T) {
	t.Parallel()

	codes := RetryableStatusCodes()

	expectedCodes := []int{429, 500, 502, 503, 504}

	if len(codes) != len(expectedCodes) {
		t.Errorf("RetryableStatusCodes() returned %d codes, want %d", len(codes), len(expectedCodes))
	}

	// Verify each expected code is present
	for _, expected := range expectedCodes {
		found := false
		for _, code := range codes {
			if code == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("RetryableStatusCodes() missing status code %d", expected)
		}
	}

	// Verify all codes are within expected set
	codeMap := make(map[int]bool)
	for _, code := range expectedCodes {
		codeMap[code] = true
	}

	for _, code := range codes {
		if !codeMap[code] {
			t.Errorf("RetryableStatusCodes() contains unexpected code %d", code)
		}
	}
}

func TestAPIEndpoints(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{"Chat Completions Path", PathChatCompletions, "/chat/completions"},
		{"Embeddings Path", PathEmbeddings, "/embeddings"},
		{"Images Path", PathImages, "/images/generations"},
		{"Files Path", PathFiles, "/files"},
		{"Models Path", PathModels, "/models"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.got != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestLimitsAndConstraints(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"Max Request Body Size", MaxRequestBodySize, 100 * 1024 * 1024},
		{"Default Max Tokens", DefaultMaxTokens, 2048},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.got != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}

// TestConstantImmutability ensures constants cannot be changed (compile-time check)
func TestConstantImmutability(t *testing.T) {
	t.Parallel()

	// This test doesn't actually test runtime behavior,
	// but serves as documentation that these are constants
	// and will cause compile errors if someone tries to modify them.

	// If this compiles, constants are properly defined.
	_ = SDKVersion
	_ = ZaiBaseURL
	_ = DefaultTimeout
	_ = HeaderAuthorization
	_ = DefaultSourceChannel
	_ = StatusTooManyRequests
	_ = PathChatCompletions
	_ = MaxRequestBodySize
}

func TestTimeoutValues(t *testing.T) {
	t.Parallel()

	// Verify timeout values are reasonable
	if DefaultTimeout <= 0 {
		t.Error("DefaultTimeout should be positive")
	}

	if DefaultConnectTimeout <= 0 {
		t.Error("DefaultConnectTimeout should be positive")
	}

	if DefaultTimeout < DefaultConnectTimeout {
		t.Error("DefaultTimeout should be >= DefaultConnectTimeout")
	}

	if InitialRetryDelay <= 0 {
		t.Error("InitialRetryDelay should be positive")
	}

	if MaxRetryDelay <= InitialRetryDelay {
		t.Error("MaxRetryDelay should be > InitialRetryDelay")
	}
}

func TestConnectionPoolValues(t *testing.T) {
	t.Parallel()

	// Verify connection pool values are reasonable
	if DefaultMaxConnections <= 0 {
		t.Error("DefaultMaxConnections should be positive")
	}

	if DefaultMaxIdleConns <= 0 {
		t.Error("DefaultMaxIdleConns should be positive")
	}

	if DefaultMaxIdleConns > DefaultMaxConnections {
		t.Error("DefaultMaxIdleConns should be <= DefaultMaxConnections")
	}

	if DefaultMaxIdleConnsPerHost > DefaultMaxIdleConns {
		t.Error("DefaultMaxIdleConnsPerHost should be <= DefaultMaxIdleConns")
	}
}

// Benchmark for GetUserAgent
func BenchmarkGetUserAgent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetUserAgent()
	}
}

// Benchmark for RetryableStatusCodes
func BenchmarkRetryableStatusCodes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RetryableStatusCodes()
	}
}
