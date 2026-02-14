package openrouter

import (
	"context"
	"os"

	"github.com/MetaDiv-AI/logger"
	"github.com/MetaDiv-AI/openrouter/batch"
	"github.com/MetaDiv-AI/openrouter/chat"
	"github.com/MetaDiv-AI/openrouter/cost"
	"github.com/MetaDiv-AI/openrouter/embeddings"
	"github.com/MetaDiv-AI/openrouter/errors"
	"github.com/MetaDiv-AI/openrouter/internal"
	"github.com/MetaDiv-AI/openrouter/models"
)

// Client is the OpenRouter API client.
type Client struct {
	caller     *internal.Caller
	Chat       *chat.Service
	Embeddings *embeddings.Service
	Models     *models.Service
	Cost       *cost.Service
}

// NewClient creates a new OpenRouter client with the given options.
func NewClient(opts ...Option) (*Client, error) {
	cfg := &Config{
		BaseURL:    DefaultBaseURL,
		Timeout:    DefaultTimeout,
		MaxRetries: DefaultMaxRetries,
		Headers:    make(map[string]string),
	}
	for _, opt := range opts {
		opt(cfg)
	}

	apiKey := cfg.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("OPENROUTER_API_KEY")
	}
	if apiKey == "" {
		return nil, &errors.OpenRouterError{Code: 401, Message: "missing API key: set OPENROUTER_API_KEY or use WithAPIKey()"}
	}
	if cfg.Timeout < 0 {
		return nil, &errors.OpenRouterError{Code: 400, Message: "timeout must be non-negative"}
	}
	if cfg.MaxRetries < 0 {
		return nil, &errors.OpenRouterError{Code: 400, Message: "max retries must be non-negative"}
	}

	if cfg.Debug && cfg.Logger == nil {
		cfg.Logger = logger.New().Development().Build()
	}

	caller := internal.NewCaller(
		cfg.BaseURL,
		apiKey,
		cfg.Timeout,
		cfg.Headers,
		cfg.Logger,
		cfg.MaxRetries,
	)

	modelsSvc := models.NewService(caller)
	return &Client{
		caller:     caller,
		Chat:       chat.NewService(caller),
		Embeddings: embeddings.NewService(caller),
		Models:     modelsSvc,
		Cost:       cost.NewService(modelsSvc),
	}, nil
}

// EstimateCost estimates the cost for a given model and token counts.
func (c *Client) EstimateCost(ctx context.Context, model string, inputTokens, outputTokens int) (float64, error) {
	return c.Cost.Estimate(ctx, model, inputTokens, outputTokens)
}

// BatchChat runs multiple chat requests concurrently.
func (c *Client) BatchChat(ctx context.Context, requests []*chat.ChatRequest, concurrency int) ([]*chat.ChatResponse, []error) {
	return batch.NewChatBatchProcessor(c.Chat, concurrency).Run(ctx, requests)
}

// DebugCurl returns a curl command string for a chat request (for debugging).
func (c *Client) DebugCurl(req *chat.ChatRequest) string {
	if req == nil {
		req = &chat.ChatRequest{}
	}
	return c.caller.DebugCurl("/chat/completions", req)
}
