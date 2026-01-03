// Package errors provides error types for the Z.ai SDK.
package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// ZaiError is the base error type for all Z.ai SDK errors.
type ZaiError struct {
	Message string
}

// Error implements the error interface for ZaiError.
func (e *ZaiError) Error() string {
	return e.Message
}

// NewZaiError creates a new ZaiError with the given message.
func NewZaiError(message string) *ZaiError {
	return &ZaiError{Message: message}
}

// APIStatusError represents an API error with an HTTP status code.
// This is the base type for all API errors that have an HTTP response.
type APIStatusError struct {
	*ZaiError
	StatusCode int
	Response   *http.Response
	RequestID  string // Optional request ID for tracing
}

// Error implements the error interface for APIStatusError.
func (e *APIStatusError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("API error (status %d, request_id: %s): %s", e.StatusCode, e.RequestID, e.Message)
	}
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}

// Unwrap implements error unwrapping for APIStatusError.
func (e *APIStatusError) Unwrap() error {
	return e.ZaiError
}

// NewAPIStatusError creates a new APIStatusError.
func NewAPIStatusError(message string, statusCode int, response *http.Response) *APIStatusError {
	return &APIStatusError{
		ZaiError:   &ZaiError{Message: message},
		StatusCode: statusCode,
		Response:   response,
	}
}

// APIRequestFailedError indicates a general API request failure.
type APIRequestFailedError struct {
	*APIStatusError
}

// Unwrap implements error unwrapping for APIRequestFailedError.
func (e *APIRequestFailedError) Unwrap() error {
	return e.APIStatusError
}

// NewAPIRequestFailedError creates a new APIRequestFailedError.
func NewAPIRequestFailedError(message string, statusCode int, response *http.Response) *APIRequestFailedError {
	return &APIRequestFailedError{
		APIStatusError: NewAPIStatusError(message, statusCode, response),
	}
}

// APIAuthenticationError indicates an authentication failure (401).
type APIAuthenticationError struct {
	*APIStatusError
}

// Unwrap implements error unwrapping for APIAuthenticationError.
func (e *APIAuthenticationError) Unwrap() error {
	return e.APIStatusError
}

// NewAPIAuthenticationError creates a new APIAuthenticationError.
func NewAPIAuthenticationError(message string, statusCode int, response *http.Response) *APIAuthenticationError {
	return &APIAuthenticationError{
		APIStatusError: NewAPIStatusError(message, statusCode, response),
	}
}

// APIReachLimitError indicates a rate limit has been exceeded (429).
type APIReachLimitError struct {
	*APIStatusError
	RetryAfter int // Seconds to wait before retrying
}

// Unwrap implements error unwrapping for APIReachLimitError.
func (e *APIReachLimitError) Unwrap() error {
	return e.APIStatusError
}

// NewAPIReachLimitError creates a new APIReachLimitError.
func NewAPIReachLimitError(message string, statusCode int, response *http.Response) *APIReachLimitError {
	return &APIReachLimitError{
		APIStatusError: NewAPIStatusError(message, statusCode, response),
	}
}

// APIInternalError indicates an internal server error (500).
type APIInternalError struct {
	*APIStatusError
}

// Unwrap implements error unwrapping for APIInternalError.
func (e *APIInternalError) Unwrap() error {
	return e.APIStatusError
}

// NewAPIInternalError creates a new APIInternalError.
func NewAPIInternalError(message string, statusCode int, response *http.Response) *APIInternalError {
	return &APIInternalError{
		APIStatusError: NewAPIStatusError(message, statusCode, response),
	}
}

// APIServerFlowExceedError indicates server flow has been exceeded (503).
type APIServerFlowExceedError struct {
	*APIStatusError
}

// Unwrap implements error unwrapping for APIServerFlowExceedError.
func (e *APIServerFlowExceedError) Unwrap() error {
	return e.APIStatusError
}

// NewAPIServerFlowExceedError creates a new APIServerFlowExceedError.
func NewAPIServerFlowExceedError(message string, statusCode int, response *http.Response) *APIServerFlowExceedError {
	return &APIServerFlowExceedError{
		APIStatusError: NewAPIStatusError(message, statusCode, response),
	}
}

// APIResponseError represents an error related to API response handling.
type APIResponseError struct {
	*ZaiError
	Request  *http.Request
	JSONData interface{} // Parsed JSON response data
}

// Error implements the error interface for APIResponseError.
func (e *APIResponseError) Error() string {
	if e.Request != nil {
		return fmt.Sprintf("API response error for %s %s: %s", e.Request.Method, e.Request.URL.Path, e.Message)
	}
	return fmt.Sprintf("API response error: %s", e.Message)
}

// Unwrap implements error unwrapping for APIResponseError.
func (e *APIResponseError) Unwrap() error {
	return e.ZaiError
}

// NewAPIResponseError creates a new APIResponseError.
func NewAPIResponseError(message string, request *http.Request, jsonData interface{}) *APIResponseError {
	return &APIResponseError{
		ZaiError: &ZaiError{Message: message},
		Request:  request,
		JSONData: jsonData,
	}
}

// APIResponseValidationError indicates the API response failed validation.
type APIResponseValidationError struct {
	*APIResponseError
	StatusCode int
	Response   *http.Response
}

// Error implements the error interface for APIResponseValidationError.
func (e *APIResponseValidationError) Error() string {
	if e.Request != nil {
		return fmt.Sprintf("API response validation error (status %d) for %s %s: %s",
			e.StatusCode, e.Request.Method, e.Request.URL.Path, e.Message)
	}
	return fmt.Sprintf("API response validation error (status %d): %s", e.StatusCode, e.Message)
}

// Unwrap implements error unwrapping for APIResponseValidationError.
func (e *APIResponseValidationError) Unwrap() error {
	return e.APIResponseError
}

// NewAPIResponseValidationError creates a new APIResponseValidationError.
func NewAPIResponseValidationError(response *http.Response, jsonData interface{}, message string) *APIResponseValidationError {
	if message == "" {
		message = "Data returned by API invalid for expected schema."
	}

	var request *http.Request
	if response != nil {
		request = response.Request
	}

	return &APIResponseValidationError{
		APIResponseError: NewAPIResponseError(message, request, jsonData),
		StatusCode:       response.StatusCode,
		Response:         response,
	}
}

// APIConnectionError indicates a connection error occurred.
type APIConnectionError struct {
	*APIResponseError
}

// Unwrap implements error unwrapping for APIConnectionError.
func (e *APIConnectionError) Unwrap() error {
	return e.APIResponseError
}

// NewAPIConnectionError creates a new APIConnectionError.
func NewAPIConnectionError(request *http.Request, message string) *APIConnectionError {
	if message == "" {
		message = "Connection error."
	}
	return &APIConnectionError{
		APIResponseError: NewAPIResponseError(message, request, nil),
	}
}

// APITimeoutError indicates a request timeout occurred.
type APITimeoutError struct {
	*APIConnectionError
}

// Unwrap implements error unwrapping for APITimeoutError.
func (e *APITimeoutError) Unwrap() error {
	return e.APIConnectionError
}

// NewAPITimeoutError creates a new APITimeoutError.
func NewAPITimeoutError(request *http.Request) *APITimeoutError {
	return &APITimeoutError{
		APIConnectionError: NewAPIConnectionError(request, "Request timed out."),
	}
}

// ConfigError represents a configuration error.
type ConfigError struct {
	*ZaiError
	Field string // The configuration field that caused the error
}

// Error implements the error interface for ConfigError.
func (e *ConfigError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("configuration error for field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("configuration error: %s", e.Message)
}

// Unwrap implements error unwrapping for ConfigError.
func (e *ConfigError) Unwrap() error {
	return e.ZaiError
}

// NewConfigError creates a new ConfigError.
func NewConfigError(field, message string) *ConfigError {
	return &ConfigError{
		ZaiError: &ZaiError{Message: message},
		Field:    field,
	}
}

// ValidationError represents an input validation error.
type ValidationError struct {
	*ZaiError
	Field string      // The field that failed validation
	Value interface{} // The invalid value
}

// Error implements the error interface for ValidationError.
func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// Unwrap implements error unwrapping for ValidationError.
func (e *ValidationError) Unwrap() error {
	return e.ZaiError
}

// NewValidationError creates a new ValidationError.
func NewValidationError(field, message string, value interface{}) *ValidationError {
	return &ValidationError{
		ZaiError: &ZaiError{Message: message},
		Field:    field,
		Value:    value,
	}
}

// Error type checking helpers

// IsAuthenticationError checks if the error is an authentication error.
func IsAuthenticationError(err error) bool {
	var authErr *APIAuthenticationError
	return errors.As(err, &authErr)
}

// IsRateLimitError checks if the error is a rate limit error.
func IsRateLimitError(err error) bool {
	var rateLimitErr *APIReachLimitError
	return errors.As(err, &rateLimitErr)
}

// IsServerError checks if the error is a server error (5xx).
func IsServerError(err error) bool {
	var internalErr *APIInternalError
	var flowErr *APIServerFlowExceedError
	return errors.As(err, &internalErr) || errors.As(err, &flowErr)
}

// IsRequestError checks if the error is a request error (4xx).
func IsRequestError(err error) bool {
	var requestErr *APIRequestFailedError
	return errors.As(err, &requestErr)
}

// IsConnectionError checks if the error is a connection error.
func IsConnectionError(err error) bool {
	var connErr *APIConnectionError
	return errors.As(err, &connErr)
}

// IsTimeoutError checks if the error is a timeout error.
func IsTimeoutError(err error) bool {
	var timeoutErr *APITimeoutError
	return errors.As(err, &timeoutErr)
}

// IsConfigError checks if the error is a configuration error.
func IsConfigError(err error) bool {
	var configErr *ConfigError
	return errors.As(err, &configErr)
}

// IsValidationError checks if the error is a validation error.
func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}
