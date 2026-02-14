package chat

import (
	"context"

	"github.com/MetaDiv-AI/openrouter/internal"
)

// Service provides chat completion operations.
type Service struct {
	caller *internal.Caller
}

// NewService creates a new chat service.
func NewService(caller *internal.Caller) *Service {
	return &Service{caller: caller}
}

// Create sends a non-streaming chat completion request.
func (s *Service) Create(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	if req == nil {
		req = &ChatRequest{}
	}
	req.Stream = false

	var resp ChatResponse
	if err := s.caller.DoPost(ctx, "/chat/completions", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateStream sends a streaming chat completion request.
func (s *Service) CreateStream(ctx context.Context, req *ChatRequest) (*StreamReader, error) {
	if req == nil {
		req = &ChatRequest{}
	}
	req.Stream = true

	sr := NewStreamReader()
	go func() {
		err := s.caller.DoStreamPost(ctx, "/chat/completions", req, func(chunk []byte) error {
			return sr.ProcessLine(chunk)
		})
		if err != nil {
			sr.SetError(err)
		} else {
			sr.Close()
		}
	}()
	return sr, nil
}
