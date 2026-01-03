package streaming

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"sync"
)

var (
	// ErrStreamClosed is returned when attempting to read from a closed stream.
	ErrStreamClosed = errors.New("stream closed")

	// ErrStreamNotStarted is returned when attempting operations before streaming starts.
	ErrStreamNotStarted = errors.New("stream not started")
)

// Stream represents a generic streaming response reader.
type Stream[T any] struct {
	parser *SSEParser
	reader io.ReadCloser

	// Current event and error
	mu      sync.RWMutex
	current *T
	err     error

	// State
	done   chan struct{}
	closed bool

	// Context for cancellation
	ctx context.Context

	// Unmarshal function for custom parsing
	unmarshal func([]byte) (*T, error)
}

// StreamConfig holds configuration for creating a stream.
type StreamConfig[T any] struct {
	// Reader is the underlying stream reader.
	Reader io.ReadCloser

	// Context for cancellation (optional).
	Context context.Context

	// Unmarshal is a custom unmarshaling function.
	// If nil, uses json.Unmarshal.
	Unmarshal func([]byte) (*T, error)
}

// NewStream creates a new typed stream reader.
func NewStream[T any](config StreamConfig[T]) *Stream[T] {
	if config.Context == nil {
		config.Context = context.Background()
	}

	if config.Unmarshal == nil {
		config.Unmarshal = func(data []byte) (*T, error) {
			var result T
			err := json.Unmarshal(data, &result)
			return &result, err
		}
	}

	return &Stream[T]{
		parser:    NewSSEParser(config.Reader),
		reader:    config.Reader,
		done:      make(chan struct{}),
		ctx:       config.Context,
		unmarshal: config.Unmarshal,
	}
}

// Next advances to the next event in the stream.
// Returns false when the stream is complete or encounters an error.
func (s *Stream[T]) Next() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		s.err = ErrStreamClosed
		return false
	}

	// Check context cancellation
	select {
	case <-s.ctx.Done():
		s.err = s.ctx.Err()
		s.closeInternal()
		return false
	default:
	}

	// Read next event
	event, err := s.parser.Next()
	if err != nil {
		if errors.Is(err, ErrStreamDone) {
			// Stream completed successfully
			s.closeInternal()
			return false
		}

		if errors.Is(err, io.EOF) {
			// End of stream
			s.closeInternal()
			return false
		}

		// Error occurred
		s.err = err
		s.closeInternal()
		return false
	}

	// Skip empty events
	if event.IsEmpty() {
		return s.Next() // Recursively read next event
	}

	// Check for done sentinel
	if event.IsDone() {
		s.closeInternal()
		return false
	}

	// Parse event data
	parsed, err := s.unmarshal([]byte(event.Data))
	if err != nil {
		s.err = err
		return true // Return true to allow Err() to be called
	}

	s.current = parsed
	return true
}

// Current returns the current event data.
// Should be called after Next() returns true.
func (s *Stream[T]) Current() *T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.current
}

// Err returns any error that occurred during streaming.
func (s *Stream[T]) Err() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.err
}

// Close closes the stream and releases resources.
func (s *Stream[T]) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.closeInternal()
}

// closeInternal closes the stream without locking (must be called with lock held).
func (s *Stream[T]) closeInternal() error {
	if s.closed {
		return nil
	}

	s.closed = true
	close(s.done)

	if s.reader != nil {
		return s.reader.Close()
	}

	return nil
}

// Done returns a channel that is closed when the stream completes.
func (s *Stream[T]) Done() <-chan struct{} {
	return s.done
}

// IsClosed returns true if the stream has been closed.
func (s *Stream[T]) IsClosed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.closed
}

// Recv receives the next item from the stream.
// This is a convenience method that combines Next() and Current().
// Returns io.EOF when the stream is complete.
func (s *Stream[T]) Recv() (*T, error) {
	if !s.Next() {
		if s.Err() != nil {
			return nil, s.Err()
		}
		return nil, io.EOF
	}

	return s.Current(), nil
}

// All reads all remaining items from the stream.
// Returns a slice of all items and any error that occurred.
func (s *Stream[T]) All() ([]*T, error) {
	var items []*T

	for s.Next() {
		items = append(items, s.Current())
	}

	if err := s.Err(); err != nil {
		return items, err
	}

	return items, nil
}

// Chan returns a channel that yields stream items.
// The channel is closed when the stream completes.
// Errors are discarded; use Err() to check for errors after reading from the channel.
func (s *Stream[T]) Chan() <-chan *T {
	ch := make(chan *T)

	go func() {
		defer close(ch)
		defer s.Close()

		for s.Next() {
			select {
			case ch <- s.Current():
			case <-s.ctx.Done():
				return
			}
		}
	}()

	return ch
}

// RawStream provides low-level access to SSE events without automatic JSON parsing.
type RawStream struct {
	parser *SSEParser
	reader io.ReadCloser

	mu      sync.RWMutex
	current *Event
	err     error

	done   chan struct{}
	closed bool
	ctx    context.Context
}

// NewRawStream creates a new raw event stream.
func NewRawStream(reader io.ReadCloser, ctx context.Context) *RawStream {
	if ctx == nil {
		ctx = context.Background()
	}

	return &RawStream{
		parser: NewSSEParser(reader),
		reader: reader,
		done:   make(chan struct{}),
		ctx:    ctx,
	}
}

// Next advances to the next event.
func (s *RawStream) Next() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		s.err = ErrStreamClosed
		return false
	}

	select {
	case <-s.ctx.Done():
		s.err = s.ctx.Err()
		s.closeInternal()
		return false
	default:
	}

	event, err := s.parser.Next()
	if err != nil {
		if errors.Is(err, ErrStreamDone) || errors.Is(err, io.EOF) {
			s.closeInternal()
			return false
		}

		s.err = err
		s.closeInternal()
		return false
	}

	if event.IsDone() {
		s.closeInternal()
		return false
	}

	s.current = event
	return true
}

// Current returns the current event.
func (s *RawStream) Current() *Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.current
}

// Err returns any error.
func (s *RawStream) Err() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.err
}

// Close closes the stream.
func (s *RawStream) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.closeInternal()
}

func (s *RawStream) closeInternal() error {
	if s.closed {
		return nil
	}

	s.closed = true
	close(s.done)

	if s.reader != nil {
		return s.reader.Close()
	}

	return nil
}

// Done returns the done channel.
func (s *RawStream) Done() <-chan struct{} {
	return s.done
}
