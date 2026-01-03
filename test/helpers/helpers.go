// Package helpers provides testing utilities for the Z.ai SDK.
package helpers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// AssertNoError is a helper that fails the test if err is not nil.
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// AssertError is a helper that fails the test if err is nil.
func AssertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

// AssertEqual is a helper that fails the test if got != want.
func AssertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

// AssertNotEqual is a helper that fails the test if got == want.
func AssertNotEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got == want {
		t.Errorf("got %v, want different value", got)
	}
}

// AssertTrue is a helper that fails the test if condition is false.
func AssertTrue(t *testing.T, condition bool, msg string) {
	t.Helper()
	if !condition {
		t.Errorf("assertion failed: %s", msg)
	}
}

// AssertFalse is a helper that fails the test if condition is true.
func AssertFalse(t *testing.T, condition bool, msg string) {
	t.Helper()
	if condition {
		t.Errorf("assertion failed: %s", msg)
	}
}

// AssertNil is a helper that fails the test if v is not nil.
func AssertNil(t *testing.T, v interface{}) {
	t.Helper()
	if v != nil {
		t.Errorf("expected nil, got %v", v)
	}
}

// AssertNotNil is a helper that fails the test if v is nil.
func AssertNotNil(t *testing.T, v interface{}) {
	t.Helper()
	if v == nil {
		t.Fatal("expected non-nil value")
	}
}

// CreateTestContext creates a context with a reasonable timeout for tests.
// The default timeout is 5 seconds.
func CreateTestContext() (context.Context, context.CancelFunc) {
	return CreateTestContextWithTimeout(5 * time.Second)
}

// CreateTestContextWithTimeout creates a context with a custom timeout.
func CreateTestContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// CreateTestContextWithDeadline creates a context with a specific deadline.
func CreateTestContextWithDeadline(deadline time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), deadline)
}

// NewTestServer creates a new httptest.Server for testing HTTP requests.
// The handler function receives the http.ResponseWriter and http.Request.
func NewTestServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

// NewTestTLSServer creates a new httptest.Server with TLS for testing HTTPS requests.
func NewTestTLSServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewTLSServer(handler)
}

// MockHTTPResponse creates a mock HTTP response for testing.
type MockHTTPResponse struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

// CreateMockResponse creates an httptest.ResponseRecorder with the given mock response.
func CreateMockResponse(mock MockHTTPResponse) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	rec.WriteHeader(mock.StatusCode)

	for key, value := range mock.Headers {
		rec.Header().Set(key, value)
	}

	if mock.Body != "" {
		rec.WriteString(mock.Body)
	}

	return rec
}

// TestTable represents a test case in a table-driven test.
type TestTable[T any] struct {
	Name     string
	Input    T
	Expected interface{}
	WantErr  bool
}

// RunTableTests runs a table-driven test with the given test cases.
// The testFunc is called for each test case and should perform the test logic.
func RunTableTests[T any](t *testing.T, tests []TestTable[T], testFunc func(*testing.T, TestTable[T])) {
	t.Helper()

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()
			testFunc(t, tt)
		})
	}
}

// SkipIfShort skips the test if the -short flag is set.
func SkipIfShort(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
}

// SkipIfCI skips the test if running in a CI environment.
func SkipIfCI(t *testing.T) {
	t.Helper()
	// Common CI environment variables
	ciVars := []string{"CI", "CONTINUOUS_INTEGRATION", "GITHUB_ACTIONS", "GITLAB_CI", "CIRCLECI"}
	for _, v := range ciVars {
		if os.Getenv(v) != "" {
			t.Skip("skipping test in CI environment")
			return
		}
	}
}

// TempEnv temporarily sets an environment variable for the duration of the test.
// The original value is restored after the test completes.
func TempEnv(t *testing.T, key, value string) {
	t.Helper()
	t.Setenv(key, value)
}

// Cleanup registers a cleanup function to be called when the test completes.
// This is a thin wrapper around t.Cleanup for consistency.
func Cleanup(t *testing.T, fn func()) {
	t.Helper()
	t.Cleanup(fn)
}
