package chat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessage_Constructors(t *testing.T) {
	t.Parallel()

	t.Run("NewUserMessage", func(t *testing.T) {
		t.Parallel()

		msg := NewUserMessage("Hello")
		assert.Equal(t, RoleUser, msg.Role)
		assert.Equal(t, "Hello", msg.Content)
	})

	t.Run("NewSystemMessage", func(t *testing.T) {
		t.Parallel()

		msg := NewSystemMessage("You are a helpful assistant")
		assert.Equal(t, RoleSystem, msg.Role)
		assert.Equal(t, "You are a helpful assistant", msg.Content)
	})

	t.Run("NewAssistantMessage", func(t *testing.T) {
		t.Parallel()

		msg := NewAssistantMessage("Hi there!")
		assert.Equal(t, RoleAssistant, msg.Role)
		assert.Equal(t, "Hi there!", msg.Content)
	})

	t.Run("NewToolMessage", func(t *testing.T) {
		t.Parallel()

		msg := NewToolMessage("call-123", "Result")
		assert.Equal(t, RoleTool, msg.Role)
		assert.Equal(t, "Result", msg.Content)
		assert.Equal(t, "call-123", msg.ToolCallID)
	})
}

func TestContentPart_Constructors(t *testing.T) {
	t.Parallel()

	t.Run("NewTextContentPart", func(t *testing.T) {
		t.Parallel()

		part := NewTextContentPart("Hello")
		assert.Equal(t, "text", part.Type)
		assert.Equal(t, "Hello", part.Text)
		assert.Nil(t, part.ImageURL)
	})

	t.Run("NewImageContentPart", func(t *testing.T) {
		t.Parallel()

		part := NewImageContentPart("https://example.com/image.jpg")
		assert.Equal(t, "image_url", part.Type)
		require.NotNil(t, part.ImageURL)
		assert.Equal(t, "https://example.com/image.jpg", part.ImageURL.URL)
	})
}

func TestFunctionCall_GetArguments(t *testing.T) {
	t.Parallel()

	t.Run("valid JSON", func(t *testing.T) {
		t.Parallel()

		fc := &FunctionCall{
			Name:      "get_weather",
			Arguments: `{"location":"San Francisco","unit":"celsius"}`,
		}

		var args map[string]string
		err := fc.GetArguments(&args)
		require.NoError(t, err)

		assert.Equal(t, "San Francisco", args["location"])
		assert.Equal(t, "celsius", args["unit"])
	})

	t.Run("invalid JSON", func(t *testing.T) {
		t.Parallel()

		fc := &FunctionCall{
			Name:      "test",
			Arguments: `{invalid json}`,
		}

		var args map[string]string
		err := fc.GetArguments(&args)
		assert.Error(t, err)
	})
}

func TestNewFunctionTool(t *testing.T) {
	t.Parallel()

	params := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"location": map[string]string{
				"type":        "string",
				"description": "The city and state",
			},
		},
		"required": []string{"location"},
	}

	tool := NewFunctionTool("get_weather", "Get the current weather", params)

	assert.Equal(t, "function", tool.Type)
	assert.Equal(t, "get_weather", tool.Function.Name)
	assert.Equal(t, "Get the current weather", tool.Function.Description)
	assert.Equal(t, params, tool.Function.Parameters)
}

func TestMessage_JSON(t *testing.T) {
	t.Parallel()

	t.Run("marshal simple message", func(t *testing.T) {
		t.Parallel()

		msg := Message{
			Role:    RoleUser,
			Content: "Hello",
		}

		data, err := json.Marshal(msg)
		require.NoError(t, err)

		var decoded Message
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		assert.Equal(t, RoleUser, decoded.Role)
		assert.Equal(t, "Hello", decoded.Content)
	})

	t.Run("marshal multimodal message", func(t *testing.T) {
		t.Parallel()

		msg := Message{
			Role: RoleUser,
			Content: []ContentPart{
				NewTextContentPart("What's in this image?"),
				NewImageContentPart("https://example.com/image.jpg"),
			},
		}

		data, err := json.Marshal(msg)
		require.NoError(t, err)
		assert.Contains(t, string(data), "What's in this image?")
		assert.Contains(t, string(data), "https://example.com/image.jpg")
	})

	t.Run("marshal message with tool calls", func(t *testing.T) {
		t.Parallel()

		msg := Message{
			Role:    RoleAssistant,
			Content: "",
			ToolCalls: []ToolCall{
				{
					ID:   "call-123",
					Type: "function",
					Function: FunctionCall{
						Name:      "get_weather",
						Arguments: `{"location":"SF"}`,
					},
				},
			},
		}

		data, err := json.Marshal(msg)
		require.NoError(t, err)
		assert.Contains(t, string(data), "call-123")
		assert.Contains(t, string(data), "get_weather")
	})
}

func TestRole_Values(t *testing.T) {
	t.Parallel()

	assert.Equal(t, Role("system"), RoleSystem)
	assert.Equal(t, Role("user"), RoleUser)
	assert.Equal(t, Role("assistant"), RoleAssistant)
	assert.Equal(t, Role("tool"), RoleTool)
	assert.Equal(t, Role("function"), RoleFunction)
}

func TestToolChoice_Values(t *testing.T) {
	t.Parallel()

	assert.Equal(t, ToolChoice("none"), ToolChoiceNone)
	assert.Equal(t, ToolChoice("auto"), ToolChoiceAuto)
	assert.Equal(t, ToolChoice("required"), ToolChoiceRequired)
}

func TestResponseFormat_Constants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "text", ResponseFormatText.Type)
	assert.Equal(t, "json_object", ResponseFormatJSON.Type)
}
