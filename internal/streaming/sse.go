// Package streaming provides Server-Sent Events (SSE) parsing and stream handling.
package streaming

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
)

const (
	// EventFieldData is the SSE data field prefix.
	EventFieldData = "data:"

	// EventFieldEvent is the SSE event type field prefix.
	EventFieldEvent = "event:"

	// EventFieldID is the SSE event ID field prefix.
	EventFieldID = "id:"

	// EventFieldRetry is the SSE retry field prefix.
	EventFieldRetry = "retry:"

	// DoneSentinel indicates the end of a stream.
	DoneSentinel = "[DONE]"
)

var (
	// ErrStreamDone is returned when the stream has completed.
	ErrStreamDone = errors.New("stream done")

	// ErrInvalidEvent is returned when an event is malformed.
	ErrInvalidEvent = errors.New("invalid SSE event")
)

// Event represents a Server-Sent Event.
type Event struct {
	// Type is the event type (from "event:" field).
	Type string

	// Data is the event data (from "data:" field).
	Data string

	// ID is the event ID (from "id:" field).
	ID string

	// Retry is the retry timeout in milliseconds.
	Retry int

	// Raw holds the complete raw event text.
	Raw string
}

// IsEmpty returns true if the event has no data.
func (e *Event) IsEmpty() bool {
	return e.Data == "" && e.Type == "" && e.ID == ""
}

// IsDone returns true if this event is a completion marker.
func (e *Event) IsDone() bool {
	return strings.TrimSpace(e.Data) == DoneSentinel
}

// SSEParser parses Server-Sent Events from a stream.
type SSEParser struct {
	scanner *bufio.Scanner
	reader  io.Reader
}

// NewSSEParser creates a new SSE parser for the given reader.
func NewSSEParser(reader io.Reader) *SSEParser {
	scanner := bufio.NewScanner(reader)
	// Increase buffer size for large events
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024) // Max 1MB per line

	return &SSEParser{
		scanner: scanner,
		reader:  reader,
	}
}

// Next reads the next event from the stream.
// Returns ErrStreamDone when the stream completes.
// Returns io.EOF when there are no more events.
func (p *SSEParser) Next() (*Event, error) {
	event := &Event{}
	var rawLines []string
	var hasData bool

	for p.scanner.Scan() {
		line := p.scanner.Text()
		rawLines = append(rawLines, line)

		// Empty line indicates end of event
		if line == "" {
			if hasData {
				event.Raw = strings.Join(rawLines, "\n")

				// Check for done sentinel
				if event.IsDone() {
					return event, ErrStreamDone
				}

				return event, nil
			}
			// Reset for next event
			rawLines = rawLines[:0]
			continue
		}

		// Parse field
		if err := p.parseField(line, event); err != nil {
			continue // Skip invalid fields per SSE spec
		}

		// Track if we have any data
		if strings.HasPrefix(line, EventFieldData) {
			hasData = true
		}
	}

	// Check for scanner errors
	if err := p.scanner.Err(); err != nil {
		return nil, err
	}

	// If we have accumulated data, return it
	if hasData {
		event.Raw = strings.Join(rawLines, "\n")
		if event.IsDone() {
			return event, ErrStreamDone
		}
		return event, nil
	}

	return nil, io.EOF
}

// parseField parses a single SSE field line.
func (p *SSEParser) parseField(line string, event *Event) error {
	// Comments start with ":"
	if strings.HasPrefix(line, ":") {
		return nil // Ignore comments
	}

	// Split on first colon
	colonIdx := strings.Index(line, ":")
	if colonIdx == -1 {
		// Field with no value
		return nil
	}

	field := line[:colonIdx]
	value := line[colonIdx+1:]

	// Remove leading space from value (per SSE spec)
	if len(value) > 0 && value[0] == ' ' {
		value = value[1:]
	}

	switch field {
	case "data":
		// Multiple data fields are concatenated with newlines
		if event.Data != "" {
			event.Data += "\n"
		}
		event.Data += value

	case "event":
		event.Type = value

	case "id":
		event.ID = value

	case "retry":
		// Retry is milliseconds as integer
		// We don't parse it here, just store as-is
		// Consumer can parse if needed

	default:
		// Unknown fields are ignored per SSE spec
	}

	return nil
}

// SSELineParser is a simpler line-by-line SSE parser.
type SSELineParser struct {
	reader *bufio.Reader
}

// NewSSELineParser creates a new line parser.
func NewSSELineParser(reader io.Reader) *SSELineParser {
	return &SSELineParser{
		reader: bufio.NewReader(reader),
	}
}

// ReadLine reads a single line from the stream.
// Returns the line without the newline character.
func (p *SSELineParser) ReadLine() (string, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// Remove trailing newline
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")

	return line, nil
}

// ReadEvent reads a complete event (until blank line).
func (p *SSELineParser) ReadEvent() ([]string, error) {
	var lines []string

	for {
		line, err := p.ReadLine()
		if err != nil {
			if err == io.EOF && len(lines) > 0 {
				return lines, nil
			}
			return nil, err
		}

		// Empty line marks end of event
		if line == "" {
			if len(lines) > 0 {
				return lines, nil
			}
			// Skip multiple empty lines
			continue
		}

		lines = append(lines, line)
	}
}

// ParseDataField extracts the data from a "data: " prefixed line.
func ParseDataField(line string) (string, bool) {
	if !strings.HasPrefix(line, EventFieldData) {
		return "", false
	}

	data := strings.TrimPrefix(line, EventFieldData)
	// Remove leading space per SSE spec
	data = strings.TrimPrefix(data, " ")

	return data, true
}

// ParseEventLines parses lines into an Event.
func ParseEventLines(lines []string) *Event {
	event := &Event{
		Raw: strings.Join(lines, "\n"),
	}

	var dataLines []string

	for _, line := range lines {
		// Skip comments
		if strings.HasPrefix(line, ":") {
			continue
		}

		colonIdx := strings.Index(line, ":")
		if colonIdx == -1 {
			continue
		}

		field := line[:colonIdx]
		value := line[colonIdx+1:]

		// Remove leading space
		if len(value) > 0 && value[0] == ' ' {
			value = value[1:]
		}

		switch field {
		case "data":
			dataLines = append(dataLines, value)
		case "event":
			event.Type = value
		case "id":
			event.ID = value
		}
	}

	// Join multiple data lines with newlines
	if len(dataLines) > 0 {
		event.Data = strings.Join(dataLines, "\n")
	}

	return event
}

// IsSSEData checks if a byte slice starts with "data: ".
func IsSSEData(data []byte) bool {
	return bytes.HasPrefix(data, []byte(EventFieldData))
}

// ExtractSSEData extracts data from "data: " prefixed bytes.
func ExtractSSEData(data []byte) []byte {
	if !IsSSEData(data) {
		return data
	}

	// Remove "data: " prefix
	result := bytes.TrimPrefix(data, []byte(EventFieldData))
	// Remove leading space
	result = bytes.TrimPrefix(result, []byte(" "))

	return result
}
