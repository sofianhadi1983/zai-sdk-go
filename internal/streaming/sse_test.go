package streaming

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvent_IsEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		event    *Event
		expected bool
	}{
		{
			name:     "empty event",
			event:    &Event{},
			expected: true,
		},
		{
			name:     "event with data",
			event:    &Event{Data: "test"},
			expected: false,
		},
		{
			name:     "event with type",
			event:    &Event{Type: "message"},
			expected: false,
		},
		{
			name:     "event with ID",
			event:    &Event{ID: "123"},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, tt.event.IsEmpty())
		})
	}
}

func TestEvent_IsDone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		event    *Event
		expected bool
	}{
		{
			name:     "done sentinel",
			event:    &Event{Data: "[DONE]"},
			expected: true,
		},
		{
			name:     "done sentinel with whitespace",
			event:    &Event{Data: " [DONE] "},
			expected: true,
		},
		{
			name:     "regular data",
			event:    &Event{Data: "test data"},
			expected: false,
		},
		{
			name:     "empty event",
			event:    &Event{},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, tt.event.IsDone())
		})
	}
}

func TestSSEParser_Next_SingleEvent(t *testing.T) {
	t.Parallel()

	data := `data: {"message": "hello"}

`

	parser := NewSSEParser(strings.NewReader(data))

	event, err := parser.Next()
	require.NoError(t, err)
	require.NotNil(t, event)

	assert.Equal(t, `{"message": "hello"}`, event.Data)
	assert.False(t, event.IsDone())
}

func TestSSEParser_Next_MultipleEvents(t *testing.T) {
	t.Parallel()

	data := `data: first event

data: second event

data: third event

`

	parser := NewSSEParser(strings.NewReader(data))

	// First event
	event, err := parser.Next()
	require.NoError(t, err)
	assert.Equal(t, "first event", event.Data)

	// Second event
	event, err = parser.Next()
	require.NoError(t, err)
	assert.Equal(t, "second event", event.Data)

	// Third event
	event, err = parser.Next()
	require.NoError(t, err)
	assert.Equal(t, "third event", event.Data)

	// EOF
	event, err = parser.Next()
	assert.ErrorIs(t, err, io.EOF)
	assert.Nil(t, event)
}

func TestSSEParser_Next_DoneSentinel(t *testing.T) {
	t.Parallel()

	data := `data: {"message": "hello"}

data: [DONE]

`

	parser := NewSSEParser(strings.NewReader(data))

	// First event
	event, err := parser.Next()
	require.NoError(t, err)
	assert.Equal(t, `{"message": "hello"}`, event.Data)

	// Done sentinel
	event, err = parser.Next()
	assert.ErrorIs(t, err, ErrStreamDone)
	require.NotNil(t, event)
	assert.True(t, event.IsDone())
}

func TestSSEParser_Next_WithEventType(t *testing.T) {
	t.Parallel()

	data := `event: message
data: test data

`

	parser := NewSSEParser(strings.NewReader(data))

	event, err := parser.Next()
	require.NoError(t, err)
	assert.Equal(t, "message", event.Type)
	assert.Equal(t, "test data", event.Data)
}

func TestSSEParser_Next_WithID(t *testing.T) {
	t.Parallel()

	data := `id: 123
data: test data

`

	parser := NewSSEParser(strings.NewReader(data))

	event, err := parser.Next()
	require.NoError(t, err)
	assert.Equal(t, "123", event.ID)
	assert.Equal(t, "test data", event.Data)
}

func TestSSEParser_Next_MultilineData(t *testing.T) {
	t.Parallel()

	data := `data: line 1
data: line 2
data: line 3

`

	parser := NewSSEParser(strings.NewReader(data))

	event, err := parser.Next()
	require.NoError(t, err)
	assert.Equal(t, "line 1\nline 2\nline 3", event.Data)
}

func TestSSEParser_Next_Comments(t *testing.T) {
	t.Parallel()

	data := `:this is a comment
data: test data
:another comment

`

	parser := NewSSEParser(strings.NewReader(data))

	event, err := parser.Next()
	require.NoError(t, err)
	assert.Equal(t, "test data", event.Data)
}

func TestSSEParser_Next_EmptyLines(t *testing.T) {
	t.Parallel()

	data := `

data: test data


`

	parser := NewSSEParser(strings.NewReader(data))

	event, err := parser.Next()
	require.NoError(t, err)
	assert.Equal(t, "test data", event.Data)
}

func TestSSELineParser_ReadLine(t *testing.T) {
	t.Parallel()

	data := "line 1\nline 2\nline 3\n"
	parser := NewSSELineParser(strings.NewReader(data))

	line, err := parser.ReadLine()
	require.NoError(t, err)
	assert.Equal(t, "line 1", line)

	line, err = parser.ReadLine()
	require.NoError(t, err)
	assert.Equal(t, "line 2", line)

	line, err = parser.ReadLine()
	require.NoError(t, err)
	assert.Equal(t, "line 3", line)

	_, err = parser.ReadLine()
	assert.ErrorIs(t, err, io.EOF)
}

func TestSSELineParser_ReadEvent(t *testing.T) {
	t.Parallel()

	data := `data: line 1
data: line 2

data: next event

`

	parser := NewSSELineParser(strings.NewReader(data))

	// First event
	lines, err := parser.ReadEvent()
	require.NoError(t, err)
	assert.Len(t, lines, 2)
	assert.Equal(t, "data: line 1", lines[0])
	assert.Equal(t, "data: line 2", lines[1])

	// Second event
	lines, err = parser.ReadEvent()
	require.NoError(t, err)
	assert.Len(t, lines, 1)
	assert.Equal(t, "data: next event", lines[0])

	// EOF
	_, err = parser.ReadEvent()
	assert.ErrorIs(t, err, io.EOF)
}

func TestParseDataField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		line      string
		wantData  string
		wantFound bool
	}{
		{
			name:      "valid data field",
			line:      "data: test data",
			wantData:  "test data",
			wantFound: true,
		},
		{
			name:      "data field without space",
			line:      "data:test",
			wantData:  "test",
			wantFound: true,
		},
		{
			name:      "not a data field",
			line:      "event: message",
			wantData:  "",
			wantFound: false,
		},
		{
			name:      "empty data field",
			line:      "data: ",
			wantData:  "",
			wantFound: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data, found := ParseDataField(tt.line)
			assert.Equal(t, tt.wantFound, found)
			assert.Equal(t, tt.wantData, data)
		})
	}
}

func TestParseEventLines(t *testing.T) {
	t.Parallel()

	lines := []string{
		"event: message",
		"id: 123",
		"data: line 1",
		"data: line 2",
	}

	event := ParseEventLines(lines)

	assert.Equal(t, "message", event.Type)
	assert.Equal(t, "123", event.ID)
	assert.Equal(t, "line 1\nline 2", event.Data)
}

func TestIsSSEData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{
			name:     "SSE data",
			data:     []byte("data: test"),
			expected: true,
		},
		{
			name:     "not SSE data",
			data:     []byte("event: test"),
			expected: false,
		},
		{
			name:     "empty",
			data:     []byte(""),
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, IsSSEData(tt.data))
		})
	}
}

func TestExtractSSEData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "with prefix and space",
			input:    []byte("data: test data"),
			expected: []byte("test data"),
		},
		{
			name:     "with prefix no space",
			input:    []byte("data:test"),
			expected: []byte("test"),
		},
		{
			name:     "without prefix",
			input:    []byte("test data"),
			expected: []byte("test data"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := ExtractSSEData(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSSEParser_RealWorldExample(t *testing.T) {
	t.Parallel()

	data := `data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"glm-4","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null}]}

data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"glm-4","choices":[{"index":0,"delta":{"content":" world"},"finish_reason":null}]}

data: {"id":"chatcmpl-123","object":"chat.completion.chunk","created":1677652288,"model":"glm-4","choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}

data: [DONE]

`

	parser := NewSSEParser(strings.NewReader(data))

	// First chunk
	event, err := parser.Next()
	require.NoError(t, err)
	assert.Contains(t, event.Data, `"content":"Hello"`)

	// Second chunk
	event, err = parser.Next()
	require.NoError(t, err)
	assert.Contains(t, event.Data, `"content":" world"`)

	// Final chunk
	event, err = parser.Next()
	require.NoError(t, err)
	assert.Contains(t, event.Data, `"finish_reason":"stop"`)

	// Done
	event, err = parser.Next()
	assert.ErrorIs(t, err, ErrStreamDone)
	assert.True(t, event.IsDone())
}
