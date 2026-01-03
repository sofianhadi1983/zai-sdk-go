package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds the SDK configuration loaded from environment variables.
type Config struct {
	// APIKey is the Z.ai API key for authentication.
	// Required. Set via ZAI_API_KEY environment variable.
	APIKey string

	// BaseURL is the base URL for the Z.ai API.
	// Optional. Set via ZAI_BASE_URL environment variable.
	// Defaults to the appropriate URL based on client type.
	BaseURL string

	// Timeout is the HTTP request timeout duration.
	// Optional. Set via ZAI_TIMEOUT environment variable (in seconds).
	// Defaults to 120 seconds.
	Timeout time.Duration

	// MaxRetries is the maximum number of retries for failed requests.
	// Optional. Set via ZAI_MAX_RETRIES environment variable.
	// Defaults to 3.
	MaxRetries int

	// DisableTokenCache controls whether JWT token caching is disabled.
	// Optional. Set via ZAI_DISABLE_TOKEN_CACHE environment variable.
	// Defaults to true (caching disabled).
	DisableTokenCache bool

	// SourceChannel identifies the source of the API requests.
	// Optional. Set via ZAI_SOURCE_CHANNEL environment variable.
	// Defaults to "go-sdk".
	SourceChannel string
}

// Default configuration values
const (
	DefaultTimeout           = 120 * time.Second
	DefaultMaxRetries        = 3
	DefaultDisableTokenCache = true
	DefaultSourceChannel     = "go-sdk"
)

// LoadConfig loads configuration from environment variables.
// It returns an error if required configuration (API key) is missing.
//
// Environment Variables:
//   - ZAI_API_KEY: Required. API key for authentication.
//   - ZAI_BASE_URL: Optional. Base URL for the API.
//   - ZAI_TIMEOUT: Optional. Request timeout in seconds (default: 120).
//   - ZAI_MAX_RETRIES: Optional. Maximum retry attempts (default: 3).
//   - ZAI_DISABLE_TOKEN_CACHE: Optional. Disable JWT token caching (default: true).
//   - ZAI_SOURCE_CHANNEL: Optional. Source channel identifier (default: "go-sdk").
func LoadConfig() (*Config, error) {
	apiKey := os.Getenv("ZAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ZAI_API_KEY environment variable is required")
	}

	config := &Config{
		APIKey:            apiKey,
		BaseURL:           os.Getenv("ZAI_BASE_URL"),
		Timeout:           DefaultTimeout,
		MaxRetries:        DefaultMaxRetries,
		DisableTokenCache: DefaultDisableTokenCache,
		SourceChannel:     DefaultSourceChannel,
	}

	// Parse timeout from environment
	if timeoutStr := os.Getenv("ZAI_TIMEOUT"); timeoutStr != "" {
		timeoutSec, err := strconv.Atoi(timeoutStr)
		if err != nil {
			return nil, fmt.Errorf("invalid ZAI_TIMEOUT value: %w", err)
		}
		if timeoutSec <= 0 {
			return nil, fmt.Errorf("ZAI_TIMEOUT must be positive, got: %d", timeoutSec)
		}
		config.Timeout = time.Duration(timeoutSec) * time.Second
	}

	// Parse max retries from environment
	if retriesStr := os.Getenv("ZAI_MAX_RETRIES"); retriesStr != "" {
		retries, err := strconv.Atoi(retriesStr)
		if err != nil {
			return nil, fmt.Errorf("invalid ZAI_MAX_RETRIES value: %w", err)
		}
		if retries < 0 {
			return nil, fmt.Errorf("ZAI_MAX_RETRIES must be non-negative, got: %d", retries)
		}
		config.MaxRetries = retries
	}

	// Parse disable token cache from environment
	if cacheStr := os.Getenv("ZAI_DISABLE_TOKEN_CACHE"); cacheStr != "" {
		disableCache, err := strconv.ParseBool(cacheStr)
		if err != nil {
			return nil, fmt.Errorf("invalid ZAI_DISABLE_TOKEN_CACHE value: %w", err)
		}
		config.DisableTokenCache = disableCache
	}

	// Parse source channel from environment
	if channel := os.Getenv("ZAI_SOURCE_CHANNEL"); channel != "" {
		config.SourceChannel = channel
	}

	return config, nil
}

// LoadConfigOrDefault loads configuration from environment variables.
// If the API key is missing, it returns a config with default values
// (but note: the API key is still required for actual API calls).
//
// This is useful for testing or scenarios where you want to provide
// the API key later via functional options.
func LoadConfigOrDefault() *Config {
	config, err := LoadConfig()
	if err != nil {
		// Return config with defaults but no API key
		return &Config{
			APIKey:            "",
			BaseURL:           os.Getenv("ZAI_BASE_URL"),
			Timeout:           DefaultTimeout,
			MaxRetries:        DefaultMaxRetries,
			DisableTokenCache: DefaultDisableTokenCache,
			SourceChannel:     DefaultSourceChannel,
		}
	}
	return config
}

// Validate validates the configuration.
// It returns an error if required fields are missing or invalid.
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("API key is required")
	}

	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive, got: %v", c.Timeout)
	}

	if c.MaxRetries < 0 {
		return fmt.Errorf("max retries must be non-negative, got: %d", c.MaxRetries)
	}

	return nil
}

// Clone creates a copy of the configuration.
func (c *Config) Clone() *Config {
	if c == nil {
		return nil
	}
	return &Config{
		APIKey:            c.APIKey,
		BaseURL:           c.BaseURL,
		Timeout:           c.Timeout,
		MaxRetries:        c.MaxRetries,
		DisableTokenCache: c.DisableTokenCache,
		SourceChannel:     c.SourceChannel,
	}
}
