package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorResponse_GetMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		err      *ErrorResponse
		expected string
	}{
		{
			name: "message at top level",
			err: &ErrorResponse{
				Message: "top level error",
			},
			expected: "top level error",
		},
		{
			name: "message in nested error",
			err: &ErrorResponse{
				Error: &ErrorDetail{
					Message: "nested error",
				},
			},
			expected: "nested error",
		},
		{
			name: "both messages (top level takes precedence)",
			err: &ErrorResponse{
				Message: "top level",
				Error: &ErrorDetail{
					Message: "nested",
				},
			},
			expected: "top level",
		},
		{
			name:     "no message",
			err:      &ErrorResponse{},
			expected: "Unknown error",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, tt.err.GetMessage())
		})
	}
}

func TestErrorResponse_GetCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		err      *ErrorResponse
		expected string
	}{
		{
			name: "code at top level",
			err: &ErrorResponse{
				Code: "ERR_TOP",
			},
			expected: "ERR_TOP",
		},
		{
			name: "code in nested error",
			err: &ErrorResponse{
				Error: &ErrorDetail{
					Code: "ERR_NESTED",
				},
			},
			expected: "ERR_NESTED",
		},
		{
			name: "both codes (top level takes precedence)",
			err: &ErrorResponse{
				Code: "ERR_TOP",
				Error: &ErrorDetail{
					Code: "ERR_NESTED",
				},
			},
			expected: "ERR_TOP",
		},
		{
			name:     "no code",
			err:      &ErrorResponse{},
			expected: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, tt.err.GetCode())
		})
	}
}

func TestErrorResponse_JSON(t *testing.T) {
	t.Parallel()

	t.Run("unmarshal with nested error", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"error": {
				"message": "Invalid request",
				"code": "invalid_request_error",
				"type": "invalid_request",
				"param": "model"
			}
		}`

		var errResp ErrorResponse
		err := json.Unmarshal([]byte(jsonData), &errResp)
		require.NoError(t, err)

		assert.Equal(t, "Invalid request", errResp.Error.Message)
		assert.Equal(t, "invalid_request_error", errResp.Error.Code)
		assert.Equal(t, "invalid_request", errResp.Error.Type)
		assert.Equal(t, "model", errResp.Error.Param)
	})

	t.Run("unmarshal with top-level fields", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"message": "Rate limit exceeded",
			"code": "rate_limit_exceeded"
		}`

		var errResp ErrorResponse
		err := json.Unmarshal([]byte(jsonData), &errResp)
		require.NoError(t, err)

		assert.Equal(t, "Rate limit exceeded", errResp.Message)
		assert.Equal(t, "rate_limit_exceeded", errResp.Code)
	})

	t.Run("marshal and unmarshal", func(t *testing.T) {
		t.Parallel()

		original := &ErrorResponse{
			Error: &ErrorDetail{
				Message: "Test error",
				Code:    "test_error",
				Type:    "test",
				Param:   "test_param",
			},
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var decoded ErrorResponse
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		assert.Equal(t, original.Error.Message, decoded.Error.Message)
		assert.Equal(t, original.Error.Code, decoded.Error.Code)
		assert.Equal(t, original.Error.Type, decoded.Error.Type)
		assert.Equal(t, original.Error.Param, decoded.Error.Param)
	})
}

func TestBatchError_JSON(t *testing.T) {
	t.Parallel()

	jsonData := `{
		"code": "invalid_json",
		"line": 5,
		"message": "Invalid JSON at line 5",
		"param": "input"
	}`

	var batchErr BatchError
	err := json.Unmarshal([]byte(jsonData), &batchErr)
	require.NoError(t, err)

	assert.Equal(t, "invalid_json", batchErr.Code)
	assert.Equal(t, 5, batchErr.Line)
	assert.Equal(t, "Invalid JSON at line 5", batchErr.Message)
	assert.Equal(t, "input", batchErr.Param)
}

func TestAgentsError_JSON(t *testing.T) {
	t.Parallel()

	jsonData := `{
		"code": "agent_timeout",
		"message": "Agent execution timed out"
	}`

	var agentErr AgentsError
	err := json.Unmarshal([]byte(jsonData), &agentErr)
	require.NoError(t, err)

	assert.Equal(t, "agent_timeout", agentErr.Code)
	assert.Equal(t, "Agent execution timed out", agentErr.Message)
}

func TestErrorDetail_AllFields(t *testing.T) {
	t.Parallel()

	detail := &ErrorDetail{
		Message: "Test message",
		Code:    "test_code",
		Type:    "test_type",
		Param:   "test_param",
	}

	assert.Equal(t, "Test message", detail.Message)
	assert.Equal(t, "test_code", detail.Code)
	assert.Equal(t, "test_type", detail.Type)
	assert.Equal(t, "test_param", detail.Param)
}

func TestBatchError_AllFields(t *testing.T) {
	t.Parallel()

	batchErr := &BatchError{
		Code:    "batch_error",
		Line:    10,
		Message: "Batch processing error",
		Param:   "batch_param",
	}

	assert.Equal(t, "batch_error", batchErr.Code)
	assert.Equal(t, 10, batchErr.Line)
	assert.Equal(t, "Batch processing error", batchErr.Message)
	assert.Equal(t, "batch_param", batchErr.Param)
}

func TestAgentsError_AllFields(t *testing.T) {
	t.Parallel()

	agentErr := &AgentsError{
		Code:    "agent_error",
		Message: "Agent execution failed",
	}

	assert.Equal(t, "agent_error", agentErr.Code)
	assert.Equal(t, "Agent execution failed", agentErr.Message)
}
