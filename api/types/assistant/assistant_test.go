package assistant

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConversationRequest(t *testing.T) {
	t.Parallel()

	messages := []ConversationMessage{
		{
			Role: "user",
			Content: []MessageContent{
				MessageTextContent{Type: "text", Text: "Hello"},
			},
		},
	}

	req := NewConversationRequest("asst_123", messages)

	assert.Equal(t, "asst_123", req.AssistantID)
	assert.Equal(t, messages, req.Messages)
}

func TestConversationRequest_Setters(t *testing.T) {
	t.Parallel()

	t.Run("SetModel", func(t *testing.T) {
		t.Parallel()

		req := &ConversationRequest{}
		req.SetModel("GLM-4-Assistant")

		assert.Equal(t, "GLM-4-Assistant", req.Model)
	})

	t.Run("SetStream", func(t *testing.T) {
		t.Parallel()

		req := &ConversationRequest{}
		req.SetStream(true)

		assert.True(t, req.Stream)
	})

	t.Run("SetConversationID", func(t *testing.T) {
		t.Parallel()

		req := &ConversationRequest{}
		req.SetConversationID("conv_456")

		assert.Equal(t, "conv_456", req.ConversationID)
	})

	t.Run("SetAttachments", func(t *testing.T) {
		t.Parallel()

		attachments := []AssistantAttachment{
			{FileID: "file_123"},
		}

		req := &ConversationRequest{}
		req.SetAttachments(attachments)

		assert.Equal(t, attachments, req.Attachments)
	})

	t.Run("SetMetadata", func(t *testing.T) {
		t.Parallel()

		metadata := map[string]interface{}{
			"key": "value",
		}

		req := &ConversationRequest{}
		req.SetMetadata(metadata)

		assert.Equal(t, metadata, req.Metadata)
	})

	t.Run("SetRequestID", func(t *testing.T) {
		t.Parallel()

		req := &ConversationRequest{}
		req.SetRequestID("req_789")

		assert.Equal(t, "req_789", req.RequestID)
	})

	t.Run("SetUserID", func(t *testing.T) {
		t.Parallel()

		req := &ConversationRequest{}
		req.SetUserID("user_101")

		assert.Equal(t, "user_101", req.UserID)
	})

	t.Run("SetExtraParameters", func(t *testing.T) {
		t.Parallel()

		params := &ExtraParameters{
			Translate: &TranslateParameters{
				FromLanguage: "en",
				ToLanguage:   "zh",
			},
		}

		req := &ConversationRequest{}
		req.SetExtraParameters(params)

		assert.Equal(t, params, req.ExtraParameters)
	})

	t.Run("chained setters", func(t *testing.T) {
		t.Parallel()

		messages := []ConversationMessage{
			{Role: "user", Content: []MessageContent{MessageTextContent{Type: "text", Text: "Hi"}}},
		}

		req := NewConversationRequest("asst_123", messages)
		req.SetModel("GLM-4-Assistant").
			SetStream(true).
			SetConversationID("conv_456").
			SetUserID("user_101")

		assert.Equal(t, "GLM-4-Assistant", req.Model)
		assert.True(t, req.Stream)
		assert.Equal(t, "conv_456", req.ConversationID)
		assert.Equal(t, "user_101", req.UserID)
	})
}

func TestAssistantCompletion_GetText(t *testing.T) {
	t.Parallel()

	t.Run("with text content", func(t *testing.T) {
		t.Parallel()

		completion := &AssistantCompletion{
			Choices: []AssistantChoice{
				{
					Index: 0,
					Delta: TextContentBlock{
						Content: "Hello, how can I help you?",
						Role:    "assistant",
						Type:    "content",
					},
				},
			},
		}

		assert.Equal(t, "Hello, how can I help you?", completion.GetText())
	})

	t.Run("with no choices", func(t *testing.T) {
		t.Parallel()

		completion := &AssistantCompletion{
			Choices: []AssistantChoice{},
		}

		assert.Empty(t, completion.GetText())
	})

	t.Run("with tools content", func(t *testing.T) {
		t.Parallel()

		completion := &AssistantCompletion{
			Choices: []AssistantChoice{
				{
					Index: 0,
					Delta: ToolsDeltaBlock{
						Type:       "tools",
						ToolCallID: "call_123",
					},
				},
			},
		}

		assert.Empty(t, completion.GetText())
	})
}

func TestAssistantCompletion_StatusChecks(t *testing.T) {
	t.Parallel()

	t.Run("IsCompleted", func(t *testing.T) {
		t.Parallel()

		completion := &AssistantCompletion{Status: "completed"}
		assert.True(t, completion.IsCompleted())

		completion = &AssistantCompletion{Status: "in_progress"}
		assert.False(t, completion.IsCompleted())
	})

	t.Run("IsInProgress", func(t *testing.T) {
		t.Parallel()

		completion := &AssistantCompletion{Status: "in_progress"}
		assert.True(t, completion.IsInProgress())

		completion = &AssistantCompletion{Status: "completed"}
		assert.False(t, completion.IsInProgress())
	})

	t.Run("IsFailed", func(t *testing.T) {
		t.Parallel()

		completion := &AssistantCompletion{Status: "failed"}
		assert.True(t, completion.IsFailed())

		completion = &AssistantCompletion{Status: "completed"}
		assert.False(t, completion.IsFailed())
	})
}

func TestAssistantCompletion_GetError(t *testing.T) {
	t.Parallel()

	t.Run("with error", func(t *testing.T) {
		t.Parallel()

		completion := &AssistantCompletion{
			LastError: &ErrorInfo{
				Code:    "rate_limit",
				Message: "Rate limit exceeded",
			},
		}

		assert.Equal(t, "Rate limit exceeded", completion.GetError())
	})

	t.Run("without error", func(t *testing.T) {
		t.Parallel()

		completion := &AssistantCompletion{}
		assert.Empty(t, completion.GetError())
	})
}

func TestAssistantSupportResponse_GetAssistants(t *testing.T) {
	t.Parallel()

	assistants := []AssistantSupport{
		{
			AssistantID: "asst_1",
			Name:        "Assistant 1",
		},
		{
			AssistantID: "asst_2",
			Name:        "Assistant 2",
		},
	}

	resp := &AssistantSupportResponse{
		Data: assistants,
	}

	result := resp.GetAssistants()
	assert.Equal(t, assistants, result)
	assert.Len(t, result, 2)
}

func TestConversationUsageResponse_Methods(t *testing.T) {
	t.Parallel()

	conversations := []ConversationUsage{
		{
			ID:          "conv_1",
			AssistantID: "asst_1",
		},
		{
			ID:          "conv_2",
			AssistantID: "asst_1",
		},
	}

	resp := &ConversationUsageResponse{
		Data: ConversationUsageList{
			AssistantID:      "asst_1",
			HasMore:          true,
			ConversationList: conversations,
		},
	}

	assert.Equal(t, conversations, resp.GetConversations())
	assert.True(t, resp.HasMore())
}

func TestAssistantCompletion_JSON(t *testing.T) {
	t.Parallel()

	t.Run("unmarshal completion response", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "req_123",
			"conversation_id": "conv_456",
			"assistant_id": "asst_789",
			"created": 1609459200,
			"status": "completed",
			"choices": [
				{
					"index": 0,
					"delta": {
						"content": "Hello!",
						"role": "assistant",
						"type": "content"
					},
					"finish_reason": "stop",
					"metadata": {}
				}
			],
			"usage": {
				"prompt_tokens": 10,
				"completion_tokens": 5,
				"total_tokens": 15
			}
		}`

		var completion AssistantCompletion
		err := json.Unmarshal([]byte(jsonData), &completion)
		require.NoError(t, err)

		assert.Equal(t, "req_123", completion.ID)
		assert.Equal(t, "conv_456", completion.ConversationID)
		assert.Equal(t, "asst_789", completion.AssistantID)
		assert.Equal(t, int64(1609459200), completion.Created)
		assert.Equal(t, "completed", completion.Status)
		assert.Len(t, completion.Choices, 1)
		require.NotNil(t, completion.Usage)
		assert.Equal(t, 10, completion.Usage.PromptTokens)
		assert.Equal(t, 5, completion.Usage.CompletionTokens)
		assert.Equal(t, 15, completion.Usage.TotalTokens)
	})

	t.Run("unmarshal with error", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "req_123",
			"conversation_id": "conv_456",
			"assistant_id": "asst_789",
			"created": 1609459200,
			"status": "failed",
			"last_error": {
				"code": "content_filter",
				"message": "Content filtered by policy"
			},
			"choices": []
		}`

		var completion AssistantCompletion
		err := json.Unmarshal([]byte(jsonData), &completion)
		require.NoError(t, err)

		assert.True(t, completion.IsFailed())
		require.NotNil(t, completion.LastError)
		assert.Equal(t, "content_filter", completion.LastError.Code)
		assert.Equal(t, "Content filtered by policy", completion.LastError.Message)
	})

	t.Run("marshal request", func(t *testing.T) {
		t.Parallel()

		messages := []ConversationMessage{
			{
				Role: "user",
				Content: []MessageContent{
					MessageTextContent{Type: "text", Text: "Hello"},
				},
			},
		}

		req := NewConversationRequest("asst_123", messages)
		req.SetModel("GLM-4-Assistant").SetStream(true)

		data, err := json.Marshal(req)
		require.NoError(t, err)

		assert.Contains(t, string(data), "asst_123")
		assert.Contains(t, string(data), "GLM-4-Assistant")
		assert.Contains(t, string(data), "true")
	})
}

func TestAssistantSupportResponse_JSON(t *testing.T) {
	t.Parallel()

	jsonData := `{
		"code": 200,
		"msg": "success",
		"data": [
			{
				"assistant_id": "asst_1",
				"created_at": 1609459200,
				"updated_at": 1609545600,
				"name": "Code Assistant",
				"avatar": "https://example.com/avatar.png",
				"description": "Helps with coding tasks",
				"status": "publish",
				"tools": ["code_interpreter", "web_search"],
				"starter_prompts": ["Help me debug this code", "Explain this algorithm"]
			}
		]
	}`

	var resp AssistantSupportResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "success", resp.Message)
	assert.Len(t, resp.Data, 1)

	assistant := resp.Data[0]
	assert.Equal(t, "asst_1", assistant.AssistantID)
	assert.Equal(t, "Code Assistant", assistant.Name)
	assert.Equal(t, "Helps with coding tasks", assistant.Description)
	assert.Equal(t, "publish", assistant.Status)
	assert.Len(t, assistant.Tools, 2)
	assert.Len(t, assistant.StarterPrompts, 2)
}

func TestConversationUsageResponse_JSON(t *testing.T) {
	t.Parallel()

	jsonData := `{
		"code": 200,
		"msg": "success",
		"data": {
			"assistant_id": "asst_1",
			"has_more": true,
			"conversation_list": [
				{
					"id": "conv_1",
					"assistant_id": "asst_1",
					"create_time": 1609459200,
					"update_time": 1609545600,
					"usage": {
						"prompt_tokens": 100,
						"completion_tokens": 50,
						"total_tokens": 150
					}
				},
				{
					"id": "conv_2",
					"assistant_id": "asst_1",
					"create_time": 1609632000,
					"update_time": 1609718400,
					"usage": {
						"prompt_tokens": 200,
						"completion_tokens": 100,
						"total_tokens": 300
					}
				}
			]
		}
	}`

	var resp ConversationUsageResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	require.NoError(t, err)

	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "success", resp.Message)
	assert.True(t, resp.HasMore())

	conversations := resp.GetConversations()
	require.Len(t, conversations, 2)

	conv1 := conversations[0]
	assert.Equal(t, "conv_1", conv1.ID)
	assert.Equal(t, "asst_1", conv1.AssistantID)
	assert.Equal(t, 100, conv1.Usage.PromptTokens)
	assert.Equal(t, 50, conv1.Usage.CompletionTokens)
	assert.Equal(t, 150, conv1.Usage.TotalTokens)
}

func TestMessageTextContent(t *testing.T) {
	t.Parallel()

	content := MessageTextContent{
		Type: "text",
		Text: "Hello, world!",
	}

	assert.Equal(t, "text", content.Type)
	assert.Equal(t, "Hello, world!", content.Text)

	// Verify it implements MessageContent interface
	var _ MessageContent = content
}

func TestTextContentBlock(t *testing.T) {
	t.Parallel()

	block := TextContentBlock{
		Content: "Assistant response",
		Role:    "assistant",
		Type:    "content",
	}

	assert.Equal(t, "Assistant response", block.Content)
	assert.Equal(t, "assistant", block.Role)
	assert.Equal(t, "content", block.Type)

	// Verify it implements MessageContent interface
	var _ MessageContent = block
}

func TestToolsDeltaBlock(t *testing.T) {
	t.Parallel()

	block := ToolsDeltaBlock{
		Type:       "tools",
		ToolCallID: "call_123",
		ToolName:   "code_interpreter",
		ToolOutput: "Result: 42",
	}

	assert.Equal(t, "tools", block.Type)
	assert.Equal(t, "call_123", block.ToolCallID)
	assert.Equal(t, "code_interpreter", block.ToolName)
	assert.Equal(t, "Result: 42", block.ToolOutput)

	// Verify it implements MessageContent interface
	var _ MessageContent = block
}

func TestTranslateParameters(t *testing.T) {
	t.Parallel()

	params := TranslateParameters{
		FromLanguage: "en",
		ToLanguage:   "zh",
	}

	assert.Equal(t, "en", params.FromLanguage)
	assert.Equal(t, "zh", params.ToLanguage)
}

func TestExtraParameters(t *testing.T) {
	t.Parallel()

	params := ExtraParameters{
		Translate: &TranslateParameters{
			FromLanguage: "en",
			ToLanguage:   "es",
		},
	}

	require.NotNil(t, params.Translate)
	assert.Equal(t, "en", params.Translate.FromLanguage)
	assert.Equal(t, "es", params.Translate.ToLanguage)
}
