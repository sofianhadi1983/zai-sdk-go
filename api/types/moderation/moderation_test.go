package moderation

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewModerationRequest(t *testing.T) {
	t.Parallel()

	t.Run("new request with string input", func(t *testing.T) {
		t.Parallel()

		req := NewModerationRequest("moderation", "test content")

		assert.Equal(t, "moderation", req.Model)
		assert.Equal(t, "test content", req.Input)
	})

	t.Run("new request with structured input", func(t *testing.T) {
		t.Parallel()

		input := map[string]interface{}{
			"type": "text",
			"text": "test content",
		}
		req := NewModerationRequest("moderation", input)

		assert.Equal(t, "moderation", req.Model)
		assert.Equal(t, input, req.Input)
	})
}

func TestNewTextModerationRequest(t *testing.T) {
	t.Parallel()

	req := NewTextModerationRequest("moderation", "hello world")

	assert.Equal(t, "moderation", req.Model)

	inputMap, ok := req.Input.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "text", inputMap["type"])
	assert.Equal(t, "hello world", inputMap["text"])
}

func TestNewBatchTextModerationRequest(t *testing.T) {
	t.Parallel()

	texts := []string{"text1", "text2", "text3"}
	req := NewBatchTextModerationRequest("moderation", texts)

	assert.Equal(t, "moderation", req.Model)
	assert.Equal(t, texts, req.Input)
}

func TestModerationRequest_JSONMarshaling(t *testing.T) {
	t.Parallel()

	t.Run("marshal and unmarshal string input", func(t *testing.T) {
		t.Parallel()

		req := NewModerationRequest("moderation", "test content")

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var decoded ModerationRequest
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		assert.Equal(t, req.Model, decoded.Model)
		assert.Equal(t, req.Input, decoded.Input)
	})

	t.Run("marshal and unmarshal array input", func(t *testing.T) {
		t.Parallel()

		req := NewBatchTextModerationRequest("moderation", []string{"text1", "text2"})

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var decoded ModerationRequest
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		assert.Equal(t, req.Model, decoded.Model)
	})
}

func TestModerationResponse(t *testing.T) {
	t.Parallel()

	t.Run("unmarshal from API response", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "modr-abc123",
			"model": "moderation",
			"results": [
				{
					"flagged": true,
					"categories": {
						"harassment": false,
						"harassment/threatening": false,
						"hate": true,
						"hate/threatening": false,
						"self-harm": false,
						"self-harm/instructions": false,
						"self-harm/intent": false,
						"sexual": false,
						"sexual/minors": false,
						"violence": false,
						"violence/graphic": false
					},
					"category_scores": {
						"harassment": 0.1,
						"harassment/threatening": 0.05,
						"hate": 0.85,
						"hate/threatening": 0.03,
						"self-harm": 0.01,
						"self-harm/instructions": 0.02,
						"self-harm/intent": 0.01,
						"sexual": 0.04,
						"sexual/minors": 0.01,
						"violence": 0.15,
						"violence/graphic": 0.08
					}
				}
			]
		}`

		var resp ModerationResponse
		err := json.Unmarshal([]byte(jsonData), &resp)
		require.NoError(t, err)

		assert.Equal(t, "modr-abc123", resp.ID)
		assert.Equal(t, "moderation", resp.Model)
		assert.Len(t, resp.Results, 1)

		result := resp.Results[0]
		assert.True(t, result.Flagged)
		assert.False(t, result.Categories.Harassment)
		assert.True(t, result.Categories.Hate)
		assert.False(t, result.Categories.Violence)

		assert.InDelta(t, 0.85, result.CategoryScores.Hate, 0.01)
		assert.InDelta(t, 0.1, result.CategoryScores.Harassment, 0.01)
		assert.InDelta(t, 0.15, result.CategoryScores.Violence, 0.01)
	})

	t.Run("multiple results", func(t *testing.T) {
		t.Parallel()

		jsonData := `{
			"id": "modr-xyz789",
			"model": "moderation",
			"results": [
				{
					"flagged": false,
					"categories": {
						"harassment": false,
						"harassment/threatening": false,
						"hate": false,
						"hate/threatening": false,
						"self-harm": false,
						"self-harm/instructions": false,
						"self-harm/intent": false,
						"sexual": false,
						"sexual/minors": false,
						"violence": false,
						"violence/graphic": false
					},
					"category_scores": {
						"harassment": 0.01,
						"harassment/threatening": 0.01,
						"hate": 0.01,
						"hate/threatening": 0.01,
						"self-harm": 0.01,
						"self-harm/instructions": 0.01,
						"self-harm/intent": 0.01,
						"sexual": 0.01,
						"sexual/minors": 0.01,
						"violence": 0.01,
						"violence/graphic": 0.01
					}
				},
				{
					"flagged": true,
					"categories": {
						"harassment": false,
						"harassment/threatening": false,
						"hate": false,
						"hate/threatening": false,
						"self-harm": false,
						"self-harm/instructions": false,
						"self-harm/intent": false,
						"sexual": true,
						"sexual/minors": false,
						"violence": false,
						"violence/graphic": false
					},
					"category_scores": {
						"harassment": 0.05,
						"harassment/threatening": 0.02,
						"hate": 0.03,
						"hate/threatening": 0.01,
						"self-harm": 0.01,
						"self-harm/instructions": 0.01,
						"self-harm/intent": 0.01,
						"sexual": 0.92,
						"sexual/minors": 0.05,
						"violence": 0.04,
						"violence/graphic": 0.02
					}
				}
			]
		}`

		var resp ModerationResponse
		err := json.Unmarshal([]byte(jsonData), &resp)
		require.NoError(t, err)

		assert.Len(t, resp.Results, 2)
		assert.False(t, resp.Results[0].Flagged)
		assert.True(t, resp.Results[1].Flagged)
		assert.True(t, resp.Results[1].Categories.Sexual)
	})
}

func TestModerationResponse_GetResults(t *testing.T) {
	t.Parallel()

	t.Run("with results", func(t *testing.T) {
		t.Parallel()

		resp := &ModerationResponse{
			Results: []ModerationResult{
				{Flagged: false},
				{Flagged: true},
			},
		}

		results := resp.GetResults()
		assert.Len(t, results, 2)
	})

	t.Run("empty results", func(t *testing.T) {
		t.Parallel()

		resp := &ModerationResponse{}
		results := resp.GetResults()
		assert.Empty(t, results)
	})
}

func TestModerationResponse_IsFlagged(t *testing.T) {
	t.Parallel()

	t.Run("with flagged result", func(t *testing.T) {
		t.Parallel()

		resp := &ModerationResponse{
			Results: []ModerationResult{
				{Flagged: false},
				{Flagged: true},
			},
		}

		assert.True(t, resp.IsFlagged())
	})

	t.Run("without flagged results", func(t *testing.T) {
		t.Parallel()

		resp := &ModerationResponse{
			Results: []ModerationResult{
				{Flagged: false},
				{Flagged: false},
			},
		}

		assert.False(t, resp.IsFlagged())
	})

	t.Run("empty results", func(t *testing.T) {
		t.Parallel()

		resp := &ModerationResponse{}
		assert.False(t, resp.IsFlagged())
	})
}

func TestModerationResult_IsSafe(t *testing.T) {
	t.Parallel()

	t.Run("safe result", func(t *testing.T) {
		t.Parallel()

		result := ModerationResult{Flagged: false}
		assert.True(t, result.IsSafe())
	})

	t.Run("unsafe result", func(t *testing.T) {
		t.Parallel()

		result := ModerationResult{Flagged: true}
		assert.False(t, result.IsSafe())
	})
}

func TestModerationCategories_HasCategory(t *testing.T) {
	t.Parallel()

	categories := ModerationCategories{
		Harassment: false,
		Hate:       true,
		Sexual:     false,
		Violence:   true,
	}

	t.Run("check hate category", func(t *testing.T) {
		t.Parallel()

		hasHate := categories.HasCategory(func(c *ModerationCategories) bool {
			return c.Hate
		})
		assert.True(t, hasHate)
	})

	t.Run("check harassment category", func(t *testing.T) {
		t.Parallel()

		hasHarassment := categories.HasCategory(func(c *ModerationCategories) bool {
			return c.Harassment
		})
		assert.False(t, hasHarassment)
	})

	t.Run("check violence category", func(t *testing.T) {
		t.Parallel()

		hasViolence := categories.HasCategory(func(c *ModerationCategories) bool {
			return c.Violence
		})
		assert.True(t, hasViolence)
	})
}

func TestModerationCategoryScores(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		t.Parallel()

		scores := ModerationCategoryScores{
			Harassment:            0.1,
			HarassmentThreatening: 0.05,
			Hate:                  0.85,
			HateThreatening:       0.03,
			SelfHarm:              0.01,
			SelfHarmInstructions:  0.02,
			SelfHarmIntent:        0.01,
			Sexual:                0.04,
			SexualMinors:          0.01,
			Violence:              0.15,
			ViolenceGraphic:       0.08,
		}

		data, err := json.Marshal(scores)
		require.NoError(t, err)

		var decoded ModerationCategoryScores
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		assert.InDelta(t, scores.Hate, decoded.Hate, 0.01)
		assert.InDelta(t, scores.Violence, decoded.Violence, 0.01)
		assert.InDelta(t, scores.Sexual, decoded.Sexual, 0.01)
	})
}
