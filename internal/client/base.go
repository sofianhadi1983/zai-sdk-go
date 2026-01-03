// Package client provides the base API client implementation.
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/z-ai/zai-sdk-go/internal/auth"
	"github.com/z-ai/zai-sdk-go/internal/constants"
	"github.com/z-ai/zai-sdk-go/internal/logger"
	"github.com/z-ai/zai-sdk-go/internal/models"
	"github.com/z-ai/zai-sdk-go/internal/streaming"
	"github.com/z-ai/zai-sdk-go/internal/transport"
	"github.com/z-ai/zai-sdk-go/pkg/zai/errors"
)

// Config holds configuration for the API client.
type Config struct {
	// APIKey is the API key for authentication (format: "key.secret").
	APIKey string

	// BaseURL is the base URL for API requests.
	// If empty, uses the default Z.ai API URL.
	BaseURL string

	// Timeout is the request timeout.
	// If zero, uses the default timeout.
	Timeout time.Duration

	// MaxRetries is the maximum number of retry attempts.
	// If zero, uses the default max retries.
	MaxRetries int

	// DisableTokenCache disables JWT token caching.
	// When true, uses raw API key for authentication.
	DisableTokenCache bool

	// HTTPClient is a custom HTTP client.
	// If nil, creates a default client.
	HTTPClient *http.Client

	// Logger is a custom logger.
	// If nil, uses the default logger.
	Logger *logger.Logger
}

// BaseClient is the base client for making API requests.
type BaseClient struct {
	config         *Config
	httpClient     *transport.RetryableHTTPClient
	tokenGenerator *auth.TokenGenerator
	logger         *logger.Logger
}

// NewBaseClient creates a new base API client.
func NewBaseClient(config *Config) (*BaseClient, error) {
	if config == nil {
		config = &Config{}
	}

	// Validate API key
	if config.APIKey == "" {
		return nil, errors.NewConfigError("APIKey", "API key is required")
	}

	// Set defaults
	if config.BaseURL == "" {
		config.BaseURL = constants.ZaiBaseURL
	}

	if config.Timeout == 0 {
		config.Timeout = constants.DefaultTimeout
	}

	if config.MaxRetries == 0 {
		config.MaxRetries = constants.DefaultMaxRetries
	}

	// Create logger
	log := config.Logger
	if log == nil {
		log = logger.Default()
	}

	// Create HTTP client
	httpConfig := &transport.HTTPClientConfig{
		BaseURL:        config.BaseURL,
		Timeout:        config.Timeout,
		ConnectTimeout: constants.DefaultConnectTimeout,
	}

	if config.HTTPClient != nil {
		// Use custom HTTP client
		httpConfig.MaxIdleConns = constants.DefaultMaxIdleConns
		httpConfig.MaxIdleConnsPerHost = constants.DefaultMaxIdleConnsPerHost
		httpConfig.IdleConnTimeout = constants.DefaultIdleConnTimeout
	}

	httpClient := transport.NewHTTPClient(httpConfig)
	httpClient.SetLogger(log)

	// Create retryable client
	retryConfig := &transport.RetryConfig{
		MaxRetries:           config.MaxRetries,
		InitialBackoff:       constants.InitialRetryDelay,
		MaxBackoff:           constants.MaxRetryDelay,
		BackoffMultiplier:    constants.RetryBackoffMultiplier,
		RetryableStatusCodes: constants.RetryableStatusCodes(),
		EnableJitter:         true,
	}

	retryableClient := transport.NewRetryableHTTPClient(httpClient, retryConfig)
	retryableClient.SetLogger(log)

	// Create token generator
	tokenGen := auth.NewTokenGenerator()
	if config.DisableTokenCache {
		tokenGen.DisableCache()
	}

	return &BaseClient{
		config:         config,
		httpClient:     retryableClient,
		tokenGenerator: tokenGen,
		logger:         log,
	}, nil
}

// Do executes an HTTP request with retry and authentication.
func (c *BaseClient) Do(ctx context.Context, req *http.Request) (*models.APIResponse, error) {
	// Add authentication
	if err := c.addAuth(req); err != nil {
		return nil, err
	}

	// Execute with retry
	start := time.Now()
	resp, err := c.httpClient.DoWithRetry(ctx, req)
	elapsed := time.Since(start)

	if err != nil {
		return nil, err
	}

	// Wrap response
	apiResp := models.NewAPIResponse(resp, elapsed)

	// Check for errors
	if apiResp.IsError() {
		return apiResp, c.handleErrorResponse(apiResp)
	}

	return apiResp, nil
}

// Get performs a GET request.
func (c *BaseClient) Get(ctx context.Context, path string, query map[string]string) (*models.APIResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// Add query parameters
	if len(query) > 0 {
		q := req.URL.Query()
		for k, v := range query {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	return c.Do(ctx, req)
}

// Post performs a POST request with JSON body.
func (c *BaseClient) Post(ctx context.Context, path string, body interface{}) (*models.APIResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}

	return c.Do(ctx, req)
}

// PostMultipart performs a POST request with multipart/form-data body.
func (c *BaseClient) PostMultipart(ctx context.Context, path string, body io.Reader, contentType string) (*models.APIResponse, error) {
	req, err := c.httpClient.GetClient().NewRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}

	// Set content type for multipart form data
	req.Header.Set("Content-Type", contentType)

	return c.Do(ctx, req)
}

// Put performs a PUT request with JSON body.
func (c *BaseClient) Put(ctx context.Context, path string, body interface{}) (*models.APIResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}

	return c.Do(ctx, req)
}

// Delete performs a DELETE request.
func (c *BaseClient) Delete(ctx context.Context, path string) (*models.APIResponse, error) {
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(ctx, req)
}

// Stream performs a streaming request.
func (c *BaseClient) Stream(ctx context.Context, path string, body interface{}) (*models.StreamResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}

	// Add authentication
	if err := c.addAuth(req); err != nil {
		return nil, err
	}

	// Execute request (no retry for streaming)
	start := time.Now()
	resp, err := c.httpClient.GetClient().Do(ctx, req)
	elapsed := time.Since(start)

	if err != nil {
		return nil, err
	}

	// Wrap response
	apiResp := models.NewAPIResponse(resp, elapsed)

	// Check for errors
	if apiResp.IsError() {
		return nil, c.handleErrorResponse(apiResp)
	}

	return models.NewStreamResponse(apiResp), nil
}

// ParseJSON parses a JSON response into the given type.
func (c *BaseClient) ParseJSON(resp *models.APIResponse, v interface{}) error {
	defer resp.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// NewTypedStream creates a new typed stream from a stream response.
// This is a package-level function due to Go's limitation on generic methods.
func NewTypedStream[T any](streamResp *models.StreamResponse, ctx context.Context) *streaming.Stream[T] {
	if ctx == nil {
		ctx = context.Background()
	}

	return streaming.NewStream[T](streaming.StreamConfig[T]{
		Reader:  streamResp.Body,
		Context: ctx,
	})
}

// Close closes the client and releases resources.
func (c *BaseClient) Close() {
	c.httpClient.Close()
}

// newRequest creates a new HTTP request.
func (c *BaseClient) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	var bodyReader io.Reader

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		// bytes package is imported at the top
		bodyReader = newBytesReader(data)
	}

	return c.httpClient.GetClient().NewRequest(ctx, method, path, bodyReader)
}

// newBytesReader creates a bytes.Reader from data.
func newBytesReader(data []byte) io.Reader {
	// Import bytes in the package imports
	r := &bytesReader{data: data}
	return r
}

type bytesReader struct {
	data []byte
	pos  int
}

func (r *bytesReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

// addAuth adds authentication to the request.
func (c *BaseClient) addAuth(req *http.Request) error {
	var token string

	if c.config.DisableTokenCache {
		// Use raw API key
		token = c.config.APIKey
	} else {
		// Generate JWT token
		var err error
		token, err = c.tokenGenerator.GenerateToken(c.config.APIKey)
		if err != nil {
			return fmt.Errorf("failed to generate auth token: %w", err)
		}
	}

	req.Header.Set(constants.HeaderAuthorization, "Bearer "+token)
	return nil
}

// handleErrorResponse converts an error response to an error.
func (c *BaseClient) handleErrorResponse(resp *models.APIResponse) error {
	defer resp.Close()

	// Try to parse error response
	var errResp models.ErrorResponse
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		// Fallback to generic error
		return errors.NewAPIStatusError(
			fmt.Sprintf("HTTP %d: failed to read error response", resp.StatusCode),
			resp.StatusCode,
			resp.HTTPResponse,
		)
	}

	if err := json.Unmarshal(data, &errResp); err != nil {
		// Fallback to generic error with body
		return errors.NewAPIStatusError(
			fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(data)),
			resp.StatusCode,
			resp.HTTPResponse,
		)
	}

	// Create specific error based on status code
	message := errResp.GetMessage()
	statusCode := resp.StatusCode

	switch statusCode {
	case http.StatusBadRequest:
		return errors.NewAPIRequestFailedError(message, statusCode, resp.HTTPResponse)

	case http.StatusUnauthorized:
		return errors.NewAPIAuthenticationError(message, statusCode, resp.HTTPResponse)

	case http.StatusTooManyRequests:
		return errors.NewAPIReachLimitError(message, statusCode, resp.HTTPResponse)

	case http.StatusInternalServerError:
		return errors.NewAPIInternalError(message, statusCode, resp.HTTPResponse)

	case http.StatusServiceUnavailable:
		return errors.NewAPIServerFlowExceedError(message, statusCode, resp.HTTPResponse)

	default:
		return errors.NewAPIStatusError(message, statusCode, resp.HTTPResponse)
	}
}

// GetConfig returns the client configuration.
func (c *BaseClient) GetConfig() *Config {
	return c.config
}

// GetLogger returns the client logger.
func (c *BaseClient) GetLogger() *logger.Logger {
	return c.logger
}
