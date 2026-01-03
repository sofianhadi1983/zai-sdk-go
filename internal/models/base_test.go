package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBaseModel_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	jsonData := `{"key1": "value1", "key2": 42, "key3": true}`

	var base BaseModel
	err := json.Unmarshal([]byte(jsonData), &base)

	require.NoError(t, err)
	assert.NotNil(t, base.Extra)
	assert.Len(t, base.Extra, 3)
	assert.Equal(t, "value1", base.Extra["key1"])
	assert.Equal(t, float64(42), base.Extra["key2"]) // JSON numbers unmarshal as float64
	assert.Equal(t, true, base.Extra["key3"])
}

func TestBaseModel_MarshalJSON(t *testing.T) {
	t.Parallel()

	base := BaseModel{
		Extra: map[string]interface{}{
			"name":  "test",
			"count": 10,
			"active": true,
		},
	}

	data, err := json.Marshal(base)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "test", result["name"])
	assert.Equal(t, float64(10), result["count"])
	assert.Equal(t, true, result["active"])
}

func TestBaseModel_Get(t *testing.T) {
	t.Parallel()

	base := BaseModel{
		Extra: map[string]interface{}{
			"existing": "value",
		},
	}

	// Get existing key
	val, ok := base.Get("existing")
	assert.True(t, ok)
	assert.Equal(t, "value", val)

	// Get non-existing key
	val, ok = base.Get("nonexistent")
	assert.False(t, ok)
	assert.Nil(t, val)
}

func TestBaseModel_Set(t *testing.T) {
	t.Parallel()

	base := BaseModel{}

	// Set on nil Extra
	base.Set("key1", "value1")
	assert.NotNil(t, base.Extra)
	assert.Equal(t, "value1", base.Extra["key1"])

	// Set on existing Extra
	base.Set("key2", 42)
	assert.Equal(t, 42, base.Extra["key2"])
}

func TestBaseModel_Keys(t *testing.T) {
	t.Parallel()

	t.Run("with extra fields", func(t *testing.T) {
		t.Parallel()

		base := BaseModel{
			Extra: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
		}

		keys := base.Keys()
		assert.Len(t, keys, 3)
		assert.Contains(t, keys, "key1")
		assert.Contains(t, keys, "key2")
		assert.Contains(t, keys, "key3")
	})

	t.Run("without extra fields", func(t *testing.T) {
		t.Parallel()

		base := BaseModel{}
		keys := base.Keys()
		assert.Empty(t, keys)
	})
}

func TestBaseModel_ToMap(t *testing.T) {
	t.Parallel()

	base := BaseModel{
		Extra: map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		},
	}

	result := base.ToMap()
	assert.Len(t, result, 2)
	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, 42, result["key2"])

	// Verify it's a copy (mutations don't affect original)
	result["key3"] = "value3"
	assert.NotContains(t, base.Extra, "key3")
}

func TestCommonRequestFields_Validate(t *testing.T) {
	t.Parallel()

	fields := CommonRequestFields{
		Model:     "glm-4",
		RequestID: "test-123",
	}

	err := fields.Validate()
	assert.NoError(t, err)
}

func TestCommonResponseFields_GetCreatedTime(t *testing.T) {
	t.Parallel()

	t.Run("with valid timestamp", func(t *testing.T) {
		t.Parallel()

		now := time.Now().Unix()
		fields := CommonResponseFields{
			Created: now,
		}

		createdTime := fields.GetCreatedTime()
		assert.Equal(t, now, createdTime.Unix())
	})

	t.Run("with zero timestamp", func(t *testing.T) {
		t.Parallel()

		fields := CommonResponseFields{
			Created: 0,
		}

		createdTime := fields.GetCreatedTime()
		assert.True(t, createdTime.IsZero())
	})
}

func TestCommonRequestFields_Stream(t *testing.T) {
	t.Parallel()

	streamTrue := true
	streamFalse := false

	tests := []struct {
		name     string
		stream   *bool
		expected *bool
	}{
		{
			name:     "stream enabled",
			stream:   &streamTrue,
			expected: &streamTrue,
		},
		{
			name:     "stream disabled",
			stream:   &streamFalse,
			expected: &streamFalse,
		},
		{
			name:     "stream not set",
			stream:   nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fields := CommonRequestFields{
				Stream: tt.stream,
			}

			if tt.expected == nil {
				assert.Nil(t, fields.Stream)
			} else {
				require.NotNil(t, fields.Stream)
				assert.Equal(t, *tt.expected, *fields.Stream)
			}
		})
	}
}
