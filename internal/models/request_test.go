package models

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRequestOptions(t *testing.T) {
	t.Parallel()

	opts := NewRequestOptions()

	assert.NotNil(t, opts)
	assert.NotNil(t, opts.Query)
	assert.NotNil(t, opts.Headers)
	assert.Empty(t, opts.Query)
	assert.Empty(t, opts.Headers)
}

func TestRequestOptions_WithMethod(t *testing.T) {
	t.Parallel()

	opts := NewRequestOptions().WithMethod("POST")

	assert.Equal(t, "POST", opts.Method)
}

func TestRequestOptions_WithURL(t *testing.T) {
	t.Parallel()

	opts := NewRequestOptions().WithURL("https://api.example.com/v1/chat")

	assert.Equal(t, "https://api.example.com/v1/chat", opts.URL)
}

func TestRequestOptions_WithQuery(t *testing.T) {
	t.Parallel()

	opts := NewRequestOptions().
		WithQuery("param1", "value1").
		WithQuery("param2", "value2")

	assert.Len(t, opts.Query, 2)
	assert.Equal(t, "value1", opts.Query["param1"])
	assert.Equal(t, "value2", opts.Query["param2"])
}

func TestRequestOptions_WithHeader(t *testing.T) {
	t.Parallel()

	opts := NewRequestOptions().
		WithHeader("Content-Type", "application/json").
		WithHeader("Authorization", "Bearer token")

	assert.Len(t, opts.Headers, 2)
	assert.Equal(t, "application/json", opts.Headers["Content-Type"])
	assert.Equal(t, "Bearer token", opts.Headers["Authorization"])
}

func TestRequestOptions_WithBody(t *testing.T) {
	t.Parallel()

	body := strings.NewReader("test body")
	opts := NewRequestOptions().WithBody(body)

	assert.NotNil(t, opts.Body)
	assert.Equal(t, body, opts.Body)
}

func TestRequestOptions_WithJSONData(t *testing.T) {
	t.Parallel()

	data := map[string]interface{}{
		"key": "value",
	}

	opts := NewRequestOptions().WithJSONData(data)

	assert.NotNil(t, opts.JSONData)
	assert.Equal(t, data, opts.JSONData)
}

func TestRequestOptions_WithTimeout(t *testing.T) {
	t.Parallel()

	timeout := 30 * time.Second
	opts := NewRequestOptions().WithTimeout(timeout)

	assert.Equal(t, timeout, opts.Timeout)
}

func TestRequestOptions_WithMaxRetries(t *testing.T) {
	t.Parallel()

	opts := NewRequestOptions().WithMaxRetries(5)

	assert.Equal(t, 5, opts.MaxRetries)
}

func TestRequestOptions_WithIdempotencyKey(t *testing.T) {
	t.Parallel()

	opts := NewRequestOptions().WithIdempotencyKey("key-123")

	assert.Equal(t, "key-123", opts.IdempotencyKey)
}

func TestRequestOptions_WithStream(t *testing.T) {
	t.Parallel()

	opts := NewRequestOptions().WithStream(true)

	assert.True(t, opts.Stream)
}

func TestRequestOptions_Chaining(t *testing.T) {
	t.Parallel()

	opts := NewRequestOptions().
		WithMethod("POST").
		WithURL("https://api.example.com/v1/chat").
		WithQuery("model", "glm-4").
		WithHeader("Content-Type", "application/json").
		WithTimeout(30 * time.Second).
		WithMaxRetries(3).
		WithStream(true)

	assert.Equal(t, "POST", opts.Method)
	assert.Equal(t, "https://api.example.com/v1/chat", opts.URL)
	assert.Equal(t, "glm-4", opts.Query["model"])
	assert.Equal(t, "application/json", opts.Headers["Content-Type"])
	assert.Equal(t, 30*time.Second, opts.Timeout)
	assert.Equal(t, 3, opts.MaxRetries)
	assert.True(t, opts.Stream)
}
