package zai

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sofianhadi1983/zai-sdk-go/internal/constants"
	"github.com/sofianhadi1983/zai-sdk-go/internal/logger"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	t.Run("with valid API key", func(t *testing.T) {
		t.Parallel()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
		)

		require.NoError(t, err)
		require.NotNil(t, client)
		defer client.Close()

		assert.NotNil(t, client.baseClient)
		assert.NotNil(t, client.config)
		assert.Equal(t, "test-key.test-secret", client.config.APIKey)
		assert.Equal(t, constants.ZaiBaseURL, client.config.BaseURL)
	})

	t.Run("without API key", func(t *testing.T) {
		t.Parallel()

		client, err := NewClient()

		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "API key is required")
	})

	t.Run("with all options", func(t *testing.T) {
		t.Parallel()

		customLogger := logger.Default()
		customTimeout := 30 * time.Second
		customRetries := 5

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL("https://custom.api.com"),
			WithTimeout(customTimeout),
			WithMaxRetries(customRetries),
			WithDisableTokenCache(),
			WithLogger(customLogger),
		)

		require.NoError(t, err)
		require.NotNil(t, client)
		defer client.Close()

		assert.Equal(t, "test-key.test-secret", client.config.APIKey)
		assert.Equal(t, "https://custom.api.com", client.config.BaseURL)
		assert.Equal(t, customTimeout, client.config.Timeout)
		assert.Equal(t, customRetries, client.config.MaxRetries)
		assert.True(t, client.config.DisableTokenCache)
		assert.Equal(t, customLogger, client.config.Logger)
	})

	t.Run("with partial options", func(t *testing.T) {
		t.Parallel()

		client, err := NewClient(
			WithAPIKey("test-key.test-secret"),
			WithTimeout(45*time.Second),
		)

		require.NoError(t, err)
		require.NotNil(t, client)
		defer client.Close()

		assert.Equal(t, "test-key.test-secret", client.config.APIKey)
		assert.Equal(t, constants.ZaiBaseURL, client.config.BaseURL)
		assert.Equal(t, 45*time.Second, client.config.Timeout)
		assert.False(t, client.config.DisableTokenCache)
	})
}

func TestNewZhipuClient(t *testing.T) {
	t.Parallel()

	t.Run("with valid API key", func(t *testing.T) {
		t.Parallel()

		client, err := NewZhipuClient(
			WithAPIKey("test-key.test-secret"),
		)

		require.NoError(t, err)
		require.NotNil(t, client)
		defer client.Close()

		assert.Equal(t, "test-key.test-secret", client.config.APIKey)
		assert.Equal(t, constants.ZhipuBaseURL, client.config.BaseURL)
	})

	t.Run("with custom base URL override", func(t *testing.T) {
		t.Parallel()

		client, err := NewZhipuClient(
			WithAPIKey("test-key.test-secret"),
			WithBaseURL("https://custom.zhipu.com"),
		)

		require.NoError(t, err)
		require.NotNil(t, client)
		defer client.Close()

		// The last option should win
		assert.Equal(t, "https://custom.zhipu.com", client.config.BaseURL)
	})

	t.Run("without API key", func(t *testing.T) {
		t.Parallel()

		client, err := NewZhipuClient()

		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "API key is required")
	})
}

func TestNewClientFromEnv(t *testing.T) {
	t.Parallel()

	t.Run("with environment variables", func(t *testing.T) {
		// Save original values
		originalAPIKey := os.Getenv("ZAI_API_KEY")
		originalBaseURL := os.Getenv("ZAI_BASE_URL")
		defer func() {
			if originalAPIKey != "" {
				os.Setenv("ZAI_API_KEY", originalAPIKey)
			} else {
				os.Unsetenv("ZAI_API_KEY")
			}
			if originalBaseURL != "" {
				os.Setenv("ZAI_BASE_URL", originalBaseURL)
			} else {
				os.Unsetenv("ZAI_BASE_URL")
			}
		}()

		// Set test environment variables
		os.Setenv("ZAI_API_KEY", "env-key.env-secret")
		os.Setenv("ZAI_BASE_URL", "https://env.api.com")

		client, err := NewClientFromEnv()

		require.NoError(t, err)
		require.NotNil(t, client)
		defer client.Close()

		assert.Equal(t, "env-key.env-secret", client.config.APIKey)
		assert.Equal(t, "https://env.api.com", client.config.BaseURL)
	})

	t.Run("with API key only", func(t *testing.T) {
		// Save original value
		originalAPIKey := os.Getenv("ZAI_API_KEY")
		originalBaseURL := os.Getenv("ZAI_BASE_URL")
		defer func() {
			if originalAPIKey != "" {
				os.Setenv("ZAI_API_KEY", originalAPIKey)
			} else {
				os.Unsetenv("ZAI_API_KEY")
			}
			if originalBaseURL != "" {
				os.Setenv("ZAI_BASE_URL", originalBaseURL)
			} else {
				os.Unsetenv("ZAI_BASE_URL")
			}
		}()

		// Set only API key
		os.Setenv("ZAI_API_KEY", "env-key.env-secret")
		os.Unsetenv("ZAI_BASE_URL")

		client, err := NewClientFromEnv()

		require.NoError(t, err)
		require.NotNil(t, client)
		defer client.Close()

		assert.Equal(t, "env-key.env-secret", client.config.APIKey)
		assert.Equal(t, constants.ZaiBaseURL, client.config.BaseURL) // Should use default
	})

	t.Run("with additional options", func(t *testing.T) {
		// Save original value
		originalAPIKey := os.Getenv("ZAI_API_KEY")
		defer func() {
			if originalAPIKey != "" {
				os.Setenv("ZAI_API_KEY", originalAPIKey)
			} else {
				os.Unsetenv("ZAI_API_KEY")
			}
		}()

		os.Setenv("ZAI_API_KEY", "env-key.env-secret")

		client, err := NewClientFromEnv(
			WithTimeout(60*time.Second),
			WithMaxRetries(5),
		)

		require.NoError(t, err)
		require.NotNil(t, client)
		defer client.Close()

		assert.Equal(t, "env-key.env-secret", client.config.APIKey)
		assert.Equal(t, 60*time.Second, client.config.Timeout)
		assert.Equal(t, 5, client.config.MaxRetries)
	})

	t.Run("without API key", func(t *testing.T) {
		// Save original value
		originalAPIKey := os.Getenv("ZAI_API_KEY")
		defer func() {
			if originalAPIKey != "" {
				os.Setenv("ZAI_API_KEY", originalAPIKey)
			} else {
				os.Unsetenv("ZAI_API_KEY")
			}
		}()

		os.Unsetenv("ZAI_API_KEY")

		client, err := NewClientFromEnv()

		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "API key is required")
	})
}

func TestNewZhipuClientFromEnv(t *testing.T) {
	t.Parallel()

	t.Run("with environment variable", func(t *testing.T) {
		// Save original value
		originalAPIKey := os.Getenv("ZAI_API_KEY")
		defer func() {
			if originalAPIKey != "" {
				os.Setenv("ZAI_API_KEY", originalAPIKey)
			} else {
				os.Unsetenv("ZAI_API_KEY")
			}
		}()

		os.Setenv("ZAI_API_KEY", "env-key.env-secret")

		client, err := NewZhipuClientFromEnv()

		require.NoError(t, err)
		require.NotNil(t, client)
		defer client.Close()

		assert.Equal(t, "env-key.env-secret", client.config.APIKey)
		assert.Equal(t, constants.ZhipuBaseURL, client.config.BaseURL)
	})

	t.Run("with additional options", func(t *testing.T) {
		// Save original value
		originalAPIKey := os.Getenv("ZAI_API_KEY")
		defer func() {
			if originalAPIKey != "" {
				os.Setenv("ZAI_API_KEY", originalAPIKey)
			} else {
				os.Unsetenv("ZAI_API_KEY")
			}
		}()

		os.Setenv("ZAI_API_KEY", "env-key.env-secret")

		client, err := NewZhipuClientFromEnv(
			WithDisableTokenCache(),
		)

		require.NoError(t, err)
		require.NotNil(t, client)
		defer client.Close()

		assert.True(t, client.config.DisableTokenCache)
	})
}

func TestClient_GetConfig(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithTimeout(30*time.Second),
	)

	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()

	config := client.GetConfig()
	assert.Equal(t, "test-key.test-secret", config.APIKey)
	assert.Equal(t, 30*time.Second, config.Timeout)
}

func TestClient_GetLogger(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
	)

	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()

	logger := client.GetLogger()
	assert.NotNil(t, logger)
}

func TestClient_Close(t *testing.T) {
	t.Parallel()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
	)

	require.NoError(t, err)
	require.NotNil(t, client)

	// Should not panic
	client.Close()

	// Closing again should also not panic
	client.Close()
}

func TestClientOptions(t *testing.T) {
	t.Parallel()

	t.Run("WithAPIKey", func(t *testing.T) {
		t.Parallel()

		config := &ClientConfig{}
		opt := WithAPIKey("test-key")
		opt(config)

		assert.Equal(t, "test-key", config.APIKey)
	})

	t.Run("WithBaseURL", func(t *testing.T) {
		t.Parallel()

		config := &ClientConfig{}
		opt := WithBaseURL("https://test.com")
		opt(config)

		assert.Equal(t, "https://test.com", config.BaseURL)
	})

	t.Run("WithTimeout", func(t *testing.T) {
		t.Parallel()

		config := &ClientConfig{}
		opt := WithTimeout(30 * time.Second)
		opt(config)

		assert.Equal(t, 30*time.Second, config.Timeout)
	})

	t.Run("WithMaxRetries", func(t *testing.T) {
		t.Parallel()

		config := &ClientConfig{}
		opt := WithMaxRetries(5)
		opt(config)

		assert.Equal(t, 5, config.MaxRetries)
	})

	t.Run("WithDisableTokenCache", func(t *testing.T) {
		t.Parallel()

		config := &ClientConfig{}
		opt := WithDisableTokenCache()
		opt(config)

		assert.True(t, config.DisableTokenCache)
	})

	t.Run("WithLogger", func(t *testing.T) {
		t.Parallel()

		config := &ClientConfig{}
		customLogger := logger.Default()
		opt := WithLogger(customLogger)
		opt(config)

		assert.Equal(t, customLogger, config.Logger)
	})
}
