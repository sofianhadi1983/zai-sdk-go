package configs

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		want    *Config
		wantErr bool
	}{
		{
			name: "valid config with all values",
			envVars: map[string]string{
				"ZAI_API_KEY":              "test-api-key",
				"ZAI_BASE_URL":             "https://api.test.com",
				"ZAI_TIMEOUT":              "60",
				"ZAI_MAX_RETRIES":          "5",
				"ZAI_DISABLE_TOKEN_CACHE":  "false",
				"ZAI_SOURCE_CHANNEL":       "test-sdk",
			},
			want: &Config{
				APIKey:            "test-api-key",
				BaseURL:           "https://api.test.com",
				Timeout:           60 * time.Second,
				MaxRetries:        5,
				DisableTokenCache: false,
				SourceChannel:     "test-sdk",
			},
			wantErr: false,
		},
		{
			name: "valid config with defaults",
			envVars: map[string]string{
				"ZAI_API_KEY": "test-api-key",
			},
			want: &Config{
				APIKey:            "test-api-key",
				BaseURL:           "",
				Timeout:           DefaultTimeout,
				MaxRetries:        DefaultMaxRetries,
				DisableTokenCache: DefaultDisableTokenCache,
				SourceChannel:     DefaultSourceChannel,
			},
			wantErr: false,
		},
		{
			name:    "missing API key",
			envVars: map[string]string{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid timeout",
			envVars: map[string]string{
				"ZAI_API_KEY": "test-api-key",
				"ZAI_TIMEOUT": "invalid",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "negative timeout",
			envVars: map[string]string{
				"ZAI_API_KEY": "test-api-key",
				"ZAI_TIMEOUT": "-10",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid max retries",
			envVars: map[string]string{
				"ZAI_API_KEY":     "test-api-key",
				"ZAI_MAX_RETRIES": "invalid",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "negative max retries",
			envVars: map[string]string{
				"ZAI_API_KEY":     "test-api-key",
				"ZAI_MAX_RETRIES": "-1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "zero max retries is valid",
			envVars: map[string]string{
				"ZAI_API_KEY":     "test-api-key",
				"ZAI_MAX_RETRIES": "0",
			},
			want: &Config{
				APIKey:            "test-api-key",
				BaseURL:           "",
				Timeout:           DefaultTimeout,
				MaxRetries:        0,
				DisableTokenCache: DefaultDisableTokenCache,
				SourceChannel:     DefaultSourceChannel,
			},
			wantErr: false,
		},
		{
			name: "invalid disable token cache",
			envVars: map[string]string{
				"ZAI_API_KEY":             "test-api-key",
				"ZAI_DISABLE_TOKEN_CACHE": "invalid",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			got, err := LoadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != nil {
				if got.APIKey != tt.want.APIKey {
					t.Errorf("APIKey = %v, want %v", got.APIKey, tt.want.APIKey)
				}
				if got.BaseURL != tt.want.BaseURL {
					t.Errorf("BaseURL = %v, want %v", got.BaseURL, tt.want.BaseURL)
				}
				if got.Timeout != tt.want.Timeout {
					t.Errorf("Timeout = %v, want %v", got.Timeout, tt.want.Timeout)
				}
				if got.MaxRetries != tt.want.MaxRetries {
					t.Errorf("MaxRetries = %v, want %v", got.MaxRetries, tt.want.MaxRetries)
				}
				if got.DisableTokenCache != tt.want.DisableTokenCache {
					t.Errorf("DisableTokenCache = %v, want %v", got.DisableTokenCache, tt.want.DisableTokenCache)
				}
				if got.SourceChannel != tt.want.SourceChannel {
					t.Errorf("SourceChannel = %v, want %v", got.SourceChannel, tt.want.SourceChannel)
				}
			}
		})
	}
}

func TestLoadConfigOrDefault(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		want    *Config
	}{
		{
			name: "with API key",
			envVars: map[string]string{
				"ZAI_API_KEY": "test-api-key",
			},
			want: &Config{
				APIKey:            "test-api-key",
				BaseURL:           "",
				Timeout:           DefaultTimeout,
				MaxRetries:        DefaultMaxRetries,
				DisableTokenCache: DefaultDisableTokenCache,
				SourceChannel:     DefaultSourceChannel,
			},
		},
		{
			name:    "without API key returns defaults",
			envVars: map[string]string{},
			want: &Config{
				APIKey:            "",
				BaseURL:           "",
				Timeout:           DefaultTimeout,
				MaxRetries:        DefaultMaxRetries,
				DisableTokenCache: DefaultDisableTokenCache,
				SourceChannel:     DefaultSourceChannel,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			got := LoadConfigOrDefault()
			if got.APIKey != tt.want.APIKey {
				t.Errorf("APIKey = %v, want %v", got.APIKey, tt.want.APIKey)
			}
			if got.Timeout != tt.want.Timeout {
				t.Errorf("Timeout = %v, want %v", got.Timeout, tt.want.Timeout)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				APIKey:     "test-api-key",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			wantErr: false,
		},
		{
			name: "missing API key",
			config: &Config{
				APIKey:     "",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			wantErr: true,
		},
		{
			name: "zero timeout",
			config: &Config{
				APIKey:     "test-api-key",
				Timeout:    0,
				MaxRetries: 3,
			},
			wantErr: true,
		},
		{
			name: "negative timeout",
			config: &Config{
				APIKey:     "test-api-key",
				Timeout:    -10 * time.Second,
				MaxRetries: 3,
			},
			wantErr: true,
		},
		{
			name: "negative max retries",
			config: &Config{
				APIKey:     "test-api-key",
				Timeout:    30 * time.Second,
				MaxRetries: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_Clone(t *testing.T) {
	t.Parallel()

	original := &Config{
		APIKey:            "test-api-key",
		BaseURL:           "https://api.test.com",
		Timeout:           60 * time.Second,
		MaxRetries:        5,
		DisableTokenCache: false,
		SourceChannel:     "test-sdk",
	}

	cloned := original.Clone()

	// Verify values are the same
	if cloned.APIKey != original.APIKey {
		t.Errorf("Clone APIKey = %v, want %v", cloned.APIKey, original.APIKey)
	}
	if cloned.BaseURL != original.BaseURL {
		t.Errorf("Clone BaseURL = %v, want %v", cloned.BaseURL, original.BaseURL)
	}

	// Modify clone and verify original is unchanged
	cloned.APIKey = "modified-key"
	if original.APIKey == "modified-key" {
		t.Error("Modifying clone affected original")
	}

	// Test nil clone
	var nilConfig *Config
	clonedNil := nilConfig.Clone()
	if clonedNil != nil {
		t.Error("Clone of nil should return nil")
	}
}

func TestDefaultConstants(t *testing.T) {
	t.Parallel()

	if DefaultTimeout != 120*time.Second {
		t.Errorf("DefaultTimeout = %v, want 120s", DefaultTimeout)
	}
	if DefaultMaxRetries != 3 {
		t.Errorf("DefaultMaxRetries = %d, want 3", DefaultMaxRetries)
	}
	if DefaultDisableTokenCache != true {
		t.Errorf("DefaultDisableTokenCache = %v, want true", DefaultDisableTokenCache)
	}
	if DefaultSourceChannel != "go-sdk" {
		t.Errorf("DefaultSourceChannel = %s, want go-sdk", DefaultSourceChannel)
	}
}

// Benchmark for LoadConfig
func BenchmarkLoadConfig(b *testing.B) {
	os.Setenv("ZAI_API_KEY", "benchmark-key")
	defer os.Unsetenv("ZAI_API_KEY")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = LoadConfig()
	}
}
