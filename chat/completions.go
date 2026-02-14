package chat

import "context"

// CompletionsRequest is the legacy prompt-based completions request.
type CompletionsRequest struct {
	Prompt      string   `json:"prompt"`
	Model       string   `json:"model,omitempty"`
	MaxTokens   *int     `json:"max_tokens,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
	Stream      bool     `json:"stream,omitempty"`
}

// CompletionsResponse is the legacy completions response.
type CompletionsResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   *Usage   `json:"usage,omitempty"`
}

// CreateCompletions sends a legacy prompt-based completion request.
func (s *Service) CreateCompletions(ctx context.Context, req *CompletionsRequest) (*CompletionsResponse, error) {
	if req == nil {
		req = &CompletionsRequest{}
	}

	var resp CompletionsResponse
	if err := s.caller.DoPost(ctx, "/chat/completions", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
