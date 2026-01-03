// Package constants defines SDK-wide constants and default values.
package constants

import "time"

const (
	// SDKVersion is the current version of the Z.ai Go SDK.
	SDKVersion = "0.1.0"

	// SDKTitle is the name of the SDK.
	SDKTitle = "Z.ai Go SDK"
)

// API Base URLs

const (
	// ZaiBaseURL is the default base URL for Z.ai API (international).
	// Used by ZaiClient for users outside mainland China.
	ZaiBaseURL = "https://api.z.ai/api/paas/v4"

	// ZhipuBaseURL is the default base URL for Zhipu AI API (mainland China).
	// Used by ZhipuAiClient for users in mainland China.
	ZhipuBaseURL = "https://open.bigmodel.cn/api/paas/v4"
)

// HTTP Client Defaults

const (
	// DefaultTimeout is the default timeout for HTTP requests.
	// Equivalent to Python SDK's timeout=300.0 (total timeout).
	DefaultTimeout = 300 * time.Second

	// DefaultConnectTimeout is the default connection timeout.
	// Equivalent to Python SDK's connect=8.0.
	DefaultConnectTimeout = 8 * time.Second

	// DefaultMaxRetries is the default maximum number of retry attempts.
	DefaultMaxRetries = 3

	// DefaultMaxConnections is the default maximum number of connections.
	// Equivalent to Python SDK's max_connections=50.
	DefaultMaxConnections = 50

	// DefaultMaxIdleConns is the default maximum number of idle connections.
	// Equivalent to Python SDK's max_keepalive_connections=10.
	DefaultMaxIdleConns = 10

	// DefaultMaxIdleConnsPerHost is the maximum idle connections per host.
	DefaultMaxIdleConnsPerHost = 10

	// DefaultIdleConnTimeout is the timeout for idle connections.
	DefaultIdleConnTimeout = 90 * time.Second
)

// Retry Configuration

const (
	// InitialRetryDelay is the initial delay before first retry.
	// Equivalent to Python SDK's INITIAL_RETRY_DELAY = 0.5.
	InitialRetryDelay = 500 * time.Millisecond

	// MaxRetryDelay is the maximum delay between retries.
	// Equivalent to Python SDK's MAX_RETRY_DELAY = 8.0.
	MaxRetryDelay = 8 * time.Second

	// RetryBackoffMultiplier is the exponential backoff multiplier.
	RetryBackoffMultiplier = 2.0
)

// HTTP Headers

const (
	// HeaderAuthorization is the HTTP Authorization header name.
	HeaderAuthorization = "Authorization"

	// HeaderContentType is the HTTP Content-Type header name.
	HeaderContentType = "Content-Type"

	// HeaderAccept is the HTTP Accept header name.
	HeaderAccept = "Accept"

	// HeaderAcceptLanguage is the HTTP Accept-Language header name.
	HeaderAcceptLanguage = "Accept-Language"

	// HeaderUserAgent is the HTTP User-Agent header name.
	HeaderUserAgent = "User-Agent"

	// HeaderSourceChannel is the custom header for source channel tracking.
	// Equivalent to Python SDK's x-source-channel.
	HeaderSourceChannel = "x-source-channel"

	// HeaderRawResponse is used for raw response handling.
	// Equivalent to Python SDK's X-Stainless-Raw-Response.
	HeaderRawResponse = "X-Stainless-Raw-Response"
)

// Header Values

const (
	// DefaultSourceChannel is the default value for source channel header.
	// Changed from "python-sdk" to "go-sdk" for Go SDK.
	DefaultSourceChannel = "go-sdk"

	// ContentTypeJSON is the JSON content type.
	ContentTypeJSON = "application/json"

	// ContentTypeFormURLEncoded is the form URL encoded content type.
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"

	// ContentTypeMultipartFormData is the multipart form data content type prefix.
	// Actual value will include boundary parameter.
	ContentTypeMultipartFormData = "multipart/form-data"

	// AcceptJSON is the JSON accept type.
	AcceptJSON = "application/json"

	// AcceptLanguageEnglish is the English accept language.
	AcceptLanguageEnglish = "en-US,en"

	// AcceptLanguageChinese is the Chinese accept language.
	AcceptLanguageChinese = "zh-CN,zh"
)

// User Agent

// GetUserAgent returns the user agent string for HTTP requests.
func GetUserAgent() string {
	return SDKTitle + "/" + SDKVersion
}

// Status Codes for Retry Logic

const (
	// StatusTooManyRequests indicates rate limiting (429).
	StatusTooManyRequests = 429

	// StatusInternalServerError indicates server error (500).
	StatusInternalServerError = 500

	// StatusBadGateway indicates bad gateway (502).
	StatusBadGateway = 502

	// StatusServiceUnavailable indicates service unavailable (503).
	StatusServiceUnavailable = 503

	// StatusGatewayTimeout indicates gateway timeout (504).
	StatusGatewayTimeout = 504
)

// RetryableStatusCodes returns the list of HTTP status codes that should trigger a retry.
func RetryableStatusCodes() []int {
	return []int{
		StatusTooManyRequests,
		StatusInternalServerError,
		StatusBadGateway,
		StatusServiceUnavailable,
		StatusGatewayTimeout,
	}
}

// API Endpoints (common paths)

const (
	// PathChatCompletions is the path for chat completions API.
	PathChatCompletions = "/chat/completions"

	// PathEmbeddings is the path for embeddings API.
	PathEmbeddings = "/embeddings"

	// PathImages is the path for image generation API.
	PathImages = "/images/generations"

	// PathFiles is the path for file management API.
	PathFiles = "/files"

	// PathModels is the path for models API.
	PathModels = "/models"
)

// Limits and Constraints

const (
	// MaxRequestBodySize is the maximum size of a request body (100MB).
	MaxRequestBodySize = 100 * 1024 * 1024

	// DefaultMaxTokens is a reasonable default for max tokens if not specified.
	DefaultMaxTokens = 2048
)
