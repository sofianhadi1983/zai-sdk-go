package models

import (
	"io"
	"time"
)

// RequestOptions holds configuration for making HTTP requests.
type RequestOptions struct {
	// Method is the HTTP method (GET, POST, etc.).
	Method string

	// URL is the endpoint URL.
	URL string

	// Query holds query parameters.
	Query map[string]string

	// Headers holds custom headers.
	Headers map[string]string

	// Body is the request body reader.
	Body io.Reader

	// JSONData is the request body as a JSON-serializable object.
	// If set, this takes precedence over Body.
	JSONData interface{}

	// Timeout is the request timeout.
	Timeout time.Duration

	// MaxRetries is the maximum number of retry attempts.
	MaxRetries int

	// IdempotencyKey is used for idempotent requests.
	IdempotencyKey string

	// Stream indicates if this is a streaming request.
	Stream bool
}

// NewRequestOptions creates a new RequestOptions with defaults.
func NewRequestOptions() *RequestOptions {
	return &RequestOptions{
		Query:   make(map[string]string),
		Headers: make(map[string]string),
	}
}

// WithMethod sets the HTTP method.
func (r *RequestOptions) WithMethod(method string) *RequestOptions {
	r.Method = method
	return r
}

// WithURL sets the URL.
func (r *RequestOptions) WithURL(url string) *RequestOptions {
	r.URL = url
	return r
}

// WithQuery adds a query parameter.
func (r *RequestOptions) WithQuery(key, value string) *RequestOptions {
	if r.Query == nil {
		r.Query = make(map[string]string)
	}
	r.Query[key] = value
	return r
}

// WithHeader adds a header.
func (r *RequestOptions) WithHeader(key, value string) *RequestOptions {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers[key] = value
	return r
}

// WithBody sets the request body.
func (r *RequestOptions) WithBody(body io.Reader) *RequestOptions {
	r.Body = body
	return r
}

// WithJSONData sets the JSON data.
func (r *RequestOptions) WithJSONData(data interface{}) *RequestOptions {
	r.JSONData = data
	return r
}

// WithTimeout sets the timeout.
func (r *RequestOptions) WithTimeout(timeout time.Duration) *RequestOptions {
	r.Timeout = timeout
	return r
}

// WithMaxRetries sets the max retries.
func (r *RequestOptions) WithMaxRetries(maxRetries int) *RequestOptions {
	r.MaxRetries = maxRetries
	return r
}

// WithIdempotencyKey sets the idempotency key.
func (r *RequestOptions) WithIdempotencyKey(key string) *RequestOptions {
	r.IdempotencyKey = key
	return r
}

// WithStream sets the stream flag.
func (r *RequestOptions) WithStream(stream bool) *RequestOptions {
	r.Stream = stream
	return r
}
