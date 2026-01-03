package models

import (
	"io"
	"net/http"
	"time"
)

// APIResponse wraps an HTTP response with additional metadata.
type APIResponse struct {
	// HTTPResponse is the underlying HTTP response.
	HTTPResponse *http.Response

	// Body is the response body reader.
	Body io.ReadCloser

	// Headers are the response headers.
	Headers http.Header

	// StatusCode is the HTTP status code.
	StatusCode int

	// URL is the request URL.
	URL string

	// Method is the HTTP method used.
	Method string

	// HTTPVersion is the HTTP protocol version.
	HTTPVersion string

	// Elapsed is the request duration.
	Elapsed time.Duration

	// RequestID is extracted from the response headers.
	RequestID string

	// IsClosed indicates if the response body has been closed.
	IsClosed bool
}

// NewAPIResponse creates a new APIResponse from an http.Response.
func NewAPIResponse(resp *http.Response, elapsed time.Duration) *APIResponse {
	requestID := resp.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = resp.Header.Get("Request-ID")
	}

	return &APIResponse{
		HTTPResponse: resp,
		Body:         resp.Body,
		Headers:      resp.Header,
		StatusCode:   resp.StatusCode,
		URL:          resp.Request.URL.String(),
		Method:       resp.Request.Method,
		HTTPVersion:  resp.Proto,
		Elapsed:      elapsed,
		RequestID:    requestID,
		IsClosed:     false,
	}
}

// Close closes the response body.
func (r *APIResponse) Close() error {
	if r.IsClosed || r.Body == nil {
		return nil
	}

	r.IsClosed = true
	return r.Body.Close()
}

// GetHeader retrieves a header value.
func (r *APIResponse) GetHeader(key string) string {
	return r.Headers.Get(key)
}

// IsSuccess returns true if the status code indicates success (2xx).
func (r *APIResponse) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsClientError returns true if the status code indicates a client error (4xx).
func (r *APIResponse) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

// IsServerError returns true if the status code indicates a server error (5xx).
func (r *APIResponse) IsServerError() bool {
	return r.StatusCode >= 500 && r.StatusCode < 600
}

// IsError returns true if the status code indicates an error (4xx or 5xx).
func (r *APIResponse) IsError() bool {
	return r.IsClientError() || r.IsServerError()
}

// StreamResponse represents a streaming API response.
type StreamResponse struct {
	*APIResponse

	// Reader is the stream reader.
	Reader io.Reader

	// Done is a channel that closes when the stream is complete.
	Done chan struct{}

	// Err holds any error that occurred during streaming.
	Err error
}

// NewStreamResponse creates a new StreamResponse.
func NewStreamResponse(apiResp *APIResponse) *StreamResponse {
	return &StreamResponse{
		APIResponse: apiResp,
		Reader:      apiResp.Body,
		Done:        make(chan struct{}),
	}
}

// Close closes the stream and the underlying response.
func (s *StreamResponse) Close() error {
	close(s.Done)
	return s.APIResponse.Close()
}
