package models

// ErrorResponse represents an error response from the API.
type ErrorResponse struct {
	// Error contains the error details.
	Error *ErrorDetail `json:"error,omitempty"`

	// Message is a direct error message (some APIs return this at top level).
	Message string `json:"message,omitempty"`

	// Code is a direct error code (some APIs return this at top level).
	Code string `json:"code,omitempty"`
}

// ErrorDetail contains detailed information about an error.
type ErrorDetail struct {
	// Message is the error message.
	Message string `json:"message,omitempty"`

	// Code is the error code.
	Code string `json:"code,omitempty"`

	// Type is the error type.
	Type string `json:"type,omitempty"`

	// Param is the parameter that caused the error.
	Param string `json:"param,omitempty"`
}

// GetMessage returns the error message, checking both top-level and nested error.
func (e *ErrorResponse) GetMessage() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Error != nil && e.Error.Message != "" {
		return e.Error.Message
	}
	return "Unknown error"
}

// GetCode returns the error code, checking both top-level and nested error.
func (e *ErrorResponse) GetCode() string {
	if e.Code != "" {
		return e.Code
	}
	if e.Error != nil && e.Error.Code != "" {
		return e.Error.Code
	}
	return ""
}

// BatchError represents an error in batch processing.
type BatchError struct {
	// Code is the error code.
	Code string `json:"code,omitempty"`

	// Line is the line number where the error occurred.
	Line int `json:"line,omitempty"`

	// Message is the error message.
	Message string `json:"message,omitempty"`

	// Param is the parameter that caused the error.
	Param string `json:"param,omitempty"`
}

// AgentsError represents an error from the agents API.
type AgentsError struct {
	// Code is the error code.
	Code string `json:"code,omitempty"`

	// Message is the error message.
	Message string `json:"message,omitempty"`
}
