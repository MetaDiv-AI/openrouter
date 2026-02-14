package openrouter

import (
	"time"

	"github.com/MetaDiv-AI/logger"
)

const (
	DefaultBaseURL    = "https://openrouter.ai/api/v1"
	DefaultTimeout    = 60 * time.Second
	DefaultMaxRetries = 3
)

// Config holds the client configuration.
type Config struct {
	APIKey     string
	BaseURL    string
	Timeout    time.Duration
	MaxRetries int
	Headers    map[string]string
	Debug      bool
	Logger     logger.Logger
}

// Option is a functional option for configuring the client.
type Option func(*Config)

// WithAPIKey sets the API key. If not set, OPENROUTER_API_KEY env var is used.
func WithAPIKey(apiKey string) Option {
	return func(c *Config) {
		c.APIKey = apiKey
	}
}

// WithBaseURL sets the base URL for API requests.
func WithBaseURL(baseURL string) Option {
	return func(c *Config) {
		c.BaseURL = baseURL
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithMaxRetries sets the maximum number of retries for retryable errors.
func WithMaxRetries(n int) Option {
	return func(c *Config) {
		c.MaxRetries = n
	}
}

// WithHeaders sets custom headers merged with defaults.
func WithHeaders(headers map[string]string) Option {
	return func(c *Config) {
		if c.Headers == nil {
			c.Headers = make(map[string]string)
		}
		for k, v := range headers {
			c.Headers[k] = v
		}
	}
}

// WithReferer sets the HTTP-Referer header for app attribution on openrouter.ai.
func WithReferer(url string) Option {
	return WithHeaders(map[string]string{"HTTP-Referer": url})
}

// WithTitle sets the X-Title header for app title on openrouter.ai.
func WithTitle(title string) Option {
	return WithHeaders(map[string]string{"X-Title": title})
}

// WithForwardedFor sets the X-Forwarded-For header.
func WithForwardedFor(ip string) Option {
	return WithHeaders(map[string]string{"X-Forwarded-For": ip})
}

// WithDebug enables debug logging. Uses logger.New().Development().Build() if Logger is nil.
func WithDebug(debug bool) Option {
	return func(c *Config) {
		c.Debug = debug
	}
}

// WithLogger sets a custom logger (used by http_caller.WithDebugLogger).
func WithLogger(log logger.Logger) Option {
	return func(c *Config) {
		c.Logger = log
	}
}
