package zai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/z-ai/zai-sdk-go/api/types/moderation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModerationsService_Create(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/moderations", r.URL.Path)

		// Parse request body
		var reqBody moderation.ModerationRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		// Verify request fields
		assert.Equal(t, "moderation", reqBody.Model)

		// Send mock response
		resp := moderation.ModerationResponse{
			ID:    "modr-abc123",
			Model: "moderation",
			Results: []moderation.ModerationResult{
				{
					Flagged: true,
					Categories: moderation.ModerationCategories{
						Harassment:            false,
						HarassmentThreatening: false,
						Hate:                  true,
						HateThreatening:       false,
						SelfHarm:              false,
						SelfHarmInstructions:  false,
						SelfHarmIntent:        false,
						Sexual:                false,
						SexualMinors:          false,
						Violence:              false,
						ViolenceGraphic:       false,
					},
					CategoryScores: moderation.ModerationCategoryScores{
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
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	req := moderation.NewTextModerationRequest("moderation", "test content")

	resp, err := client.Moderations.Create(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "modr-abc123", resp.ID)
	assert.Equal(t, "moderation", resp.Model)

	// Verify results
	results := resp.GetResults()
	assert.Len(t, results, 1)

	result := results[0]
	assert.True(t, result.Flagged)
	assert.False(t, result.IsSafe())

	// Verify categories
	assert.False(t, result.Categories.Harassment)
	assert.True(t, result.Categories.Hate)
	assert.False(t, result.Categories.Violence)

	// Verify scores
	assert.InDelta(t, 0.85, result.CategoryScores.Hate, 0.01)
	assert.InDelta(t, 0.1, result.CategoryScores.Harassment, 0.01)
	assert.InDelta(t, 0.15, result.CategoryScores.Violence, 0.01)

	// Check overall flagged status
	assert.True(t, resp.IsFlagged())
}

func TestModerationsService_Create_SafeContent(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/moderations", r.URL.Path)

		// Send mock response with safe content
		resp := moderation.ModerationResponse{
			ID:    "modr-safe123",
			Model: "moderation",
			Results: []moderation.ModerationResult{
				{
					Flagged: false,
					Categories: moderation.ModerationCategories{
						Harassment:            false,
						HarassmentThreatening: false,
						Hate:                  false,
						HateThreatening:       false,
						SelfHarm:              false,
						SelfHarmInstructions:  false,
						SelfHarmIntent:        false,
						Sexual:                false,
						SexualMinors:          false,
						Violence:              false,
						ViolenceGraphic:       false,
					},
					CategoryScores: moderation.ModerationCategoryScores{
						Harassment:            0.01,
						HarassmentThreatening: 0.01,
						Hate:                  0.01,
						HateThreatening:       0.01,
						SelfHarm:              0.01,
						SelfHarmInstructions:  0.01,
						SelfHarmIntent:        0.01,
						Sexual:                0.01,
						SexualMinors:          0.01,
						Violence:              0.01,
						ViolenceGraphic:       0.01,
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	req := moderation.NewTextModerationRequest("moderation", "hello world")

	resp, err := client.Moderations.Create(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	results := resp.GetResults()
	assert.Len(t, results, 1)

	result := results[0]
	assert.False(t, result.Flagged)
	assert.True(t, result.IsSafe())
	assert.False(t, resp.IsFlagged())
}

func TestModerationsService_Create_BatchInput(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/moderations", r.URL.Path)

		// Parse request body
		var reqBody moderation.ModerationRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		assert.Equal(t, "moderation", reqBody.Model)

		// Send mock response with multiple results
		resp := moderation.ModerationResponse{
			ID:    "modr-batch123",
			Model: "moderation",
			Results: []moderation.ModerationResult{
				{
					Flagged: false,
					Categories: moderation.ModerationCategories{
						Sexual: false,
					},
					CategoryScores: moderation.ModerationCategoryScores{
						Sexual: 0.05,
					},
				},
				{
					Flagged: true,
					Categories: moderation.ModerationCategories{
						Sexual: true,
					},
					CategoryScores: moderation.ModerationCategoryScores{
						Sexual: 0.92,
					},
				},
				{
					Flagged: false,
					Categories: moderation.ModerationCategories{
						Violence: false,
					},
					CategoryScores: moderation.ModerationCategoryScores{
						Violence: 0.08,
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	texts := []string{"safe text", "unsafe text", "another safe text"}
	req := moderation.NewBatchTextModerationRequest("moderation", texts)

	resp, err := client.Moderations.Create(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	results := resp.GetResults()
	assert.Len(t, results, 3)

	assert.False(t, results[0].Flagged)
	assert.True(t, results[1].Flagged)
	assert.False(t, results[2].Flagged)

	// Overall flagged status should be true because one result is flagged
	assert.True(t, resp.IsFlagged())
}

func TestModerationsService_CheckText(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/moderations", r.URL.Path)

		// Parse request body
		var reqBody moderation.ModerationRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		// Verify text moderation request structure
		assert.Equal(t, "moderation", reqBody.Model)
		inputMap, ok := reqBody.Input.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "text", inputMap["type"])
		assert.Equal(t, "test content", inputMap["text"])

		// Send mock response
		resp := moderation.ModerationResponse{
			ID:    "modr-text123",
			Model: "moderation",
			Results: []moderation.ModerationResult{
				{
					Flagged: false,
					Categories: moderation.ModerationCategories{
						Harassment: false,
					},
					CategoryScores: moderation.ModerationCategoryScores{
						Harassment: 0.02,
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	resp, err := client.Moderations.CheckText(context.Background(), "moderation", "test content")
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "modr-text123", resp.ID)
	assert.False(t, resp.IsFlagged())
}

func TestModerationsService_CheckBatch(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/moderations", r.URL.Path)

		// Parse request body
		var reqBody moderation.ModerationRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		// Verify batch request structure
		assert.Equal(t, "moderation", reqBody.Model)
		texts, ok := reqBody.Input.([]interface{})
		require.True(t, ok)
		assert.Len(t, texts, 2)

		// Send mock response
		resp := moderation.ModerationResponse{
			ID:    "modr-batch123",
			Model: "moderation",
			Results: []moderation.ModerationResult{
				{Flagged: false},
				{Flagged: false},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	texts := []string{"text1", "text2"}
	resp, err := client.Moderations.CheckBatch(context.Background(), "moderation", texts)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "modr-batch123", resp.ID)
	assert.Len(t, resp.GetResults(), 2)
	assert.False(t, resp.IsFlagged())
}

func TestModerationsService_Create_APIError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{
				"message": "Invalid moderation request",
				"type":    "invalid_request_error",
				"code":    "invalid_moderation_request",
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(
		WithAPIKey("test-key.test-secret"),
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	req := moderation.NewModerationRequest("moderation", nil)

	_, err = client.Moderations.Create(context.Background(), req)
	require.Error(t, err)
}
