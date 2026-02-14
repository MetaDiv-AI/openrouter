package cost

import (
	"context"
	"strconv"

	"github.com/MetaDiv-AI/openrouter/errors"
	"github.com/MetaDiv-AI/openrouter/models"
)

// Service provides cost estimation.
type Service struct {
	models *models.Service
}

// NewService creates a new cost service that uses the models service for model data.
func NewService(modelsService *models.Service) *Service {
	return &Service{
		models: modelsService,
	}
}

// Estimate computes the estimated cost for a model and token counts.
func (s *Service) Estimate(ctx context.Context, modelID string, inputTokens, outputTokens int) (float64, error) {
	list, err := s.models.List(ctx)
	if err != nil {
		return 0, err
	}

	var m *models.Model
	for i := range list {
		if list[i].ID == modelID {
			m = &list[i]
			break
		}
	}
	if m == nil {
		return 0, errors.ErrModelNotFound
	}
	if m.Pricing == nil {
		return 0, errors.ErrPricingUnavailable
	}

	promptPrice, _ := strconv.ParseFloat(m.Pricing.Prompt, 64)
	completionPrice, _ := strconv.ParseFloat(m.Pricing.Completion, 64)

	cost := float64(inputTokens)*promptPrice/1e6 + float64(outputTokens)*completionPrice/1e6
	return cost, nil
}
