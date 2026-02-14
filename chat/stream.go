package chat

import (
	"encoding/json"
	"errors"
	"io"
	"sync"

	oerrors "github.com/MetaDiv-AI/openrouter/errors"
	"github.com/MetaDiv-AI/openrouter/internal"
)

// StreamError represents an error in a streaming chunk (e.g. provider disconnect).
type StreamError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// StreamChunk represents a single streaming chunk.
type StreamChunk struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []Choice     `json:"choices"`
	Usage   *Usage       `json:"usage,omitempty"`
	Error   *StreamError `json:"error,omitempty"`
}

// StreamReader reads SSE chat completion stream.
type StreamReader struct {
	ch     chan StreamChunk
	done   chan struct{}
	usage  *Usage
	err    error
	closed bool
	mu     sync.Mutex
}

// NewStreamReader creates a new StreamReader with a 256-chunk buffer.
// Consumers should call Next() promptly to avoid blocking the producer.
func NewStreamReader() *StreamReader {
	return &StreamReader{
		ch:   make(chan StreamChunk, 256),
		done: make(chan struct{}),
	}
}

// ProcessLine processes a single SSE line (from http_caller ChunkHandler).
func (sr *StreamReader) ProcessLine(line []byte) error {
	payload, done := internal.ParseSSELine(line)
	if done {
		sr.Close()
		return nil
	}
	if len(payload) == 0 {
		return nil
	}

	var chunk StreamChunk
	if err := json.Unmarshal(payload, &chunk); err != nil {
		sr.SetError(&oerrors.OpenRouterError{Code: 500, Message: "malformed stream chunk: " + err.Error()})
		return nil
	}
	if chunk.Error != nil {
		sr.SetError(&oerrors.OpenRouterError{
			Code:    chunk.Error.Code,
			Message: chunk.Error.Message,
		})
		return nil
	}
	if chunk.Usage != nil {
		sr.mu.Lock()
		sr.usage = chunk.Usage
		sr.mu.Unlock()
	}
	if len(chunk.Choices) > 0 {
		select {
		case sr.ch <- chunk:
		case <-sr.done:
			return nil
		}
	}
	return nil
}

// Next returns the next chunk or io.EOF when done.
func (sr *StreamReader) Next() (*StreamChunk, error) {
	sr.mu.Lock()
	if sr.err != nil {
		err := sr.err
		sr.mu.Unlock()
		return nil, err
	}
	sr.mu.Unlock()

	chunk, ok := <-sr.ch
	if !ok {
		sr.mu.Lock()
		usage := sr.usage
		sr.mu.Unlock()
		if usage != nil {
			return &StreamChunk{Usage: usage}, io.EOF
		}
		return nil, io.EOF
	}
	return &chunk, nil
}

// ReadAll consumes the stream and returns the full content and usage.
func (sr *StreamReader) ReadAll() (string, *Usage, error) {
	var content string
	var usage *Usage
	for {
		chunk, err := sr.Next()
		if errors.Is(err, io.EOF) {
			if chunk != nil && chunk.Usage != nil {
				usage = chunk.Usage
			}
			sr.mu.Lock()
			if sr.usage != nil {
				usage = sr.usage
			}
			sr.mu.Unlock()
			return content, usage, nil
		}
		if err != nil {
			return content, usage, err
		}
		if chunk != nil && len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil && chunk.Choices[0].Delta.Content != nil {
			switch c := chunk.Choices[0].Delta.Content.(type) {
			case string:
				content += c
			}
		}
		if chunk != nil && chunk.Usage != nil {
			usage = chunk.Usage
		}
	}
}

// Close closes the stream.
func (sr *StreamReader) Close() {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	if !sr.closed {
		sr.closed = true
		close(sr.done)
		close(sr.ch)
	}
}

// SetError sets the stream error.
func (sr *StreamReader) SetError(err error) {
	sr.mu.Lock()
	sr.err = err
	sr.mu.Unlock()
	sr.Close()
}
