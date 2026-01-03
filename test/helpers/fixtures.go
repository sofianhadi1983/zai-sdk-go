package helpers

import "time"

// Test constants
const (
	// MockAPIKey is a fake API key for testing.
	MockAPIKey = "test-api-key-12345"

	// MockInvalidAPIKey is an invalid API key for testing error cases.
	MockInvalidAPIKey = "invalid-key"

	// TestBaseURL is the base URL for test API endpoints.
	TestBaseURL = "https://api.test.z.ai/api/paas/v4"

	// TestZhipuBaseURL is the Zhipu base URL for testing.
	TestZhipuBaseURL = "https://test.bigmodel.cn/api/paas/v4"

	// MockRequestID is a sample request ID for testing.
	MockRequestID = "req-test-12345"

	// MockTraceID is a sample trace ID for testing.
	MockTraceID = "trace-test-67890"

	// MockUserID is a sample user ID for testing.
	MockUserID = "user-test-abc123"
)

// Test timeouts
const (
	// ShortTimeout is a short timeout for quick tests (1 second).
	ShortTimeout = 1 * time.Second

	// MediumTimeout is a medium timeout for most tests (5 seconds).
	MediumTimeout = 5 * time.Second

	// LongTimeout is a long timeout for slow tests (30 seconds).
	LongTimeout = 30 * time.Second
)

// MockChatCompletionRequest returns a sample chat completion request.
func MockChatCompletionRequest() map[string]interface{} {
	return map[string]interface{}{
		"model": "glm-4",
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": "Hello, how are you?",
			},
		},
		"temperature": 0.7,
		"max_tokens":  100,
	}
}

// MockChatCompletionResponse returns a sample chat completion response.
func MockChatCompletionResponse() string {
	return `{
		"id": "chatcmpl-123",
		"object": "chat.completion",
		"created": 1677652288,
		"model": "glm-4",
		"choices": [
			{
				"index": 0,
				"message": {
					"role": "assistant",
					"content": "I'm doing well, thank you for asking!"
				},
				"finish_reason": "stop"
			}
		],
		"usage": {
			"prompt_tokens": 10,
			"completion_tokens": 12,
			"total_tokens": 22
		}
	}`
}

// MockStreamingChatCompletionResponse returns a sample streaming response chunk.
func MockStreamingChatCompletionResponse() string {
	return `data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"glm-4","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null}]}

`
}

// MockEmbeddingRequest returns a sample embedding request.
func MockEmbeddingRequest() map[string]interface{} {
	return map[string]interface{}{
		"model": "embedding-2",
		"input": "The quick brown fox jumps over the lazy dog.",
	}
}

// MockEmbeddingResponse returns a sample embedding response.
func MockEmbeddingResponse() string {
	return `{
		"object": "list",
		"data": [
			{
				"object": "embedding",
				"embedding": [0.1, 0.2, 0.3, 0.4, 0.5],
				"index": 0
			}
		],
		"model": "embedding-2",
		"usage": {
			"prompt_tokens": 8,
			"total_tokens": 8
		}
	}`
}

// MockImageGenerationRequest returns a sample image generation request.
func MockImageGenerationRequest() map[string]interface{} {
	return map[string]interface{}{
		"model":  "cogview-3",
		"prompt": "A beautiful sunset over the ocean",
		"n":      1,
		"size":   "1024x1024",
	}
}

// MockImageGenerationResponse returns a sample image generation response.
func MockImageGenerationResponse() string {
	return `{
		"created": 1677652288,
		"data": [
			{
				"url": "https://example.com/image.png"
			}
		]
	}`
}

// MockErrorResponse returns a sample error response.
func MockErrorResponse(statusCode int, message string) string {
	return `{
		"error": {
			"message": "` + message + `",
			"type": "invalid_request_error",
			"code": "` + string(rune(statusCode)) + `"
		}
	}`
}

// MockAuthErrorResponse returns a sample authentication error response.
func MockAuthErrorResponse() string {
	return `{
		"error": {
			"message": "Invalid API key",
			"type": "authentication_error",
			"code": "401"
		}
	}`
}

// MockRateLimitErrorResponse returns a sample rate limit error response.
func MockRateLimitErrorResponse() string {
	return `{
		"error": {
			"message": "Rate limit exceeded",
			"type": "rate_limit_error",
			"code": "429"
		}
	}`
}

// MockInternalErrorResponse returns a sample internal server error response.
func MockInternalErrorResponse() string {
	return `{
		"error": {
			"message": "Internal server error",
			"type": "server_error",
			"code": "500"
		}
	}`
}

// MockValidationErrorResponse returns a sample validation error response.
func MockValidationErrorResponse(field, message string) string {
	return `{
		"error": {
			"message": "Validation failed for field '` + field + `': ` + message + `",
			"type": "validation_error",
			"code": "400"
		}
	}`
}

// MockJWTToken returns a sample JWT token for testing.
func MockJWTToken() string {
	// This is a mock JWT token - not a real one, just for testing structure
	return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlRlc3QgVXNlciIsImlhdCI6MTUxNjIzOTAyMn0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}

// MockFileUploadResponse returns a sample file upload response.
func MockFileUploadResponse() string {
	return `{
		"id": "file-abc123",
		"object": "file",
		"bytes": 1024,
		"created_at": 1677652288,
		"filename": "test.txt",
		"purpose": "fine-tune"
	}`
}

// MockModelListResponse returns a sample model list response.
func MockModelListResponse() string {
	return `{
		"object": "list",
		"data": [
			{
				"id": "glm-4",
				"object": "model",
				"created": 1677610602,
				"owned_by": "zhipu-ai"
			},
			{
				"id": "glm-3-turbo",
				"object": "model",
				"created": 1677610602,
				"owned_by": "zhipu-ai"
			}
		]
	}`
}

// MockHeaders returns a sample set of HTTP headers for testing.
func MockHeaders() map[string]string {
	return map[string]string{
		"Content-Type":     "application/json",
		"Authorization":    "Bearer " + MockAPIKey,
		"User-Agent":       "Z.ai Go SDK/0.1.0",
		"x-source-channel": "go-sdk",
		"x-request-id":     MockRequestID,
	}
}
