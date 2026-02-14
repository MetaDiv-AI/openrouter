package models

import (
	"context"
	"strconv"
	"strings"

	"github.com/MetaDiv-AI/openrouter/errors"
)

// Get returns a single model by ID from the list.
func (s *Service) Get(ctx context.Context, id string) (*Model, error) {
	list, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if list[i].ID == id {
			return &list[i], nil
		}
	}
	return nil, errors.ErrModelNotFound
}

// ByProvider filters models by provider (e.g. "anthropic", "openai").
func (s *Service) ByProvider(ctx context.Context, provider string) ([]Model, error) {
	list, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	provider = strings.ToLower(provider)
	var out []Model
	for _, m := range list {
		if strings.HasPrefix(strings.ToLower(m.ID), provider+"/") {
			out = append(out, m)
		}
	}
	return out, nil
}

// ByContextLength filters models with context_length >= minTokens.
func (s *Service) ByContextLength(ctx context.Context, minTokens int) ([]Model, error) {
	list, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	var out []Model
	for _, m := range list {
		if m.ContextLength >= minTokens {
			out = append(out, m)
		}
	}
	return out, nil
}

// Cheapest returns the cheapest model by prompt+completion price.
func (s *Service) Cheapest(ctx context.Context) (*Model, error) {
	list, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	var cheapest *Model
	var bestPrice float64 = -1
	for i := range list {
		m := &list[i]
		if m.Pricing == nil {
			continue
		}
		p, _ := strconv.ParseFloat(m.Pricing.Prompt, 64)
		c, _ := strconv.ParseFloat(m.Pricing.Completion, 64)
		price := p + c
		if bestPrice < 0 || price < bestPrice {
			bestPrice = price
			cheapest = m
		}
	}
	return cheapest, nil
}

// SupportsVision returns models that support image input.
func (s *Service) SupportsVision(ctx context.Context) ([]Model, error) {
	list, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	var out []Model
	for _, m := range list {
		if m.Architecture != nil {
			for _, mod := range m.Architecture.InputModalities {
				if strings.EqualFold(mod, "image") {
					out = append(out, m)
					break
				}
			}
		}
	}
	return out, nil
}
