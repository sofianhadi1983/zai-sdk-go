package helpers_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/sofianhadi1983/zai-sdk-go/test/helpers"
)

// This file demonstrates table-driven test patterns for the Z.ai Go SDK.
// Use this as a template for writing your own tests.

// Example 1: Simple table-driven test with basic assertions
func TestExampleSimpleTableDriven(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{
			name:     "positive number",
			input:    5,
			expected: 10,
		},
		{
			name:     "zero",
			input:    0,
			expected: 0,
		},
		{
			name:     "negative number",
			input:    -3,
			expected: -6,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Perform your test logic here
			result := tt.input * 2

			// Use helper assertions
			helpers.AssertEqual(t, result, tt.expected)
		})
	}
}

// Example 2: Table-driven test with error handling
func TestExampleWithErrorHandling(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid input",
			input:   "hello",
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
			errMsg:  "input cannot be empty",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Simulated function that might return an error
			var err error
			if tt.input == "" {
				err = &simpleError{msg: "input cannot be empty"}
			}

			if tt.wantErr {
				helpers.AssertError(t, err)
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("error message = %q, want %q", err.Error(), tt.errMsg)
				}
			} else {
				helpers.AssertNoError(t, err)
			}
		})
	}
}

// Example 3: Using the generic TestTable helper
func TestExampleUsingTestTableHelper(t *testing.T) {
	t.Parallel()

	tests := []helpers.TestTable[string]{
		{
			Name:     "uppercase",
			Input:    "hello",
			Expected: "HELLO",
			WantErr:  false,
		},
		{
			Name:     "already uppercase",
			Input:    "WORLD",
			Expected: "WORLD",
			WantErr:  false,
		},
	}

	helpers.RunTableTests(t, tests, func(t *testing.T, tt helpers.TestTable[string]) {
		// Simulated uppercase function
		result := mockUppercase(tt.Input)

		// Cast Expected to string for type-safe comparison
		expected, ok := tt.Expected.(string)
		if !ok {
			t.Fatal("Expected is not a string")
		}
		helpers.AssertEqual(t, result, expected)
	})
}

// Example 4: Test with context and timeout
func TestExampleWithContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := helpers.CreateTestContext()
	defer cancel()

	// Use the context in your test
	select {
	case <-ctx.Done():
		t.Fatal("context cancelled unexpectedly")
	default:
		// Test logic here
		helpers.AssertTrue(t, true, "context is active")
	}
}

// Example 5: Test with mock HTTP server
func TestExampleWithMockServer(t *testing.T) {
	t.Parallel()

	// Create a mock server
	server := helpers.NewTestServer(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		helpers.AssertEqual(t, r.Method, "GET")
		helpers.AssertEqual(t, r.Header.Get("Authorization"), "Bearer "+helpers.MockAPIKey)

		// Send response
		w.WriteHeader(200)
		w.Write([]byte(`{"status": "ok"}`))
	})
	defer server.Close()

	// Your test logic using server.URL
	helpers.AssertNotNil(t, server)
}

// Example 6: Test with temporary environment variable
func TestExampleWithTempEnv(t *testing.T) {
	// Note: Cannot use t.Parallel() when using t.Setenv()

	newValue := "temporary"

	// Set temporary environment variable
	helpers.TempEnv(t, "TEST_ENV_VAR", newValue)

	// Verify the value
	helpers.AssertEqual(t, os.Getenv("TEST_ENV_VAR"), newValue)

	// The original value will be restored automatically after the test
}

// Example 7: Test with cleanup function
func TestExampleWithCleanup(t *testing.T) {
	t.Parallel()

	// Simulated resource allocation
	resource := &mockResource{name: "test"}

	// Register cleanup
	helpers.Cleanup(t, func() {
		resource.Close()
	})

	// Use the resource
	helpers.AssertNotNil(t, resource)
}

// Example 8: Skip test in short mode
func TestExampleSkipInShortMode(t *testing.T) {
	helpers.SkipIfShort(t)

	// This test will be skipped when running with -short flag
	// Useful for slow integration tests
	t.Log("Running slow test...")
}

// Example 9: Test with mock response
func TestExampleWithMockResponse(t *testing.T) {
	t.Parallel()

	mock := helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       helpers.MockChatCompletionResponse(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	rec := helpers.CreateMockResponse(mock)

	helpers.AssertEqual(t, rec.Code, 200)
	helpers.AssertNotNil(t, rec.Body)
}

// Example 10: Comprehensive example combining multiple patterns
func TestExampleComprehensive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      map[string]interface{}
		mockStatus int
		mockBody   string
		wantErr    bool
	}{
		{
			name:       "successful chat completion",
			input:      helpers.MockChatCompletionRequest(),
			mockStatus: 200,
			mockBody:   helpers.MockChatCompletionResponse(),
			wantErr:    false,
		},
		{
			name:       "authentication error",
			input:      helpers.MockChatCompletionRequest(),
			mockStatus: 401,
			mockBody:   helpers.MockAuthErrorResponse(),
			wantErr:    true,
		},
		{
			name:       "rate limit error",
			input:      helpers.MockChatCompletionRequest(),
			mockStatus: 429,
			mockBody:   helpers.MockRateLimitErrorResponse(),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create test context
			ctx, cancel := helpers.CreateTestContext()
			defer cancel()

			// Create mock server
			server := helpers.NewTestServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockBody))
			})
			defer server.Close()

			// Perform your API call here using server.URL
			// For now, just verify the setup
			helpers.AssertNotNil(t, ctx)
			helpers.AssertNotNil(t, server)
			helpers.AssertNotNil(t, tt.input)
		})
	}
}

// Mock types for examples

type simpleError struct {
	msg string
}

func (e *simpleError) Error() string {
	return e.msg
}

type mockResource struct {
	name   string
	closed bool
}

func (r *mockResource) Close() {
	r.closed = true
}

func mockUppercase(s string) string {
	// Simplified uppercase implementation
	result := ""
	for _, c := range s {
		if c >= 'a' && c <= 'z' {
			result += string(c - 32)
		} else {
			result += string(c)
		}
	}
	return result
}
