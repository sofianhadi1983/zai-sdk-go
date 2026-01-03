package chat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatCompletionResponse_GetFirstChoice(t *testing.T) {
	t.Parallel()

	t.Run("with choices", func(t *testing.T) {
		t.Parallel()

		resp := &ChatCompletionResponse{
			Choices: []Choice{
				{Index: 0, Message: NewAssistantMessage("First")},
				{Index: 1, Message: NewAssistantMessage("Second")},
			},
		}

		choice := resp.GetFirstChoice()
		require.NotNil(t, choice)
		assert.Equal(t, 0, choice.Index)
	})

	t.Run("without choices", func(t *testing.T) {
		t.Parallel()

		resp := &ChatCompletionResponse{
			Choices: []Choice{},
		}

		choice := resp.GetFirstChoice()
		assert.Nil(t, choice)
	})
}

func TestChatCompletionResponse_GetContent(t *testing.T) {
	t.Parallel()

	t.Run("with string content", func(t *testing.T) {
		t.Parallel()

		resp := &ChatCompletionResponse{
			Choices: []Choice{
				{
					Message: Message{
						Role:    RoleAssistant,
						Content: "Hello, world!",
					},
				},
			},
		}

		content := resp.GetContent()
		assert.Equal(t, "Hello, world!", content)
	})

	t.Run("without choices", func(t *testing.T) {
		t.Parallel()

		resp := &ChatCompletionResponse{
			Choices: []Choice{},
		}

		content := resp.GetContent()
		assert.Equal(t, "", content)
	})

	t.Run("with non-string content", func(t *testing.T) {
		t.Parallel()

		resp := &ChatCompletionResponse{
			Choices: []Choice{
				{
					Message: Message{
						Role: RoleAssistant,
						Content: []ContentPart{
							NewTextContentPart("Hello"),
						},
					},
				},
			},
		}

		content := resp.GetContent()
		assert.Equal(t, "", content)
	})
}

func TestChatCompletionResponse_JSON(t *testing.T) {
	t.Parallel()

	t.Run("unmarshal complete response", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "chatcmpl-123",
			"object": "chat.completion",
			"created": 1677652288,
			"model": "glm-4",
			"choices": [
				{
					"index": 0,
					"message": {
						"role": "assistant",
						"content": "Hello there!"
					},
					"finish_reason": "stop"
				}
			],
			"usage": {
				"prompt_tokens": 10,
				"completion_tokens": 20,
				"total_tokens": 30
			}
		}`

		var resp ChatCompletionResponse
		err := json.Unmarshal([]byte(jsonData), &resp)
		require.NoError(t, err)

		assert.Equal(t, "chatcmpl-123", resp.ID)
		assert.Equal(t, "chat.completion", resp.Object)
		assert.Equal(t, int64(1677652288), resp.Created)
		assert.Equal(t, "glm-4", resp.Model)
		assert.Len(t, resp.Choices, 1)

		choice := resp.Choices[0]
		assert.Equal(t, 0, choice.Index)
		assert.Equal(t, "stop", choice.FinishReason)
		assert.Equal(t, RoleAssistant, choice.Message.Role)
		assert.Equal(t, "Hello there!", choice.Message.Content)

		require.NotNil(t, resp.Usage)
		assert.Equal(t, 10, resp.Usage.PromptTokens)
		assert.Equal(t, 20, resp.Usage.CompletionTokens)
		assert.Equal(t, 30, resp.Usage.TotalTokens)
	})

	t.Run("unmarshal response with tool calls", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "chatcmpl-123",
			"object": "chat.completion",
			"created": 1677652288,
			"model": "glm-4",
			"choices": [
				{
					"index": 0,
					"message": {
						"role": "assistant",
						"content": null,
						"tool_calls": [
							{
								"id": "call-123",
								"type": "function",
								"function": {
									"name": "get_weather",
									"arguments": "{\"location\":\"SF\"}"
								}
							}
						]
					},
					"finish_reason": "tool_calls"
				}
			]
		}`

		var resp ChatCompletionResponse
		err := json.Unmarshal([]byte(jsonData), &resp)
		require.NoError(t, err)

		assert.Len(t, resp.Choices, 1)
		choice := resp.Choices[0]
		assert.Equal(t, "tool_calls", choice.FinishReason)
		require.Len(t, choice.Message.ToolCalls, 1)

		toolCall := choice.Message.ToolCalls[0]
		assert.Equal(t, "call-123", toolCall.ID)
		assert.Equal(t, "function", toolCall.Type)
		assert.Equal(t, "get_weather", toolCall.Function.Name)
		assert.Equal(t, `{"location":"SF"}`, toolCall.Function.Arguments)
	})
}

func TestChatCompletionChunk_GetContent(t *testing.T) {
	t.Parallel()

	t.Run("with content", func(t *testing.T) {
		t.Parallel()

		chunk := &ChatCompletionChunk{
			Choices: []ChunkChoice{
				{
					Delta: Delta{
						Content: "Hello",
					},
				},
			},
		}

		content := chunk.GetContent()
		assert.Equal(t, "Hello", content)
	})

	t.Run("without choices", func(t *testing.T) {
		t.Parallel()

		chunk := &ChatCompletionChunk{
			Choices: []ChunkChoice{},
		}

		content := chunk.GetContent()
		assert.Equal(t, "", content)
	})
}

func TestChatCompletionChunk_IsFinished(t *testing.T) {
	t.Parallel()

	t.Run("finished chunk", func(t *testing.T) {
		t.Parallel()

		chunk := &ChatCompletionChunk{
			Choices: []ChunkChoice{
				{
					Delta:        Delta{Content: ""},
					FinishReason: "stop",
				},
			},
		}

		assert.True(t, chunk.IsFinished())
	})

	t.Run("unfinished chunk", func(t *testing.T) {
		t.Parallel()

		chunk := &ChatCompletionChunk{
			Choices: []ChunkChoice{
				{
					Delta:        Delta{Content: "Hello"},
					FinishReason: "",
				},
			},
		}

		assert.False(t, chunk.IsFinished())
	})

	t.Run("without choices", func(t *testing.T) {
		t.Parallel()

		chunk := &ChatCompletionChunk{
			Choices: []ChunkChoice{},
		}

		assert.False(t, chunk.IsFinished())
	})
}

func TestChatCompletionChunk_JSON(t *testing.T) {
	t.Parallel()

	t.Run("unmarshal first chunk with role", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "chatcmpl-123",
			"object": "chat.completion.chunk",
			"created": 1677652288,
			"model": "glm-4",
			"choices": [
				{
					"index": 0,
					"delta": {
						"role": "assistant",
						"content": ""
					}
				}
			]
		}`

		var chunk ChatCompletionChunk
		err := json.Unmarshal([]byte(jsonData), &chunk)
		require.NoError(t, err)

		assert.Equal(t, "chatcmpl-123", chunk.ID)
		assert.Equal(t, "chat.completion.chunk", chunk.Object)
		assert.Len(t, chunk.Choices, 1)

		choice := chunk.Choices[0]
		assert.Equal(t, RoleAssistant, choice.Delta.Role)
		assert.Equal(t, "", choice.Delta.Content)
	})

	t.Run("unmarshal content chunk", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "chatcmpl-123",
			"object": "chat.completion.chunk",
			"created": 1677652288,
			"model": "glm-4",
			"choices": [
				{
					"index": 0,
					"delta": {
						"content": "Hello"
					}
				}
			]
		}`

		var chunk ChatCompletionChunk
		err := json.Unmarshal([]byte(jsonData), &chunk)
		require.NoError(t, err)

		assert.Equal(t, "Hello", chunk.Choices[0].Delta.Content)
	})

	t.Run("unmarshal final chunk", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "chatcmpl-123",
			"object": "chat.completion.chunk",
			"created": 1677652288,
			"model": "glm-4",
			"choices": [
				{
					"index": 0,
					"delta": {},
					"finish_reason": "stop"
				}
			],
			"usage": {
				"prompt_tokens": 10,
				"completion_tokens": 20,
				"total_tokens": 30
			}
		}`

		var chunk ChatCompletionChunk
		err := json.Unmarshal([]byte(jsonData), &chunk)
		require.NoError(t, err)

		assert.Equal(t, "stop", chunk.Choices[0].FinishReason)
		require.NotNil(t, chunk.Usage)
		assert.Equal(t, 10, chunk.Usage.PromptTokens)
		assert.True(t, chunk.IsFinished())
	})
}

func TestChoice_FinishReasons(t *testing.T) {
	t.Parallel()

	reasons := []string{"stop", "length", "tool_calls", "content_filter", "function_call"}

	for _, reason := range reasons {
		t.Run(reason, func(t *testing.T) {
			choice := Choice{
				FinishReason: reason,
			}

			data, err := json.Marshal(choice)
			require.NoError(t, err)
			assert.Contains(t, string(data), reason)
		})
	}
}

func TestLogProbs_JSON(t *testing.T) {
	t.Parallel()

	jsonData := `{
		"content": [
			{
				"token": "Hello",
				"logprob": -0.1,
				"bytes": [72, 101, 108, 108, 111],
				"top_logprobs": [
					{
						"token": "Hello",
						"logprob": -0.1
					},
					{
						"token": "Hi",
						"logprob": -2.3
					}
				]
			}
		]
	}`

	var logProbs LogProbs
	err := json.Unmarshal([]byte(jsonData), &logProbs)
	require.NoError(t, err)

	require.Len(t, logProbs.Content, 1)
	token := logProbs.Content[0]

	assert.Equal(t, "Hello", token.Token)
	assert.Equal(t, -0.1, token.LogProb)
	assert.Equal(t, []int{72, 101, 108, 108, 111}, token.Bytes)
	require.Len(t, token.TopLogProbs, 2)
	assert.Equal(t, "Hello", token.TopLogProbs[0].Token)
	assert.Equal(t, "Hi", token.TopLogProbs[1].Token)
}

func TestChatCompletionResponse_RealWorldExample(t *testing.T) {
	t.Parallel()

	// Simulate a real API response
	jsonData := `{
		"id": "8497903547395645060",
		"object": "chat.completion",
		"created": 1700000000,
		"model": "glm-4.7",
		"choices": [
			{
				"index": 0,
				"message": {
					"role": "assistant",
					"content": "I am an AI assistant developed by Z.ai. How can I help you today?"
				},
				"finish_reason": "stop"
			}
		],
		"usage": {
			"prompt_tokens": 15,
			"completion_tokens": 18,
			"total_tokens": 33
		}
	}`

	var resp ChatCompletionResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	require.NoError(t, err)

	// Verify all fields
	assert.Equal(t, "8497903547395645060", resp.ID)
	assert.Equal(t, "glm-4.7", resp.Model)

	// Verify helper methods work
	content := resp.GetContent()
	assert.Contains(t, content, "Z.ai")

	choice := resp.GetFirstChoice()
	require.NotNil(t, choice)
	assert.Equal(t, "stop", choice.FinishReason)

	// Verify usage
	require.NotNil(t, resp.Usage)
	assert.Equal(t, 33, resp.Usage.TotalTokens)
}
