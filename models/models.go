package models

import (
	"context"
	"sync"
	"time"

	"github.com/MetaDiv-AI/openrouter/internal"
)

// DefaultListCacheTTL is the default duration to cache the models list.
const DefaultListCacheTTL = 5 * time.Minute

// Service provides model listing and discovery.
type Service struct {
	caller      *internal.Caller
	cache       []Model
	cacheExpiry time.Time
	cacheMu     sync.RWMutex
	cacheTTL    time.Duration
}

// NewService creates a new models service.
func NewService(caller *internal.Caller) *Service {
	return &Service{
		caller:   caller,
		cacheTTL: DefaultListCacheTTL,
	}
}

// ListResponse is the response from listing models.
type ListResponse struct {
	Data []Model `json:"data"`
	Next string  `json:"next,omitempty"` // pagination cursor, if any
}

// Model represents an OpenRouter model.
type Model struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Created       int64         `json:"created"`
	Description   string        `json:"description,omitempty"`
	ContextLength int           `json:"context_length"`
	Architecture  *Architecture `json:"architecture,omitempty"`
	Pricing       *Pricing      `json:"pricing,omitempty"`
	TopProvider   *TopProvider  `json:"top_provider,omitempty"`
}

// Architecture describes model capabilities.
type Architecture struct {
	Modality         string   `json:"modality,omitempty"`
	InputModalities  []string `json:"input_modalities,omitempty"`
	OutputModalities []string `json:"output_modalities,omitempty"`
}

// Pricing holds model pricing (per-token, as string for precision).
type Pricing struct {
	Prompt          string `json:"prompt,omitempty"`
	Completion      string `json:"completion,omitempty"`
	InputCacheRead  string `json:"input_cache_read,omitempty"`
	InputCacheWrite string `json:"input_cache_write,omitempty"`
}

// TopProvider holds top provider info.
type TopProvider struct {
	ContextLength       int  `json:"context_length"`
	MaxCompletionTokens int  `json:"max_completion_tokens"`
	IsModerated         bool `json:"is_moderated"`
}

// List returns all available models. Results are cached for DefaultListCacheTTL.
func (s *Service) List(ctx context.Context) ([]Model, error) {
	s.cacheMu.RLock()
	if time.Now().Before(s.cacheExpiry) && len(s.cache) > 0 {
		list := s.cache
		s.cacheMu.RUnlock()
		return list, nil
	}
	s.cacheMu.RUnlock()

	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	if time.Now().Before(s.cacheExpiry) && len(s.cache) > 0 {
		return s.cache, nil
	}

	var resp ListResponse
	if err := s.caller.DoGet(ctx, "/models", &resp); err != nil {
		return nil, err
	}
	s.cache = resp.Data
	s.cacheExpiry = time.Now().Add(s.cacheTTL)
	return s.cache, nil
}
