package ocr

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOCRRequest(t *testing.T) {
	t.Parallel()

	file := strings.NewReader("test image data")
	fileName := "test.jpg"

	req := NewOCRRequest(file, fileName, ToolTypeHandWrite)

	assert.NotNil(t, req)
	assert.Equal(t, file, req.File)
	assert.Equal(t, fileName, req.FileName)
	assert.Equal(t, ToolTypeHandWrite, req.ToolType)
	assert.Empty(t, req.LanguageType)
	assert.False(t, req.Probability)
}

func TestOCRRequest_SetMethods(t *testing.T) {
	t.Parallel()

	file := strings.NewReader("test data")
	req := NewOCRRequest(file, "test.jpg", ToolTypeHandWrite).
		SetLanguageType("zh-CN").
		SetProbability(true)

	assert.Equal(t, "zh-CN", req.LanguageType)
	assert.True(t, req.Probability)
}

func TestToolType_Constants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, ToolType("hand_write"), ToolTypeHandWrite)
}

func TestLocation_JSON(t *testing.T) {
	t.Parallel()

	loc := Location{
		Left:   10,
		Top:    20,
		Width:  100,
		Height: 50,
	}

	data, err := json.Marshal(loc)
	require.NoError(t, err)

	var decoded Location
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, loc.Left, decoded.Left)
	assert.Equal(t, loc.Top, decoded.Top)
	assert.Equal(t, loc.Width, decoded.Width)
	assert.Equal(t, loc.Height, decoded.Height)
}

func TestProbability_JSON(t *testing.T) {
	t.Parallel()

	prob := Probability{
		Average:  0.95,
		Variance: 0.02,
		Min:      0.90,
	}

	data, err := json.Marshal(prob)
	require.NoError(t, err)

	var decoded Probability
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, prob.Average, decoded.Average)
	assert.Equal(t, prob.Variance, decoded.Variance)
	assert.Equal(t, prob.Min, decoded.Min)
}

func TestWordsResult_JSON(t *testing.T) {
	t.Parallel()

	t.Run("with probability", func(t *testing.T) {
		result := WordsResult{
			Location: Location{
				Left:   10,
				Top:    20,
				Width:  100,
				Height: 50,
			},
			Words: "Hello World",
			Probability: &Probability{
				Average:  0.95,
				Variance: 0.02,
				Min:      0.90,
			},
		}

		data, err := json.Marshal(result)
		require.NoError(t, err)

		var decoded WordsResult
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		assert.Equal(t, result.Location, decoded.Location)
		assert.Equal(t, result.Words, decoded.Words)
		require.NotNil(t, decoded.Probability)
		assert.Equal(t, result.Probability.Average, decoded.Probability.Average)
	})

	t.Run("without probability", func(t *testing.T) {
		result := WordsResult{
			Location: Location{
				Left:   10,
				Top:    20,
				Width:  100,
				Height: 50,
			},
			Words: "Test",
		}

		data, err := json.Marshal(result)
		require.NoError(t, err)

		var decoded WordsResult
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		assert.Equal(t, result.Location, decoded.Location)
		assert.Equal(t, result.Words, decoded.Words)
		assert.Nil(t, decoded.Probability)
	})
}

func TestOCRResponse_JSON(t *testing.T) {
	t.Parallel()

	resp := OCRResponse{
		TaskID:         "task_123",
		Message:        "success",
		Status:         "completed",
		WordsResultNum: 2,
		WordsResult: []WordsResult{
			{
				Location: Location{Left: 10, Top: 20, Width: 100, Height: 50},
				Words:    "Hello",
				Probability: &Probability{
					Average:  0.95,
					Variance: 0.02,
					Min:      0.90,
				},
			},
			{
				Location: Location{Left: 120, Top: 20, Width: 100, Height: 50},
				Words:    "World",
				Probability: &Probability{
					Average:  0.98,
					Variance: 0.01,
					Min:      0.95,
				},
			},
		},
	}

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var decoded OCRResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, resp.TaskID, decoded.TaskID)
	assert.Equal(t, resp.Message, decoded.Message)
	assert.Equal(t, resp.Status, decoded.Status)
	assert.Equal(t, resp.WordsResultNum, decoded.WordsResultNum)
	assert.Len(t, decoded.WordsResult, 2)
	assert.Equal(t, "Hello", decoded.WordsResult[0].Words)
	assert.Equal(t, "World", decoded.WordsResult[1].Words)
}

func TestOCRResponse_GetResults(t *testing.T) {
	t.Parallel()

	t.Run("with results", func(t *testing.T) {
		resp := OCRResponse{
			WordsResultNum: 2,
			WordsResult: []WordsResult{
				{Words: "Hello"},
				{Words: "World"},
			},
		}

		results := resp.GetResults()
		assert.Len(t, results, 2)
		assert.Equal(t, "Hello", results[0].Words)
		assert.Equal(t, "World", results[1].Words)
	})

	t.Run("with nil results", func(t *testing.T) {
		resp := OCRResponse{
			WordsResultNum: 0,
			WordsResult:    nil,
		}

		results := resp.GetResults()
		assert.NotNil(t, results)
		assert.Len(t, results, 0)
	})

	t.Run("with empty results", func(t *testing.T) {
		resp := OCRResponse{
			WordsResultNum: 0,
			WordsResult:    []WordsResult{},
		}

		results := resp.GetResults()
		assert.NotNil(t, results)
		assert.Len(t, results, 0)
	})
}

func TestOCRResponse_HasResults(t *testing.T) {
	t.Parallel()

	t.Run("with results", func(t *testing.T) {
		resp := OCRResponse{
			WordsResultNum: 2,
			WordsResult: []WordsResult{
				{Words: "Hello"},
				{Words: "World"},
			},
		}

		assert.True(t, resp.HasResults())
	})

	t.Run("without results", func(t *testing.T) {
		resp := OCRResponse{
			WordsResultNum: 0,
			WordsResult:    []WordsResult{},
		}

		assert.False(t, resp.HasResults())
	})

	t.Run("with nil results", func(t *testing.T) {
		resp := OCRResponse{
			WordsResultNum: 0,
			WordsResult:    nil,
		}

		assert.False(t, resp.HasResults())
	})

	t.Run("mismatch count and results", func(t *testing.T) {
		resp := OCRResponse{
			WordsResultNum: 2,
			WordsResult:    []WordsResult{},
		}

		assert.False(t, resp.HasResults())
	})
}

func TestOCRResponse_GetText(t *testing.T) {
	t.Parallel()

	t.Run("with multiple results", func(t *testing.T) {
		resp := OCRResponse{
			WordsResultNum: 3,
			WordsResult: []WordsResult{
				{Words: "Hello"},
				{Words: "World"},
				{Words: "Test"},
			},
		}

		text := resp.GetText()
		assert.Equal(t, "Hello World Test", text)
	})

	t.Run("with single result", func(t *testing.T) {
		resp := OCRResponse{
			WordsResultNum: 1,
			WordsResult: []WordsResult{
				{Words: "Hello"},
			},
		}

		text := resp.GetText()
		assert.Equal(t, "Hello", text)
	})

	t.Run("without results", func(t *testing.T) {
		resp := OCRResponse{
			WordsResultNum: 0,
			WordsResult:    []WordsResult{},
		}

		text := resp.GetText()
		assert.Equal(t, "", text)
	})

	t.Run("with nil results", func(t *testing.T) {
		resp := OCRResponse{
			WordsResultNum: 0,
			WordsResult:    nil,
		}

		text := resp.GetText()
		assert.Equal(t, "", text)
	})
}
