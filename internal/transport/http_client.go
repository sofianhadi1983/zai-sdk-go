// Package transport provides HTTP client functionality for the Z.ai SDK.
package transport

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/z-ai/zai-sdk-go/internal/constants"
	"github.com/z-ai/zai-sdk-go/internal/logger"
)

// HTTPClientConfig holds configuration for the HTTP client.
type HTTPClientConfig struct {
	// BaseURL is the base URL for all requests.
	BaseURL string

	// Timeout is the maximum duration for a request.
	Timeout time.Duration

	// ConnectTimeout is the maximum duration for establishing a connection.
	ConnectTimeout time.Duration

	// MaxIdleConns is the maximum number of idle connections across all hosts.
	MaxIdleConns int

	// MaxIdleConnsPerHost is the maximum number of idle connections per host.
	MaxIdleConnsPerHost int

	// IdleConnTimeout is the maximum duration an idle connection is kept alive.
	IdleConnTimeout time.Duration

	// TLSHandshakeTimeout is the maximum duration for the TLS handshake.
	TLSHandshakeTimeout time.Duration

	// ExpectContinueTimeout is the maximum duration to wait for a server's first
	// response headers after fully writing the request headers if the request has
	// an "Expect: 100-continue" header.
	ExpectContinueTimeout time.Duration

	// TLSConfig allows customizing TLS configuration.
	TLSConfig *tls.Config

	// DisableKeepAlives disables HTTP keep-alives.
	DisableKeepAlives bool

	// DisableCompression disables compression.
	DisableCompression bool

	// MaxConnsPerHost limits the total number of connections per host.
	MaxConnsPerHost int
}

// DefaultHTTPClientConfig returns the default HTTP client configuration.
func DefaultHTTPClientConfig() *HTTPClientConfig {
	return &HTTPClientConfig{
		BaseURL:               constants.ZaiBaseURL,
		Timeout:               constants.DefaultTimeout,
		ConnectTimeout:        constants.DefaultConnectTimeout,
		MaxIdleConns:          constants.DefaultMaxIdleConns,
		MaxIdleConnsPerHost:   constants.DefaultMaxIdleConnsPerHost,
		IdleConnTimeout:       constants.DefaultIdleConnTimeout,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     false,
		DisableCompression:    false,
		MaxConnsPerHost:       0, // 0 means no limit
	}
}

// RequestMiddleware is a function that can modify a request before it's sent.
type RequestMiddleware func(*http.Request) error

// ResponseMiddleware is a function that can process a response after it's received.
type ResponseMiddleware func(*http.Response) error

// HTTPClient is a wrapper around http.Client with additional functionality.
type HTTPClient struct {
	client              *http.Client
	config              *HTTPClientConfig
	requestMiddlewares  []RequestMiddleware
	responseMiddlewares []ResponseMiddleware
	logger              *logger.Logger
}

// NewHTTPClient creates a new HTTP client with the given configuration.
func NewHTTPClient(config *HTTPClientConfig) *HTTPClient {
	if config == nil {
		config = DefaultHTTPClientConfig()
	}

	// Create custom transport
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   config.ConnectTimeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          config.MaxIdleConns,
		MaxIdleConnsPerHost:   config.MaxIdleConnsPerHost,
		MaxConnsPerHost:       config.MaxConnsPerHost,
		IdleConnTimeout:       config.IdleConnTimeout,
		TLSHandshakeTimeout:   config.TLSHandshakeTimeout,
		ExpectContinueTimeout: config.ExpectContinueTimeout,
		DisableKeepAlives:     config.DisableKeepAlives,
		DisableCompression:    config.DisableCompression,
	}

	// Apply custom TLS config if provided
	if config.TLSConfig != nil {
		transport.TLSClientConfig = config.TLSConfig
	}

	return &HTTPClient{
		client: &http.Client{
			Transport: transport,
			Timeout:   config.Timeout,
		},
		config:              config,
		requestMiddlewares:  make([]RequestMiddleware, 0),
		responseMiddlewares: make([]ResponseMiddleware, 0),
		logger:              logger.Default(),
	}
}

// SetLogger sets a custom logger for the HTTP client.
func (c *HTTPClient) SetLogger(l *logger.Logger) {
	c.logger = l
}

// AddRequestMiddleware adds a middleware to process requests before they're sent.
func (c *HTTPClient) AddRequestMiddleware(middleware RequestMiddleware) {
	c.requestMiddlewares = append(c.requestMiddlewares, middleware)
}

// AddResponseMiddleware adds a middleware to process responses after they're received.
func (c *HTTPClient) AddResponseMiddleware(middleware ResponseMiddleware) {
	c.responseMiddlewares = append(c.responseMiddlewares, middleware)
}

// NewRequest creates a new HTTP request with the given parameters.
// The URL can be absolute or relative to the base URL.
func (c *HTTPClient) NewRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	// Build the full URL
	fullURL, err := c.buildURL(path)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	c.setDefaultHeaders(req)

	return req, nil
}

// Do executes an HTTP request and returns the response.
// The response body must be closed by the caller.
func (c *HTTPClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Ensure the request has a context
	if req.Context() == nil {
		req = req.WithContext(ctx)
	}

	// Apply request middlewares
	for _, middleware := range c.requestMiddlewares {
		if err := middleware(req); err != nil {
			return nil, fmt.Errorf("request middleware error: %w", err)
		}
	}

	// Log the request
	if c.logger != nil {
		c.logger.DebugContext(ctx, "HTTP request",
			slog.String("method", req.Method),
			slog.String("url", req.URL.String()),
		)
	}

	// Execute the request
	resp, err := c.client.Do(req)
	if err != nil {
		if c.logger != nil {
			c.logger.ErrorContext(ctx, "HTTP request failed",
				slog.String("error", err.Error()),
			)
		}
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Apply response middlewares
	for _, middleware := range c.responseMiddlewares {
		if err := middleware(resp); err != nil {
			resp.Body.Close() // Clean up on middleware error
			return nil, fmt.Errorf("response middleware error: %w", err)
		}
	}

	// Log the response
	if c.logger != nil {
		c.logger.DebugContext(ctx, "HTTP response",
			slog.Int("status_code", resp.StatusCode),
			slog.String("status", resp.Status),
		)
	}

	return resp, nil
}

// Close closes idle connections in the HTTP client.
func (c *HTTPClient) Close() {
	c.client.CloseIdleConnections()
}

// buildURL constructs the full URL from the base URL and path.
func (c *HTTPClient) buildURL(path string) (string, error) {
	// If path is already an absolute URL, use it as-is
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path, nil
	}

	// Parse the base URL
	baseURL, err := url.Parse(c.config.BaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	// Ensure base URL ends with /
	basePath := baseURL.Path
	if !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}

	// Remove leading / from path if present
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	// Combine base path and request path
	fullPath := basePath + path

	// Build the final URL
	fullURL := &url.URL{
		Scheme: baseURL.Scheme,
		Host:   baseURL.Host,
		Path:   fullPath,
	}

	return fullURL.String(), nil
}

// setDefaultHeaders sets default headers on the request.
func (c *HTTPClient) setDefaultHeaders(req *http.Request) {
	// Set default headers if not already set
	if req.Header.Get(constants.HeaderContentType) == "" {
		req.Header.Set(constants.HeaderContentType, constants.ContentTypeJSON)
	}
	if req.Header.Get(constants.HeaderAccept) == "" {
		req.Header.Set(constants.HeaderAccept, constants.AcceptJSON)
	}
	if req.Header.Get(constants.HeaderUserAgent) == "" {
		req.Header.Set(constants.HeaderUserAgent, constants.GetUserAgent())
	}
	if req.Header.Get(constants.HeaderSourceChannel) == "" {
		req.Header.Set(constants.HeaderSourceChannel, constants.DefaultSourceChannel)
	}
}
