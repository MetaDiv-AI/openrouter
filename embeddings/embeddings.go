package embeddings

import (
	"context"

	"github.com/MetaDiv-AI/openrouter/errors"
	"github.com/MetaDiv-AI/openrouter/internal"
)

// Service provides embedding operations.
type Service struct {
	caller *internal.Caller
}

// NewService creates a new embeddings service.
func NewService(caller *internal.Caller) *Service {
	return &Service{caller: caller}
}

// CreateRequest is the request for creating embeddings.
type CreateRequest struct {
	Model string `json:"model"`
	Input any    `json:"input"` // string or []string
}

// CreateResponse is the response from creating embeddings.
type CreateResponse struct {
	Data  []EmbeddingData `json:"data"`
	Usage *Usage          `json:"usage,omitempty"`
}

// EmbeddingData holds a single embedding.
type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// Usage represents token usage.
type Usage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// Create generates embeddings for the given input.
func (s *Service) Create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	if req == nil {
		return nil, &errors.OpenRouterError{Code: 400, Message: "request cannot be nil"}
	}

	var resp CreateResponse
	if err := s.caller.DoPost(ctx, "/embeddings", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
