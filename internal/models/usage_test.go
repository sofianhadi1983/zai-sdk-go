package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsage_IsEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		usage    *Usage
		expected bool
	}{
		{
			name:     "nil usage",
			usage:    nil,
			expected: true,
		},
		{
			name:     "empty usage",
			usage:    &Usage{},
			expected: true,
		},
		{
			name: "non-empty usage",
			usage: &Usage{
				PromptTokens:     10,
				CompletionTokens: 20,
				TotalTokens:      30,
			},
			expected: false,
		},
		{
			name: "partial usage",
			usage: &Usage{
				PromptTokens: 10,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, tt.usage.IsEmpty())
		})
	}
}

func TestUsage_HasCachedTokens(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		usage    *Usage
		expected bool
	}{
		{
			name:     "no details",
			usage:    &Usage{},
			expected: false,
		},
		{
			name: "with cached tokens",
			usage: &Usage{
				PromptTokensDetails: &PromptTokensDetails{
					CachedTokens: 10,
				},
			},
			expected: true,
		},
		{
			name: "without cached tokens",
			usage: &Usage{
				PromptTokensDetails: &PromptTokensDetails{
					CachedTokens: 0,
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, tt.usage.HasCachedTokens())
		})
	}
}

func TestUsage_HasReasoningTokens(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		usage    *Usage
		expected bool
	}{
		{
			name:     "no details",
			usage:    &Usage{},
			expected: false,
		},
		{
			name: "with reasoning tokens",
			usage: &Usage{
				CompletionTokensDetails: &CompletionTokensDetails{
					ReasoningTokens: 15,
				},
			},
			expected: true,
		},
		{
			name: "without reasoning tokens",
			usage: &Usage{
				CompletionTokensDetails: &CompletionTokensDetails{
					ReasoningTokens: 0,
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, tt.usage.HasReasoningTokens())
		})
	}
}

func TestUsage_GetCachedTokens(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		usage    *Usage
		expected int
	}{
		{
			name:     "no details",
			usage:    &Usage{},
			expected: 0,
		},
		{
			name: "with cached tokens",
			usage: &Usage{
				PromptTokensDetails: &PromptTokensDetails{
					CachedTokens: 25,
				},
			},
			expected: 25,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, tt.usage.GetCachedTokens())
		})
	}
}

func TestUsage_GetReasoningTokens(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		usage    *Usage
		expected int
	}{
		{
			name:     "no details",
			usage:    &Usage{},
			expected: 0,
		},
		{
			name: "with reasoning tokens",
			usage: &Usage{
				CompletionTokensDetails: &CompletionTokensDetails{
					ReasoningTokens: 50,
				},
			},
			expected: 50,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, tt.usage.GetReasoningTokens())
		})
	}
}

func TestUsage_JSON(t *testing.T) {
	t.Parallel()

	t.Run("marshal and unmarshal", func(t *testing.T) {
		t.Parallel()

		usage := &Usage{
			PromptTokens:     100,
			CompletionTokens: 200,
			TotalTokens:      300,
			PromptTokensDetails: &PromptTokensDetails{
				CachedTokens: 10,
				TextTokens:   90,
			},
			CompletionTokensDetails: &CompletionTokensDetails{
				ReasoningTokens: 50,
				TextTokens:      150,
			},
		}

		// Marshal to JSON
		data, err := json.Marshal(usage)
		require.NoError(t, err)

		// Unmarshal from JSON
		var decoded Usage
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		assert.Equal(t, usage.PromptTokens, decoded.PromptTokens)
		assert.Equal(t, usage.CompletionTokens, decoded.CompletionTokens)
		assert.Equal(t, usage.TotalTokens, decoded.TotalTokens)
		assert.Equal(t, usage.PromptTokensDetails.CachedTokens, decoded.PromptTokensDetails.CachedTokens)
		assert.Equal(t, usage.CompletionTokensDetails.ReasoningTokens, decoded.CompletionTokensDetails.ReasoningTokens)
	})

	t.Run("unmarshal without optional fields", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"prompt_tokens": 50,
			"completion_tokens": 100,
			"total_tokens": 150
		}`

		var usage Usage
		err := json.Unmarshal([]byte(jsonData), &usage)
		require.NoError(t, err)

		assert.Equal(t, 50, usage.PromptTokens)
		assert.Equal(t, 100, usage.CompletionTokens)
		assert.Equal(t, 150, usage.TotalTokens)
		assert.Nil(t, usage.PromptTokensDetails)
		assert.Nil(t, usage.CompletionTokensDetails)
	})
}

func TestPromptTokensDetails_AllFields(t *testing.T) {
	t.Parallel()

	details := &PromptTokensDetails{
		CachedTokens: 10,
		AudioTokens:  5,
		TextTokens:   85,
		ImageTokens:  0,
	}

	assert.Equal(t, 10, details.CachedTokens)
	assert.Equal(t, 5, details.AudioTokens)
	assert.Equal(t, 85, details.TextTokens)
	assert.Equal(t, 0, details.ImageTokens)
}

func TestCompletionTokensDetails_AllFields(t *testing.T) {
	t.Parallel()

	details := &CompletionTokensDetails{
		ReasoningTokens: 25,
		AudioTokens:     0,
		TextTokens:      75,
	}

	assert.Equal(t, 25, details.ReasoningTokens)
	assert.Equal(t, 0, details.AudioTokens)
	assert.Equal(t, 75, details.TextTokens)
}
