package agents

import (
	"encoding/json"
	"testing"

	"github.com/sofianhadi1983/zai-sdk-go/api/types/chat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAgentInvokeRequest(t *testing.T) {
	t.Parallel()

	messages := []chat.Message{
		{Role: chat.RoleUser, Content: "Hello agent"},
	}

	req := NewAgentInvokeRequest("agent-123", messages)

	assert.Equal(t, "agent-123", req.AgentID)
	assert.Len(t, req.Messages, 1)
}

func TestAgentInvokeRequest_BuilderPattern(t *testing.T) {
	t.Parallel()

	messages := []chat.Message{
		{Role: chat.RoleUser, Content: "test"},
	}

	customVars := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	sensitiveCheck := &SensitiveWordCheck{
		Type:   "ALL",
		Status: "ENABLE",
	}

	req := NewAgentInvokeRequest("agent-456", messages).
		SetStream(true).
		SetRequestID("req_123").
		SetUserID("user_789").
		SetCustomVariables(customVars).
		SetSensitiveWordCheck(sensitiveCheck)

	assert.True(t, req.Stream)
	assert.Equal(t, "req_123", req.RequestID)
	assert.Equal(t, "user_789", req.UserID)
	assert.Equal(t, customVars, req.CustomVariables)
	assert.NotNil(t, req.SensitiveWordCheck)
}

func TestAgentCompletionResponse(t *testing.T) {
	t.Parallel()

	jsonData := `{
		"id": "comp_abc123",
		"agent_id": "agent_xyz",
		"conversation_id": "conv_123",
		"status": "completed",
		"request_id": "req_456",
		"choices": [
			{
				"index": 0,
				"finish_reason": "stop",
				"message": {
					"role": "assistant",
					"content": "This is the agent response"
				}
			}
		],
		"usage": {
			"prompt_tokens": 10,
			"completion_tokens": 20,
			"total_tokens": 30
		}
	}`

	var resp AgentCompletionResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	require.NoError(t, err)

	assert.Equal(t, "comp_abc123", resp.ID)
	assert.Equal(t, "agent_xyz", resp.AgentID)
	assert.Equal(t, "conv_123", resp.ConversationID)
	assert.Equal(t, "completed", resp.Status)

	choices := resp.GetChoices()
	assert.Len(t, choices, 1)
	assert.Equal(t, "stop", choices[0].FinishReason)
	assert.Equal(t, "assistant", choices[0].Message.Role)

	content := resp.GetContent()
	assert.Equal(t, "This is the agent response", content)

	assert.False(t, resp.HasError())
	assert.NotNil(t, resp.Usage)
	assert.Equal(t, 30, resp.Usage.TotalTokens)
}

func TestAgentCompletionChunk(t *testing.T) {
	t.Parallel()

	jsonData := `{
		"id": "chunk_123",
		"agent_id": "agent_xyz",
		"conversation_id": "conv_456",
		"choices": [
			{
				"index": 0,
				"delta": {
					"role": "assistant",
					"content": "streaming content"
				}
			}
		]
	}`

	var chunk AgentCompletionChunk
	err := json.Unmarshal([]byte(jsonData), &chunk)
	require.NoError(t, err)

	assert.Equal(t, "chunk_123", chunk.ID)
	assert.Len(t, chunk.Choices, 1)
	assert.Equal(t, "assistant", chunk.Choices[0].Delta.Role)

	content := chunk.GetContent()
	assert.Equal(t, "streaming content", content)
	assert.False(t, chunk.HasError())
}

func TestAgentError(t *testing.T) {
	t.Parallel()

	jsonData := `{
		"id": "comp_err",
		"error": {
			"code": "invalid_request",
			"message": "Agent not found"
		},
		"choices": []
	}`

	var resp AgentCompletionResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	require.NoError(t, err)

	assert.True(t, resp.HasError())
	assert.Equal(t, "invalid_request", resp.Error.Code)
	assert.Equal(t, "Agent not found", resp.Error.Message)
}

func TestNewAgentAsyncResultRequest(t *testing.T) {
	t.Parallel()

	req := NewAgentAsyncResultRequest("agent-123").
		SetAsyncID("async-456").
		SetConversationID("conv-789").
		SetCustomVariables(map[string]interface{}{"key": "value"})

	assert.Equal(t, "agent-123", req.AgentID)
	assert.Equal(t, "async-456", req.AsyncID)
	assert.Equal(t, "conv-789", req.ConversationID)
	assert.NotNil(t, req.CustomVariables)
}
