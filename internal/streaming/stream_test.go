package streaming

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func TestNewStream(t *testing.T) {
	t.Parallel()

	data := `data: {"content":"hello","role":"assistant"}

`
	reader := nopCloser{strings.NewReader(data)}

	stream := NewStream[testMessage](StreamConfig[testMessage]{
		Reader: reader,
	})

	assert.NotNil(t, stream)
	assert.NotNil(t, stream.parser)
	assert.NotNil(t, stream.done)
}

func TestStream_Next_SingleItem(t *testing.T) {
	t.Parallel()

	data := `data: {"content":"hello","role":"assistant"}

`
	reader := nopCloser{strings.NewReader(data)}

	stream := NewStream[testMessage](StreamConfig[testMessage]{
		Reader: reader,
	})
	defer stream.Close()

	// Should have one item
	hasNext := stream.Next()
	assert.True(t, hasNext)
	assert.NoError(t, stream.Err())

	msg := stream.Current()
	require.NotNil(t, msg)
	assert.Equal(t, "hello", msg.Content)
	assert.Equal(t, "assistant", msg.Role)

	// Should reach EOF
	hasNext = stream.Next()
	assert.False(t, hasNext)
}

func TestStream_Next_MultipleItems(t *testing.T) {
	t.Parallel()

	data := `data: {"content":"first","role":"user"}

data: {"content":"second","role":"assistant"}

data: {"content":"third","role":"user"}

`
	reader := nopCloser{strings.NewReader(data)}

	stream := NewStream[testMessage](StreamConfig[testMessage]{
		Reader: reader,
	})
	defer stream.Close()

	// First item
	assert.True(t, stream.Next())
	msg := stream.Current()
	assert.Equal(t, "first", msg.Content)

	// Second item
	assert.True(t, stream.Next())
	msg = stream.Current()
	assert.Equal(t, "second", msg.Content)

	// Third item
	assert.True(t, stream.Next())
	msg = stream.Current()
	assert.Equal(t, "third", msg.Content)

	// EOF
	assert.False(t, stream.Next())
}

func TestStream_Next_DoneSentinel(t *testing.T) {
	t.Parallel()

	data := `data: {"content":"hello","role":"assistant"}

data: [DONE]

`
	reader := nopCloser{strings.NewReader(data)}

	stream := NewStream[testMessage](StreamConfig[testMessage]{
		Reader: reader,
	})
	defer stream.Close()

	// First item
	assert.True(t, stream.Next())
	msg := stream.Current()
	assert.Equal(t, "hello", msg.Content)

	// Done sentinel should stop iteration
	assert.False(t, stream.Next())
	assert.NoError(t, stream.Err())
}

func TestStream_Recv(t *testing.T) {
	t.Parallel()

	data := `data: {"content":"hello","role":"assistant"}

data: {"content":"world","role":"user"}

`
	reader := nopCloser{strings.NewReader(data)}

	stream := NewStream[testMessage](StreamConfig[testMessage]{
		Reader: reader,
	})
	defer stream.Close()

	// First message
	msg, err := stream.Recv()
	require.NoError(t, err)
	assert.Equal(t, "hello", msg.Content)

	// Second message
	msg, err = stream.Recv()
	require.NoError(t, err)
	assert.Equal(t, "world", msg.Content)

	// EOF
	msg, err = stream.Recv()
	assert.ErrorIs(t, err, io.EOF)
	assert.Nil(t, msg)
}

func TestStream_All(t *testing.T) {
	t.Parallel()

	data := `data: {"content":"first","role":"user"}

data: {"content":"second","role":"assistant"}

data: {"content":"third","role":"user"}

`
	reader := nopCloser{strings.NewReader(data)}

	stream := NewStream[testMessage](StreamConfig[testMessage]{
		Reader: reader,
	})
	defer stream.Close()

	messages, err := stream.All()
	require.NoError(t, err)
	require.Len(t, messages, 3)

	assert.Equal(t, "first", messages[0].Content)
	assert.Equal(t, "second", messages[1].Content)
	assert.Equal(t, "third", messages[2].Content)
}

func TestStream_Chan(t *testing.T) {
	t.Parallel()

	data := `data: {"content":"first","role":"user"}

data: {"content":"second","role":"assistant"}

data: {"content":"third","role":"user"}

`
	reader := nopCloser{strings.NewReader(data)}

	stream := NewStream[testMessage](StreamConfig[testMessage]{
		Reader: reader,
	})

	// Read from channel
	var messages []*testMessage
	for msg := range stream.Chan() {
		messages = append(messages, msg)
	}

	require.Len(t, messages, 3)
	assert.Equal(t, "first", messages[0].Content)
	assert.Equal(t, "second", messages[1].Content)
	assert.Equal(t, "third", messages[2].Content)
}

func TestStream_ContextCancellation(t *testing.T) {
	t.Parallel()

	// Create infinite stream
	data := strings.Repeat(`data: {"content":"test","role":"user"}

`, 1000)

	reader := nopCloser{strings.NewReader(data)}

	ctx, cancel := context.WithCancel(context.Background())

	stream := NewStream[testMessage](StreamConfig[testMessage]{
		Reader:  reader,
		Context: ctx,
	})
	defer stream.Close()

	// Read first item
	assert.True(t, stream.Next())

	// Cancel context
	cancel()

	// Next should return false due to cancellation
	assert.False(t, stream.Next())
	assert.ErrorIs(t, stream.Err(), context.Canceled)
}

func TestStream_Close(t *testing.T) {
	t.Parallel()

	data := `data: {"content":"hello","role":"assistant"}

`
	reader := nopCloser{strings.NewReader(data)}

	stream := NewStream[testMessage](StreamConfig[testMessage]{
		Reader: reader,
	})

	// Close stream
	err := stream.Close()
	assert.NoError(t, err)

	// Should be closed
	assert.True(t, stream.IsClosed())

	// Next should return false
	assert.False(t, stream.Next())
	assert.ErrorIs(t, stream.Err(), ErrStreamClosed)

	// Close again should not error
	err = stream.Close()
	assert.NoError(t, err)
}

func TestStream_Done(t *testing.T) {
	t.Parallel()

	data := `data: {"content":"hello","role":"assistant"}

`
	reader := nopCloser{strings.NewReader(data)}

	stream := NewStream[testMessage](StreamConfig[testMessage]{
		Reader: reader,
	})

	// Done channel should not be closed initially
	select {
	case <-stream.Done():
		t.Error("Done channel should not be closed yet")
	case <-time.After(10 * time.Millisecond):
		// Expected
	}

	// Read all items
	for stream.Next() {
	}

	// Done channel should be closed now
	select {
	case <-stream.Done():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Error("Done channel should be closed")
	}
}

func TestStream_CustomUnmarshal(t *testing.T) {
	t.Parallel()

	data := `data: CUSTOM:hello:assistant

`
	reader := nopCloser{strings.NewReader(data)}

	customUnmarshal := func(data []byte) (*testMessage, error) {
		// Parse custom format: CUSTOM:content:role
		parts := strings.Split(string(data), ":")
		if len(parts) != 3 {
			return nil, assert.AnError
		}
		return &testMessage{
			Content: parts[1],
			Role:    parts[2],
		}, nil
	}

	stream := NewStream[testMessage](StreamConfig[testMessage]{
		Reader:    reader,
		Unmarshal: customUnmarshal,
	})
	defer stream.Close()

	assert.True(t, stream.Next())
	msg := stream.Current()
	assert.Equal(t, "hello", msg.Content)
	assert.Equal(t, "assistant", msg.Role)
}

func TestStream_InvalidJSON(t *testing.T) {
	t.Parallel()

	data := `data: {invalid json}

`
	reader := nopCloser{strings.NewReader(data)}

	stream := NewStream[testMessage](StreamConfig[testMessage]{
		Reader: reader,
	})
	defer stream.Close()

	// Next should return true to allow checking error
	hasNext := stream.Next()
	assert.True(t, hasNext)

	// Should have unmarshaling error
	err := stream.Err()
	assert.Error(t, err)

	var jsonErr *json.SyntaxError
	assert.ErrorAs(t, err, &jsonErr)
}

func TestNewRawStream(t *testing.T) {
	t.Parallel()

	data := `data: test data

`
	reader := nopCloser{strings.NewReader(data)}

	stream := NewRawStream(reader, nil)

	assert.NotNil(t, stream)
	assert.NotNil(t, stream.parser)
	assert.NotNil(t, stream.done)
}

func TestRawStream_Next(t *testing.T) {
	t.Parallel()

	data := `event: message
data: test data
id: 123

`
	reader := nopCloser{strings.NewReader(data)}

	stream := NewRawStream(reader, nil)
	defer stream.Close()

	assert.True(t, stream.Next())
	event := stream.Current()
	require.NotNil(t, event)

	assert.Equal(t, "message", event.Type)
	assert.Equal(t, "test data", event.Data)
	assert.Equal(t, "123", event.ID)

	// EOF
	assert.False(t, stream.Next())
}

func TestRawStream_ContextCancellation(t *testing.T) {
	t.Parallel()

	data := strings.Repeat(`data: test

`, 1000)

	reader := nopCloser{strings.NewReader(data)}

	ctx, cancel := context.WithCancel(context.Background())
	stream := NewRawStream(reader, ctx)
	defer stream.Close()

	// Read first event
	assert.True(t, stream.Next())

	// Cancel context
	cancel()

	// Should stop
	assert.False(t, stream.Next())
	assert.ErrorIs(t, stream.Err(), context.Canceled)
}

func TestRawStream_Close(t *testing.T) {
	t.Parallel()

	data := `data: test

`
	reader := nopCloser{strings.NewReader(data)}

	stream := NewRawStream(reader, nil)

	err := stream.Close()
	assert.NoError(t, err)

	// Next should return false
	assert.False(t, stream.Next())
	assert.ErrorIs(t, stream.Err(), ErrStreamClosed)
}
