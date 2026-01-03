package helpers_test

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/z-ai/zai-sdk-go/test/helpers"
)

func TestAssertNoError(t *testing.T) {
	t.Parallel()

	// This should not fail
	helpers.AssertNoError(t, nil)

	// Test that it fails when error is not nil
	// We can't directly test the failure, but we can ensure it doesn't panic
}

func TestAssertError(t *testing.T) {
	t.Parallel()

	// This should not fail
	helpers.AssertError(t, &mockError{msg: "test error"})
}

func TestAssertEqual(t *testing.T) {
	t.Parallel()

	helpers.AssertEqual(t, 1, 1)
	helpers.AssertEqual(t, "hello", "hello")
	helpers.AssertEqual(t, true, true)
}

func TestAssertNotEqual(t *testing.T) {
	t.Parallel()

	helpers.AssertNotEqual(t, 1, 2)
	helpers.AssertNotEqual(t, "hello", "world")
	helpers.AssertNotEqual(t, true, false)
}

func TestAssertTrue(t *testing.T) {
	t.Parallel()

	helpers.AssertTrue(t, true, "should be true")
	helpers.AssertTrue(t, 1 == 1, "one should equal one")
}

func TestAssertFalse(t *testing.T) {
	t.Parallel()

	helpers.AssertFalse(t, false, "should be false")
	helpers.AssertFalse(t, 1 == 2, "one should not equal two")
}

func TestAssertNil(t *testing.T) {
	t.Parallel()

	var nilError error
	helpers.AssertNil(t, nilError)
	helpers.AssertNil(t, nil)
}

func TestAssertNotNil(t *testing.T) {
	t.Parallel()

	err := &mockError{msg: "test"}
	helpers.AssertNotNil(t, err)
	helpers.AssertNotNil(t, "not nil")
	helpers.AssertNotNil(t, 123)
}

func TestCreateTestContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := helpers.CreateTestContext()
	defer cancel()

	helpers.AssertNotNil(t, ctx)

	// Verify context has a deadline
	deadline, ok := ctx.Deadline()
	helpers.AssertTrue(t, ok, "context should have a deadline")
	helpers.AssertTrue(t, deadline.After(time.Now()), "deadline should be in the future")
}

func TestCreateTestContextWithTimeout(t *testing.T) {
	t.Parallel()

	timeout := 2 * time.Second
	ctx, cancel := helpers.CreateTestContextWithTimeout(timeout)
	defer cancel()

	helpers.AssertNotNil(t, ctx)

	// Verify context has a deadline
	deadline, ok := ctx.Deadline()
	helpers.AssertTrue(t, ok, "context should have a deadline")

	// The deadline should be approximately timeout duration from now
	expectedDeadline := time.Now().Add(timeout)
	timeDiff := deadline.Sub(expectedDeadline).Abs()
	helpers.AssertTrue(t, timeDiff < 100*time.Millisecond, "deadline should be close to expected")
}

func TestCreateTestContextWithDeadline(t *testing.T) {
	t.Parallel()

	deadline := time.Now().Add(3 * time.Second)
	ctx, cancel := helpers.CreateTestContextWithDeadline(deadline)
	defer cancel()

	helpers.AssertNotNil(t, ctx)

	// Verify context has the correct deadline
	ctxDeadline, ok := ctx.Deadline()
	helpers.AssertTrue(t, ok, "context should have a deadline")
	helpers.AssertEqual(t, ctxDeadline.Unix(), deadline.Unix())
}

func TestNewTestServer(t *testing.T) {
	t.Parallel()

	server := helpers.NewTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})
	defer server.Close()

	helpers.AssertNotNil(t, server)
	helpers.AssertNotEqual(t, server.URL, "")

	// Make a request to the test server
	resp, err := http.Get(server.URL)
	helpers.AssertNoError(t, err)
	defer resp.Body.Close()

	helpers.AssertEqual(t, resp.StatusCode, http.StatusOK)
}

func TestNewTestTLSServer(t *testing.T) {
	t.Parallel()

	server := helpers.NewTestTLSServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("secure response"))
	})
	defer server.Close()

	helpers.AssertNotNil(t, server)
	helpers.AssertNotEqual(t, server.URL, "")

	// Make a request to the TLS test server
	client := server.Client()
	resp, err := client.Get(server.URL)
	helpers.AssertNoError(t, err)
	defer resp.Body.Close()

	helpers.AssertEqual(t, resp.StatusCode, http.StatusOK)
}

func TestCreateMockResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		mock helpers.MockHTTPResponse
	}{
		{
			name: "simple response",
			mock: helpers.MockHTTPResponse{
				StatusCode: 200,
				Body:       "OK",
			},
		},
		{
			name: "response with headers",
			mock: helpers.MockHTTPResponse{
				StatusCode: 201,
				Body:       `{"status": "created"}`,
				Headers: map[string]string{
					"Content-Type": "application/json",
					"X-Request-ID": "test-123",
				},
			},
		},
		{
			name: "error response",
			mock: helpers.MockHTTPResponse{
				StatusCode: 400,
				Body:       `{"error": "bad request"}`,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := helpers.CreateMockResponse(tt.mock)

			helpers.AssertNotNil(t, rec)
			helpers.AssertEqual(t, rec.Code, tt.mock.StatusCode)
			helpers.AssertEqual(t, rec.Body.String(), tt.mock.Body)

			for key, value := range tt.mock.Headers {
				helpers.AssertEqual(t, rec.Header().Get(key), value)
			}
		})
	}
}

func TestRunTableTests(t *testing.T) {
	t.Parallel()

	tests := []helpers.TestTable[int]{
		{
			Name:     "double positive",
			Input:    5,
			Expected: 10,
			WantErr:  false,
		},
		{
			Name:     "double zero",
			Input:    0,
			Expected: 0,
			WantErr:  false,
		},
		{
			Name:     "double negative",
			Input:    -3,
			Expected: -6,
			WantErr:  false,
		},
	}

	helpers.RunTableTests(t, tests, func(t *testing.T, tt helpers.TestTable[int]) {
		result := tt.Input * 2
		// Cast Expected to int for type-safe comparison
		expected, ok := tt.Expected.(int)
		if !ok {
			t.Fatal("Expected is not an int")
		}
		helpers.AssertEqual(t, result, expected)
	})
}

func TestSkipIfShort(t *testing.T) {
	// Create a sub-test to avoid skipping the entire test
	t.Run("skip in short mode", func(t *testing.T) {
		if !testing.Short() {
			// Only run this test when not in short mode
			helpers.SkipIfShort(t)
		}
	})
}

func TestTempEnv(t *testing.T) {
	// Note: Cannot use t.Parallel() with t.Setenv()

	key := "TEST_TEMP_ENV"
	value := "test-value"

	helpers.TempEnv(t, key, value)
	helpers.AssertEqual(t, os.Getenv(key), value)
}

func TestCleanup(t *testing.T) {
	t.Parallel()

	cleaned := false
	helpers.Cleanup(t, func() {
		cleaned = true
	})

	// The cleanup function will be called after the test completes
	// We can't verify it directly in the test, but we can ensure it registers
	helpers.AssertFalse(t, cleaned, "cleanup should not have run yet")
}

// Mock types for testing

type mockError struct {
	msg string
}

func (e *mockError) Error() string {
	return e.msg
}
