package transport

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/sofianhadi1983/zai-sdk-go/internal/constants"
	"github.com/sofianhadi1983/zai-sdk-go/internal/logger"
)

// RetryConfig holds configuration for retry behavior.
type RetryConfig struct {
	// MaxRetries is the maximum number of retry attempts.
	MaxRetries int

	// InitialBackoff is the initial backoff duration.
	InitialBackoff time.Duration

	// MaxBackoff is the maximum backoff duration.
	MaxBackoff time.Duration

	// BackoffMultiplier is the multiplier for exponential backoff.
	BackoffMultiplier float64

	// RetryableStatusCodes are HTTP status codes that should trigger a retry.
	RetryableStatusCodes []int

	// EnableJitter adds randomness to backoff to prevent thundering herd.
	EnableJitter bool
}

// DefaultRetryConfig returns the default retry configuration.
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:           constants.DefaultMaxRetries,
		InitialBackoff:       constants.InitialRetryDelay,
		MaxBackoff:           constants.MaxRetryDelay,
		BackoffMultiplier:    constants.RetryBackoffMultiplier,
		RetryableStatusCodes: constants.RetryableStatusCodes(),
		EnableJitter:         true,
	}
}

// RetryableHTTPClient wraps HTTPClient with retry logic.
type RetryableHTTPClient struct {
	client *HTTPClient
	config *RetryConfig
	logger *logger.Logger
}

// NewRetryableHTTPClient creates a new retryable HTTP client.
func NewRetryableHTTPClient(client *HTTPClient, config *RetryConfig) *RetryableHTTPClient {
	if config == nil {
		config = DefaultRetryConfig()
	}

	return &RetryableHTTPClient{
		client: client,
		config: config,
		logger: logger.Default(),
	}
}

// SetLogger sets a custom logger for the retryable HTTP client.
func (c *RetryableHTTPClient) SetLogger(l *logger.Logger) {
	c.logger = l
}

// DoWithRetry executes an HTTP request with retry logic.
// It will retry on retryable errors and status codes with exponential backoff.
func (c *RetryableHTTPClient) DoWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	var lastErr error
	var resp *http.Response

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		// Check if context is cancelled before attempting
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Clone the request for retry attempts (except the first one)
		var reqToSend *http.Request
		if attempt == 0 {
			reqToSend = req
		} else {
			var err error
			reqToSend, err = c.cloneRequest(ctx, req)
			if err != nil {
				return nil, fmt.Errorf("failed to clone request: %w", err)
			}
		}

		// Log retry attempt
		if attempt > 0 && c.logger != nil {
			c.logger.InfoContext(ctx, "Retrying HTTP request",
				slog.Int("attempt", attempt),
				slog.Int("max_retries", c.config.MaxRetries),
			)
		}

		// Execute the request
		resp, lastErr = c.client.Do(ctx, reqToSend)

		// Check if we should retry
		shouldRetry, retryAfter := c.shouldRetry(resp, lastErr, attempt)
		if !shouldRetry {
			// Success or non-retryable error
			return resp, lastErr
		}

		// Close response body if present (we'll retry)
		if resp != nil && resp.Body != nil {
			io.Copy(io.Discard, resp.Body) // Drain the body
			resp.Body.Close()
		}

		// Don't sleep after the last attempt
		if attempt < c.config.MaxRetries {
			backoff := c.calculateBackoff(attempt, retryAfter)

			if c.logger != nil {
				c.logger.DebugContext(ctx, "Backing off before retry",
					slog.Duration("backoff", backoff),
					slog.Int("attempt", attempt+1),
				)
			}

			// Sleep with context awareness
			select {
			case <-time.After(backoff):
				// Continue to next retry
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}

	// All retries exhausted
	if c.logger != nil {
		c.logger.WarnContext(ctx, "All retry attempts exhausted",
			slog.Int("max_retries", c.config.MaxRetries),
		)
	}

	return resp, lastErr
}

// shouldRetry determines if a request should be retried based on the response and error.
// It returns whether to retry and an optional retry-after duration from the response headers.
func (c *RetryableHTTPClient) shouldRetry(resp *http.Response, err error, attempt int) (bool, time.Duration) {
	// Don't retry if we've exhausted all attempts
	if attempt >= c.config.MaxRetries {
		return false, 0
	}

	// If there's an error, retry for network errors
	if err != nil {
		// Don't retry on context cancellation
		if err == context.Canceled || err == context.DeadlineExceeded {
			return false, 0
		}
		// Retry on network errors
		return true, 0
	}

	// No response means network error, should have been caught above
	if resp == nil {
		return true, 0
	}

	// Check if the status code is retryable
	if c.isRetryableStatusCode(resp.StatusCode) {
		// Check if the HTTP method is idempotent
		if !c.isIdempotentMethod(resp.Request.Method) {
			return false, 0
		}

		// Parse Retry-After header if present
		retryAfter := c.parseRetryAfter(resp)
		return true, retryAfter
	}

	// Not a retryable status code
	return false, 0
}

// isRetryableStatusCode checks if a status code is in the list of retryable codes.
func (c *RetryableHTTPClient) isRetryableStatusCode(statusCode int) bool {
	for _, code := range c.config.RetryableStatusCodes {
		if code == statusCode {
			return true
		}
	}
	return false
}

// isIdempotentMethod checks if an HTTP method is idempotent (safe to retry).
func (c *RetryableHTTPClient) isIdempotentMethod(method string) bool {
	// GET, HEAD, OPTIONS, TRACE are always idempotent
	// PUT and DELETE are also idempotent
	// POST is not idempotent (unless specifically marked)
	// PATCH is not idempotent
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace,
		http.MethodPut, http.MethodDelete:
		return true
	default:
		return false
	}
}

// parseRetryAfter parses the Retry-After header from the response.
// It returns the duration to wait, or 0 if the header is not present or invalid.
func (c *RetryableHTTPClient) parseRetryAfter(resp *http.Response) time.Duration {
	if resp == nil {
		return 0
	}

	retryAfterHeader := resp.Header.Get("Retry-After")
	if retryAfterHeader == "" {
		return 0
	}

	// Try to parse as seconds (integer)
	if seconds, err := strconv.Atoi(retryAfterHeader); err == nil {
		if seconds > 0 && seconds <= 300 { // Cap at 5 minutes
			return time.Duration(seconds) * time.Second
		}
	}

	// Try to parse as HTTP date
	if t, err := http.ParseTime(retryAfterHeader); err == nil {
		duration := time.Until(t)
		if duration > 0 && duration <= 5*time.Minute {
			return duration
		}
	}

	return 0
}

// calculateBackoff calculates the backoff duration for a retry attempt.
// It uses exponential backoff with optional jitter.
func (c *RetryableHTTPClient) calculateBackoff(attempt int, retryAfter time.Duration) time.Duration {
	// If server provided a Retry-After header and it's reasonable, use it
	if retryAfter > 0 && retryAfter <= c.config.MaxBackoff {
		return retryAfter
	}

	// Calculate exponential backoff
	backoff := float64(c.config.InitialBackoff) * math.Pow(c.config.BackoffMultiplier, float64(attempt))

	// Apply maximum backoff limit
	if backoff > float64(c.config.MaxBackoff) {
		backoff = float64(c.config.MaxBackoff)
	}

	// Add jitter to prevent thundering herd
	if c.config.EnableJitter {
		// Add random jitter between 0% and 25% of the backoff
		jitter := backoff * 0.25 * rand.Float64()
		backoff += jitter
	}

	return time.Duration(backoff)
}

// cloneRequest creates a copy of an HTTP request for retry purposes.
// This is necessary because http.Request.Body can only be read once.
func (c *RetryableHTTPClient) cloneRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
	// Create a new request with the same properties
	cloned := req.Clone(ctx)

	// The body is already handled by Clone for GetBody
	// If GetBody is nil and Body is not nil, we can't retry
	if req.Body != nil && req.GetBody == nil {
		return nil, fmt.Errorf("request body cannot be retried (no GetBody function)")
	}

	// If GetBody is available, use it to get a fresh body
	if cloned.GetBody != nil {
		body, err := cloned.GetBody()
		if err != nil {
			return nil, fmt.Errorf("failed to get request body: %w", err)
		}
		cloned.Body = body
	}

	return cloned, nil
}

// Close closes the underlying HTTP client.
func (c *RetryableHTTPClient) Close() {
	c.client.Close()
}

// GetClient returns the underlying HTTP client.
func (c *RetryableHTTPClient) GetClient() *HTTPClient {
	return c.client
}
