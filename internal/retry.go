package internal

import (
	"context"
	"math"
	"math/rand"
	"time"

	"github.com/MetaDiv-AI/openrouter/errors"
)

// DefaultBackoff is the default retry configuration.
var DefaultBackoff = BackoffConfig{
	InitialInterval: time.Second,
	MaxInterval:     30 * time.Second,
	Multiplier:      2,
}

// BackoffConfig configures exponential backoff behavior.
type BackoffConfig struct {
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
}

// Do retries the given operation with exponential backoff when the error is retryable.
func Do(ctx context.Context, maxRetries int, cfg BackoffConfig, fn func() error) error {
	var lastErr error
	interval := cfg.InitialInterval

	for attempt := 0; attempt <= maxRetries; attempt++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if !errors.Retryable(lastErr) {
			return lastErr
		}
		if attempt == maxRetries {
			return lastErr
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(jitter(interval)):
			interval = time.Duration(float64(interval) * cfg.Multiplier)
			if interval > cfg.MaxInterval {
				interval = cfg.MaxInterval
			}
		}
	}
	return lastErr
}

func jitter(d time.Duration) time.Duration {
	j := time.Duration(rand.Float64() * 0.3 * float64(d))
	return d + j - time.Duration(math.Round(0.15*float64(d)))
}
