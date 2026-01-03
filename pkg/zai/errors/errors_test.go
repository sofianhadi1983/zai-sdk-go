package errors

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestZaiError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "simple error",
			message: "something went wrong",
		},
		{
			name:    "empty message",
			message: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := NewZaiError(tt.message)
			if err.Error() != tt.message {
				t.Errorf("ZaiError.Error() = %q, want %q", err.Error(), tt.message)
			}

			// Test that it implements error interface
			var _ error = err
		})
	}
}

func TestAPIStatusError(t *testing.T) {
	t.Parallel()

	// Create a test HTTP response
	rec := httptest.NewRecorder()
	rec.WriteHeader(http.StatusBadRequest)
	resp := rec.Result()

	tests := []struct {
		name       string
		message    string
		statusCode int
		requestID  string
		wantErr    string
	}{
		{
			name:       "error without request ID",
			message:    "bad request",
			statusCode: 400,
			requestID:  "",
			wantErr:    "API error (status 400): bad request",
		},
		{
			name:       "error with request ID",
			message:    "bad request",
			statusCode: 400,
			requestID:  "req-123",
			wantErr:    "API error (status 400, request_id: req-123): bad request",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := NewAPIStatusError(tt.message, tt.statusCode, resp)
			err.RequestID = tt.requestID

			if err.Error() != tt.wantErr {
				t.Errorf("APIStatusError.Error() = %q, want %q", err.Error(), tt.wantErr)
			}

			if err.StatusCode != tt.statusCode {
				t.Errorf("APIStatusError.StatusCode = %d, want %d", err.StatusCode, tt.statusCode)
			}

			// Test error unwrapping
			var baseErr *ZaiError
			if !errors.As(err, &baseErr) {
				t.Error("APIStatusError should unwrap to ZaiError")
			}
		})
	}
}

func TestAPIRequestFailedError(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	rec.WriteHeader(http.StatusBadRequest)
	resp := rec.Result()

	err := NewAPIRequestFailedError("request failed", 400, resp)

	if err.StatusCode != 400 {
		t.Errorf("APIRequestFailedError.StatusCode = %d, want 400", err.StatusCode)
	}

	if !strings.Contains(err.Error(), "request failed") {
		t.Errorf("APIRequestFailedError.Error() = %q, should contain 'request failed'", err.Error())
	}

	// Test type assertion
	var statusErr *APIStatusError
	if !errors.As(err, &statusErr) {
		t.Error("APIRequestFailedError should unwrap to APIStatusError")
	}
}

func TestAPIAuthenticationError(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	rec.WriteHeader(http.StatusUnauthorized)
	resp := rec.Result()

	err := NewAPIAuthenticationError("invalid API key", 401, resp)

	if err.StatusCode != 401 {
		t.Errorf("APIAuthenticationError.StatusCode = %d, want 401", err.StatusCode)
	}

	if !strings.Contains(err.Error(), "invalid API key") {
		t.Errorf("APIAuthenticationError.Error() = %q, should contain 'invalid API key'", err.Error())
	}

	// Test type assertion
	var statusErr *APIStatusError
	if !errors.As(err, &statusErr) {
		t.Error("APIAuthenticationError should unwrap to APIStatusError")
	}
}

func TestAPIReachLimitError(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	rec.WriteHeader(http.StatusTooManyRequests)
	resp := rec.Result()

	err := NewAPIReachLimitError("rate limit exceeded", 429, resp)
	err.RetryAfter = 60

	if err.StatusCode != 429 {
		t.Errorf("APIReachLimitError.StatusCode = %d, want 429", err.StatusCode)
	}

	if err.RetryAfter != 60 {
		t.Errorf("APIReachLimitError.RetryAfter = %d, want 60", err.RetryAfter)
	}

	if !strings.Contains(err.Error(), "rate limit exceeded") {
		t.Errorf("APIReachLimitError.Error() = %q, should contain 'rate limit exceeded'", err.Error())
	}
}

func TestAPIInternalError(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	rec.WriteHeader(http.StatusInternalServerError)
	resp := rec.Result()

	err := NewAPIInternalError("internal server error", 500, resp)

	if err.StatusCode != 500 {
		t.Errorf("APIInternalError.StatusCode = %d, want 500", err.StatusCode)
	}

	if !strings.Contains(err.Error(), "internal server error") {
		t.Errorf("APIInternalError.Error() = %q, should contain 'internal server error'", err.Error())
	}
}

func TestAPIServerFlowExceedError(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	rec.WriteHeader(http.StatusServiceUnavailable)
	resp := rec.Result()

	err := NewAPIServerFlowExceedError("server flow exceeded", 503, resp)

	if err.StatusCode != 503 {
		t.Errorf("APIServerFlowExceedError.StatusCode = %d, want 503", err.StatusCode)
	}

	if !strings.Contains(err.Error(), "server flow exceeded") {
		t.Errorf("APIServerFlowExceedError.Error() = %q, should contain 'server flow exceeded'", err.Error())
	}
}

func TestAPIResponseError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		request  *http.Request
		message  string
		jsonData interface{}
		wantErr  string
	}{
		{
			name:     "error with request",
			request:  httptest.NewRequest("POST", "/api/v1/test", nil),
			message:  "invalid response",
			jsonData: map[string]string{"error": "test"},
			wantErr:  "API response error for POST /api/v1/test: invalid response",
		},
		{
			name:     "error without request",
			request:  nil,
			message:  "invalid response",
			jsonData: nil,
			wantErr:  "API response error: invalid response",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := NewAPIResponseError(tt.message, tt.request, tt.jsonData)

			if err.Error() != tt.wantErr {
				t.Errorf("APIResponseError.Error() = %q, want %q", err.Error(), tt.wantErr)
			}

			if err.Request != tt.request {
				t.Error("APIResponseError.Request not set correctly")
			}

			// Test error unwrapping
			var baseErr *ZaiError
			if !errors.As(err, &baseErr) {
				t.Error("APIResponseError should unwrap to ZaiError")
			}
		})
	}
}

func TestAPIResponseValidationError(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	rec.WriteHeader(http.StatusOK)
	resp := rec.Result()
	resp.Request = httptest.NewRequest("GET", "/api/v1/test", nil)

	tests := []struct {
		name     string
		response *http.Response
		jsonData interface{}
		message  string
		wantErr  string
	}{
		{
			name:     "with custom message",
			response: resp,
			jsonData: map[string]string{"invalid": "data"},
			message:  "custom validation error",
			wantErr:  "API response validation error (status 200) for GET /api/v1/test: custom validation error",
		},
		{
			name:     "with default message",
			response: resp,
			jsonData: nil,
			message:  "",
			wantErr:  "API response validation error (status 200) for GET /api/v1/test: Data returned by API invalid for expected schema.",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := NewAPIResponseValidationError(tt.response, tt.jsonData, tt.message)

			if err.Error() != tt.wantErr {
				t.Errorf("APIResponseValidationError.Error() = %q, want %q", err.Error(), tt.wantErr)
			}

			if err.StatusCode != tt.response.StatusCode {
				t.Errorf("APIResponseValidationError.StatusCode = %d, want %d", err.StatusCode, tt.response.StatusCode)
			}

			// Test error unwrapping
			var respErr *APIResponseError
			if !errors.As(err, &respErr) {
				t.Error("APIResponseValidationError should unwrap to APIResponseError")
			}
		})
	}
}

func TestAPIConnectionError(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/api/v1/test", nil)

	tests := []struct {
		name    string
		message string
		wantErr string
	}{
		{
			name:    "with custom message",
			message: "connection refused",
			wantErr: "API response error for GET /api/v1/test: connection refused",
		},
		{
			name:    "with default message",
			message: "",
			wantErr: "API response error for GET /api/v1/test: Connection error.",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := NewAPIConnectionError(req, tt.message)

			if err.Error() != tt.wantErr {
				t.Errorf("APIConnectionError.Error() = %q, want %q", err.Error(), tt.wantErr)
			}

			// Test error unwrapping
			var respErr *APIResponseError
			if !errors.As(err, &respErr) {
				t.Error("APIConnectionError should unwrap to APIResponseError")
			}
		})
	}
}

func TestAPITimeoutError(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/api/v1/test", nil)
	err := NewAPITimeoutError(req)

	wantErr := "API response error for GET /api/v1/test: Request timed out."
	if err.Error() != wantErr {
		t.Errorf("APITimeoutError.Error() = %q, want %q", err.Error(), wantErr)
	}

	// Test error unwrapping
	var connErr *APIConnectionError
	if !errors.As(err, &connErr) {
		t.Error("APITimeoutError should unwrap to APIConnectionError")
	}

	var respErr *APIResponseError
	if !errors.As(err, &respErr) {
		t.Error("APITimeoutError should unwrap to APIResponseError")
	}
}

func TestConfigError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		field   string
		message string
		wantErr string
	}{
		{
			name:    "with field",
			field:   "api_key",
			message: "API key is required",
			wantErr: "configuration error for field 'api_key': API key is required",
		},
		{
			name:    "without field",
			field:   "",
			message: "invalid configuration",
			wantErr: "configuration error: invalid configuration",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := NewConfigError(tt.field, tt.message)

			if err.Error() != tt.wantErr {
				t.Errorf("ConfigError.Error() = %q, want %q", err.Error(), tt.wantErr)
			}

			if err.Field != tt.field {
				t.Errorf("ConfigError.Field = %q, want %q", err.Field, tt.field)
			}

			// Test error unwrapping
			var baseErr *ZaiError
			if !errors.As(err, &baseErr) {
				t.Error("ConfigError should unwrap to ZaiError")
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		field   string
		message string
		value   interface{}
		wantErr string
	}{
		{
			name:    "with field and value",
			field:   "max_tokens",
			message: "must be positive",
			value:   -1,
			wantErr: "validation error for field 'max_tokens': must be positive",
		},
		{
			name:    "without field",
			field:   "",
			message: "invalid input",
			value:   nil,
			wantErr: "validation error: invalid input",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := NewValidationError(tt.field, tt.message, tt.value)

			if err.Error() != tt.wantErr {
				t.Errorf("ValidationError.Error() = %q, want %q", err.Error(), tt.wantErr)
			}

			if err.Field != tt.field {
				t.Errorf("ValidationError.Field = %q, want %q", err.Field, tt.field)
			}

			if err.Value != tt.value {
				t.Errorf("ValidationError.Value = %v, want %v", err.Value, tt.value)
			}

			// Test error unwrapping
			var baseErr *ZaiError
			if !errors.As(err, &baseErr) {
				t.Error("ValidationError should unwrap to ZaiError")
			}
		})
	}
}

func TestErrorInterface(t *testing.T) {
	t.Parallel()

	// Ensure all error types implement the error interface
	var _ error = &ZaiError{}
	var _ error = &APIStatusError{}
	var _ error = &APIRequestFailedError{}
	var _ error = &APIAuthenticationError{}
	var _ error = &APIReachLimitError{}
	var _ error = &APIInternalError{}
	var _ error = &APIServerFlowExceedError{}
	var _ error = &APIResponseError{}
	var _ error = &APIResponseValidationError{}
	var _ error = &APIConnectionError{}
	var _ error = &APITimeoutError{}
	var _ error = &ConfigError{}
	var _ error = &ValidationError{}
}

func TestErrorWrapping(t *testing.T) {
	t.Parallel()

	t.Run("APIStatusError wraps ZaiError", func(t *testing.T) {
		t.Parallel()

		rec := httptest.NewRecorder()
		rec.WriteHeader(http.StatusBadRequest)
		resp := rec.Result()

		err := NewAPIStatusError("test", 400, resp)

		var baseErr *ZaiError
		if !errors.As(err, &baseErr) {
			t.Error("APIStatusError should wrap ZaiError")
		}
	})

	t.Run("nested error types", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest("GET", "/test", nil)
		err := NewAPITimeoutError(req)

		// Should unwrap to APIConnectionError
		var connErr *APIConnectionError
		if !errors.As(err, &connErr) {
			t.Error("APITimeoutError should unwrap to APIConnectionError")
		}

		// Should unwrap to APIResponseError
		var respErr *APIResponseError
		if !errors.As(err, &respErr) {
			t.Error("APITimeoutError should unwrap to APIResponseError")
		}

		// Should unwrap to ZaiError
		var baseErr *ZaiError
		if !errors.As(err, &baseErr) {
			t.Error("APITimeoutError should unwrap to ZaiError")
		}
	})
}

func TestErrorHelpers(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	rec.WriteHeader(http.StatusUnauthorized)
	resp := rec.Result()

	t.Run("IsAuthenticationError", func(t *testing.T) {
		t.Parallel()

		authErr := NewAPIAuthenticationError("auth failed", 401, resp)
		if !IsAuthenticationError(authErr) {
			t.Error("IsAuthenticationError should return true for APIAuthenticationError")
		}

		otherErr := NewAPIInternalError("internal error", 500, resp)
		if IsAuthenticationError(otherErr) {
			t.Error("IsAuthenticationError should return false for non-authentication errors")
		}

		if IsAuthenticationError(nil) {
			t.Error("IsAuthenticationError should return false for nil error")
		}
	})

	t.Run("IsRateLimitError", func(t *testing.T) {
		t.Parallel()

		rec := httptest.NewRecorder()
		rec.WriteHeader(http.StatusTooManyRequests)
		resp := rec.Result()

		rateLimitErr := NewAPIReachLimitError("rate limit exceeded", 429, resp)
		if !IsRateLimitError(rateLimitErr) {
			t.Error("IsRateLimitError should return true for APIReachLimitError")
		}

		otherErr := NewAPIInternalError("internal error", 500, resp)
		if IsRateLimitError(otherErr) {
			t.Error("IsRateLimitError should return false for non-rate-limit errors")
		}

		if IsRateLimitError(nil) {
			t.Error("IsRateLimitError should return false for nil error")
		}
	})

	t.Run("IsServerError", func(t *testing.T) {
		t.Parallel()

		rec := httptest.NewRecorder()
		rec.WriteHeader(http.StatusInternalServerError)
		resp := rec.Result()

		internalErr := NewAPIInternalError("internal error", 500, resp)
		if !IsServerError(internalErr) {
			t.Error("IsServerError should return true for APIInternalError")
		}

		rec2 := httptest.NewRecorder()
		rec2.WriteHeader(http.StatusServiceUnavailable)
		resp2 := rec2.Result()

		flowErr := NewAPIServerFlowExceedError("flow exceeded", 503, resp2)
		if !IsServerError(flowErr) {
			t.Error("IsServerError should return true for APIServerFlowExceedError")
		}

		rec3 := httptest.NewRecorder()
		rec3.WriteHeader(http.StatusBadRequest)
		resp3 := rec3.Result()

		otherErr := NewAPIRequestFailedError("bad request", 400, resp3)
		if IsServerError(otherErr) {
			t.Error("IsServerError should return false for non-server errors")
		}

		if IsServerError(nil) {
			t.Error("IsServerError should return false for nil error")
		}
	})

	t.Run("IsRequestError", func(t *testing.T) {
		t.Parallel()

		rec := httptest.NewRecorder()
		rec.WriteHeader(http.StatusBadRequest)
		resp := rec.Result()

		requestErr := NewAPIRequestFailedError("bad request", 400, resp)
		if !IsRequestError(requestErr) {
			t.Error("IsRequestError should return true for APIRequestFailedError")
		}

		rec2 := httptest.NewRecorder()
		rec2.WriteHeader(http.StatusInternalServerError)
		resp2 := rec2.Result()

		otherErr := NewAPIInternalError("internal error", 500, resp2)
		if IsRequestError(otherErr) {
			t.Error("IsRequestError should return false for non-request errors")
		}

		if IsRequestError(nil) {
			t.Error("IsRequestError should return false for nil error")
		}
	})

	t.Run("IsConnectionError", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest("GET", "/test", nil)

		connErr := NewAPIConnectionError(req, "connection failed")
		if !IsConnectionError(connErr) {
			t.Error("IsConnectionError should return true for APIConnectionError")
		}

		// APITimeoutError is also a connection error
		timeoutErr := NewAPITimeoutError(req)
		if !IsConnectionError(timeoutErr) {
			t.Error("IsConnectionError should return true for APITimeoutError")
		}

		rec := httptest.NewRecorder()
		rec.WriteHeader(http.StatusBadRequest)
		resp := rec.Result()

		otherErr := NewAPIRequestFailedError("bad request", 400, resp)
		if IsConnectionError(otherErr) {
			t.Error("IsConnectionError should return false for non-connection errors")
		}

		if IsConnectionError(nil) {
			t.Error("IsConnectionError should return false for nil error")
		}
	})

	t.Run("IsTimeoutError", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest("GET", "/test", nil)

		timeoutErr := NewAPITimeoutError(req)
		if !IsTimeoutError(timeoutErr) {
			t.Error("IsTimeoutError should return true for APITimeoutError")
		}

		connErr := NewAPIConnectionError(req, "connection failed")
		if IsTimeoutError(connErr) {
			t.Error("IsTimeoutError should return false for non-timeout connection errors")
		}

		if IsTimeoutError(nil) {
			t.Error("IsTimeoutError should return false for nil error")
		}
	})

	t.Run("IsConfigError", func(t *testing.T) {
		t.Parallel()

		configErr := NewConfigError("api_key", "API key is required")
		if !IsConfigError(configErr) {
			t.Error("IsConfigError should return true for ConfigError")
		}

		otherErr := NewZaiError("generic error")
		if IsConfigError(otherErr) {
			t.Error("IsConfigError should return false for non-config errors")
		}

		if IsConfigError(nil) {
			t.Error("IsConfigError should return false for nil error")
		}
	})

	t.Run("IsValidationError", func(t *testing.T) {
		t.Parallel()

		validationErr := NewValidationError("max_tokens", "must be positive", -1)
		if !IsValidationError(validationErr) {
			t.Error("IsValidationError should return true for ValidationError")
		}

		otherErr := NewZaiError("generic error")
		if IsValidationError(otherErr) {
			t.Error("IsValidationError should return false for non-validation errors")
		}

		if IsValidationError(nil) {
			t.Error("IsValidationError should return false for nil error")
		}
	})
}
