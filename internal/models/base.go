// Package models provides core request/response models for the Z.ai SDK.
package models

import (
	"encoding/json"
	"time"
)

// Model represents a base model interface for all API models.
type Model interface {
	// Validate validates the model fields.
	Validate() error
}

// BaseModel provides common functionality for all models.
type BaseModel struct {
	// Extra holds additional fields not explicitly defined in the model.
	Extra map[string]interface{} `json:"-"`
}

// UnmarshalJSON implements custom JSON unmarshaling to capture extra fields.
func (b *BaseModel) UnmarshalJSON(data []byte) error {
	if b.Extra == nil {
		b.Extra = make(map[string]interface{})
	}
	return json.Unmarshal(data, &b.Extra)
}

// MarshalJSON implements custom JSON marshaling including extra fields.
func (b BaseModel) MarshalJSON() ([]byte, error) {
	if b.Extra == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(b.Extra)
}

// Get retrieves a value from extra fields.
func (b *BaseModel) Get(key string) (interface{}, bool) {
	if b.Extra == nil {
		return nil, false
	}
	val, ok := b.Extra[key]
	return val, ok
}

// Set stores a value in extra fields.
func (b *BaseModel) Set(key string, value interface{}) {
	if b.Extra == nil {
		b.Extra = make(map[string]interface{})
	}
	b.Extra[key] = value
}

// Keys returns all keys from extra fields.
func (b *BaseModel) Keys() []string {
	if b.Extra == nil {
		return []string{}
	}

	keys := make([]string, 0, len(b.Extra))
	for k := range b.Extra {
		keys = append(keys, k)
	}
	return keys
}

// ToMap converts the model to a map.
func (b *BaseModel) ToMap() map[string]interface{} {
	if b.Extra == nil {
		return make(map[string]interface{})
	}

	// Create a copy to avoid external mutations
	result := make(map[string]interface{}, len(b.Extra))
	for k, v := range b.Extra {
		result[k] = v
	}
	return result
}

// CommonRequestFields contains fields common to all API requests.
type CommonRequestFields struct {
	// Model is the model identifier (e.g., "glm-4", "glm-4-plus").
	Model string `json:"model,omitempty"`

	// Stream enables streaming responses.
	Stream *bool `json:"stream,omitempty"`

	// RequestID is a unique identifier for the request.
	RequestID string `json:"request_id,omitempty"`

	// UserID is the user identifier.
	UserID string `json:"user_id,omitempty"`
}

// CommonResponseFields contains fields common to all API responses.
type CommonResponseFields struct {
	// ID is a unique identifier for the response.
	ID string `json:"id,omitempty"`

	// Created is the Unix timestamp when the response was created.
	Created int64 `json:"created,omitempty"`

	// RequestID echoes the request ID if provided.
	RequestID string `json:"request_id,omitempty"`

	// Model is the model identifier used for the response.
	Model string `json:"model,omitempty"`
}

// GetCreatedTime returns the creation time as time.Time.
func (c *CommonResponseFields) GetCreatedTime() time.Time {
	if c.Created == 0 {
		return time.Time{}
	}
	return time.Unix(c.Created, 0)
}

// Validate validates common request fields.
func (c *CommonRequestFields) Validate() error {
	// Basic validation - model should not be empty for most requests
	// Individual request types can override this
	return nil
}
