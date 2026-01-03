package chat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatCompletionRequest_Setters(t *testing.T) {
	t.Parallel()

	t.Run("SetTemperature", func(t *testing.T) {
		t.Parallel()

		req := &ChatCompletionRequest{}
		req.SetTemperature(0.8)

		require.NotNil(t, req.Temperature)
		assert.Equal(t, 0.8, *req.Temperature)
	})

	t.Run("SetTopP", func(t *testing.T) {
		t.Parallel()

		req := &ChatCompletionRequest{}
		req.SetTopP(0.9)

		require.NotNil(t, req.TopP)
		assert.Equal(t, 0.9, *req.TopP)
	})

	t.Run("SetStream", func(t *testing.T) {
		t.Parallel()

		req := &ChatCompletionRequest{}
		req.SetStream(true)

		require.NotNil(t, req.Stream)
		assert.True(t, *req.Stream)
	})

	t.Run("SetMaxTokens", func(t *testing.T) {
		t.Parallel()

		req := &ChatCompletionRequest{}
		req.SetMaxTokens(100)

		require.NotNil(t, req.MaxTokens)
		assert.Equal(t, 100, *req.MaxTokens)
	})

	t.Run("chained setters", func(t *testing.T) {
		t.Parallel()

		req := &ChatCompletionRequest{
			Model: "glm-4",
		}

		req.SetTemperature(0.7).
			SetMaxTokens(500).
			SetStream(false)

		require.NotNil(t, req.Temperature)
		require.NotNil(t, req.MaxTokens)
		require.NotNil(t, req.Stream)

		assert.Equal(t, 0.7, *req.Temperature)
		assert.Equal(t, 500, *req.MaxTokens)
		assert.False(t, *req.Stream)
	})
}

func TestChatCompletionRequest_Messages(t *testing.T) {
	t.Parallel()

	t.Run("AddMessage", func(t *testing.T) {
		t.Parallel()

		req := &ChatCompletionRequest{}
		msg := NewUserMessage("Hello")

		req.AddMessage(msg)

		require.Len(t, req.Messages, 1)
		assert.Equal(t, RoleUser, req.Messages[0].Role)
		assert.Equal(t, "Hello", req.Messages[0].Content)
	})

	t.Run("AddUserMessage", func(t *testing.T) {
		t.Parallel()

		req := &ChatCompletionRequest{}
		req.AddUserMessage("Hello")

		require.Len(t, req.Messages, 1)
		assert.Equal(t, RoleUser, req.Messages[0].Role)
		assert.Equal(t, "Hello", req.Messages[0].Content)
	})

	t.Run("AddSystemMessage", func(t *testing.T) {
		t.Parallel()

		req := &ChatCompletionRequest{}
		req.AddSystemMessage("You are helpful")

		require.Len(t, req.Messages, 1)
		assert.Equal(t, RoleSystem, req.Messages[0].Role)
		assert.Equal(t, "You are helpful", req.Messages[0].Content)
	})

	t.Run("AddAssistantMessage", func(t *testing.T) {
		t.Parallel()

		req := &ChatCompletionRequest{}
		req.AddAssistantMessage("Hi there")

		require.Len(t, req.Messages, 1)
		assert.Equal(t, RoleAssistant, req.Messages[0].Role)
		assert.Equal(t, "Hi there", req.Messages[0].Content)
	})

	t.Run("chained message additions", func(t *testing.T) {
		t.Parallel()

		req := &ChatCompletionRequest{
			Model: "glm-4",
		}

		req.AddSystemMessage("You are helpful").
			AddUserMessage("Hello").
			AddAssistantMessage("Hi")

		require.Len(t, req.Messages, 3)
		assert.Equal(t, RoleSystem, req.Messages[0].Role)
		assert.Equal(t, RoleUser, req.Messages[1].Role)
		assert.Equal(t, RoleAssistant, req.Messages[2].Role)
	})
}

func TestChatCompletionRequest_Tools(t *testing.T) {
	t.Parallel()

	t.Run("AddTool", func(t *testing.T) {
		t.Parallel()

		req := &ChatCompletionRequest{}
		tool := NewFunctionTool("test", "Test function", nil)

		req.AddTool(tool)

		require.Len(t, req.Tools, 1)
		assert.Equal(t, "function", req.Tools[0].Type)
		assert.Equal(t, "test", req.Tools[0].Function.Name)
	})

	t.Run("SetToolChoice", func(t *testing.T) {
		t.Parallel()

		req := &ChatCompletionRequest{}
		req.SetToolChoice(ToolChoiceAuto)

		assert.Equal(t, "auto", req.ToolChoice)
	})

	t.Run("chained tool operations", func(t *testing.T) {
		t.Parallel()

		req := &ChatCompletionRequest{}

		req.AddTool(NewFunctionTool("func1", "First", nil)).
			AddTool(NewFunctionTool("func2", "Second", nil)).
			SetToolChoice(ToolChoiceRequired)

		require.Len(t, req.Tools, 2)
		assert.Equal(t, "required", req.ToolChoice)
	})
}

func TestChatCompletionRequest_SetResponseFormat(t *testing.T) {
	t.Parallel()

	t.Run("set text format", func(t *testing.T) {
		t.Parallel()

		req := &ChatCompletionRequest{}
		req.SetResponseFormat(ResponseFormatText)

		require.NotNil(t, req.ResponseFormat)
		assert.Equal(t, "text", req.ResponseFormat.Type)
	})

	t.Run("set JSON format", func(t *testing.T) {
		t.Parallel()

		req := &ChatCompletionRequest{}
		req.SetResponseFormat(ResponseFormatJSON)

		require.NotNil(t, req.ResponseFormat)
		assert.Equal(t, "json_object", req.ResponseFormat.Type)
	})
}

func TestChatCompletionRequest_JSON(t *testing.T) {
	t.Parallel()

	t.Run("marshal basic request", func(t *testing.T) {
		t.Parallel()

		temp := 0.7
		req := &ChatCompletionRequest{
			Model:       "glm-4",
			Messages:    []Message{NewUserMessage("Hello")},
			Temperature: &temp,
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var decoded ChatCompletionRequest
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		assert.Equal(t, "glm-4", decoded.Model)
		assert.Len(t, decoded.Messages, 1)
		require.NotNil(t, decoded.Temperature)
		assert.Equal(t, 0.7, *decoded.Temperature)
	})

	t.Run("marshal request with tools", func(t *testing.T) {
		t.Parallel()

		req := &ChatCompletionRequest{
			Model:    "glm-4",
			Messages: []Message{NewUserMessage("Get weather")},
			Tools: []Tool{
				NewFunctionTool("get_weather", "Get weather", map[string]interface{}{
					"type": "object",
				}),
			},
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)
		assert.Contains(t, string(data), "get_weather")
		assert.Contains(t, string(data), "Get weather")
	})

	t.Run("omit empty fields", func(t *testing.T) {
		t.Parallel()

		req := &ChatCompletionRequest{
			Model:    "glm-4",
			Messages: []Message{NewUserMessage("Hello")},
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		// Temperature should be omitted when nil
		assert.NotContains(t, string(data), "temperature")
		assert.NotContains(t, string(data), "tools")
	})
}

func TestChatCompletionRequest_CompleteExample(t *testing.T) {
	t.Parallel()

	// Build a complete request using chained methods
	req := &ChatCompletionRequest{
		Model: "glm-4.7",
	}

	req.AddSystemMessage("You are a helpful weather assistant").
		AddUserMessage("What's the weather in San Francisco?").
		SetTemperature(0.7).
		SetMaxTokens(500).
		AddTool(NewFunctionTool(
			"get_weather",
			"Get the current weather in a location",
			map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"location": map[string]string{
						"type":        "string",
						"description": "The city and state",
					},
				},
				"required": []string{"location"},
			},
		)).
		SetToolChoice(ToolChoiceAuto)

	// Verify the request is complete
	assert.Equal(t, "glm-4.7", req.Model)
	assert.Len(t, req.Messages, 2)
	assert.Len(t, req.Tools, 1)
	require.NotNil(t, req.Temperature)
	assert.Equal(t, 0.7, *req.Temperature)
	require.NotNil(t, req.MaxTokens)
	assert.Equal(t, 500, *req.MaxTokens)
	assert.Equal(t, "auto", req.ToolChoice)

	// Ensure it can be marshaled
	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Contains(t, string(data), "glm-4.7")
	assert.Contains(t, string(data), "get_weather")
}
