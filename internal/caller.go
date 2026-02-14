package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MetaDiv-AI/http_caller"
	"github.com/MetaDiv-AI/logger"
	"github.com/MetaDiv-AI/openrouter/errors"
)

// Caller wraps http_caller with auth, headers, retry, and error parsing.
type Caller struct {
	baseURL string
	apiKey  string
	headers map[string]string
	client  *http.Client
	logger  logger.Logger
	retries int
}

// NewCaller creates a new Caller with the given configuration.
func NewCaller(baseURL, apiKey string, timeout time.Duration, headers map[string]string, log logger.Logger, retries int) *Caller {
	baseURL = strings.TrimSuffix(baseURL, "/")
	if headers == nil {
		headers = make(map[string]string)
	}
	return &Caller{
		baseURL: baseURL,
		apiKey:  apiKey,
		headers: copyHeaders(headers),
		client: &http.Client{
			Timeout: timeout,
		},
		logger:  log,
		retries: retries,
	}
}

func copyHeaders(m map[string]string) map[string]string {
	out := make(map[string]string, len(m)+1)
	for k, v := range m {
		out[k] = v
	}
	return out
}

// errorResponse is the shape of OpenRouter error responses.
type errorResponse struct {
	Error struct {
		Code     int            `json:"code"`
		Message  string         `json:"message"`
		Metadata map[string]any `json:"metadata,omitempty"`
	} `json:"error"`
}

func parseError(statusCode int, rawBody string) *errors.OpenRouterError {
	var resp errorResponse
	if err := json.Unmarshal([]byte(rawBody), &resp); err != nil {
		return &errors.OpenRouterError{
			HTTPStatus: statusCode,
			Code:       statusCode,
			Message:    rawBody,
		}
	}
	code := resp.Error.Code
	if code == 0 {
		code = statusCode
	}
	return &errors.OpenRouterError{
		HTTPStatus: statusCode,
		Code:       code,
		Message:    resp.Error.Message,
		Metadata:   resp.Error.Metadata,
	}
}

// DoPost executes a POST request with retries and error mapping.
func (c *Caller) DoPost(ctx context.Context, path string, req, resp any) error {
	url := c.baseURL + path
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	reqBody := json.RawMessage(reqBytes)

	var lastResp *http_caller.Response[json.RawMessage]
	err = Do(ctx, c.retries, DefaultBackoff, func() error {
		builder := http_caller.New[json.RawMessage, json.RawMessage](url).
			Header("Authorization", "Bearer "+c.apiKey).
			Headers(c.headers).
			Body(&reqBody).
			WithClient(c.client)
		if c.logger != nil {
			builder = builder.WithDebugLogger(c.logger)
		}

		r, doErr := builder.Post(ctx)
		if doErr != nil {
			return doErr
		}
		lastResp = r

		if r.StatusCode >= 400 {
			return parseError(r.StatusCode, r.RawBody)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if lastResp == nil {
		return nil
	}
	if resp != nil && len(lastResp.Body) > 0 {
		return json.Unmarshal(lastResp.Body, resp)
	}
	return nil
}

// DoGet executes a GET request with retries.
func (c *Caller) DoGet(ctx context.Context, path string, resp any) error {
	url := c.baseURL + path

	var lastResp *http_caller.Response[json.RawMessage]
	err := Do(ctx, c.retries, DefaultBackoff, func() error {
		builder := http_caller.New[struct{}, json.RawMessage](url).
			Header("Authorization", "Bearer "+c.apiKey).
			Headers(c.headers).
			WithClient(c.client)
		if c.logger != nil {
			builder = builder.WithDebugLogger(c.logger)
		}

		r, doErr := builder.Get(ctx)
		if doErr != nil {
			return doErr
		}
		lastResp = r

		if r.StatusCode >= 400 {
			return parseError(r.StatusCode, r.RawBody)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if lastResp == nil {
		return nil
	}
	if resp != nil && len(lastResp.Body) > 0 {
		return json.Unmarshal(lastResp.Body, resp)
	}
	return nil
}

// DoStreamPost executes a streaming POST request (no retry).
func (c *Caller) DoStreamPost(ctx context.Context, path string, req any, handler http_caller.ChunkHandler) error {
	url := c.baseURL + path
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	reqBody := json.RawMessage(reqBytes)

	builder := http_caller.New[json.RawMessage, any](url).
		Header("Authorization", "Bearer "+c.apiKey).
		Headers(c.headers).
		Body(&reqBody).
		WithClient(c.client)
	if c.logger != nil {
		builder = builder.WithDebugLogger(c.logger)
	}

	return builder.StreamPost(ctx, handler)
}

// redactAPIKey shows only the first 7 chars (e.g. "sk-...") for safe display.
func redactAPIKey(key string) string {
	if len(key) <= 7 {
		return "[REDACTED]"
	}
	return key[:7] + "..."
}

// escapeForCurl escapes single quotes for use in shell curl -d ”.
func escapeForCurl(s string) string {
	return strings.ReplaceAll(s, "'", "'\\''")
}

// DebugCurl returns a curl command string for the given POST request.
// The API key is redacted in the output. Header values are escaped for shell safety.
func (c *Caller) DebugCurl(path string, body any) string {
	url := c.baseURL + path
	bodyBytes, _ := json.MarshalIndent(body, "", "  ")
	bodyEscaped := escapeForCurl(string(bodyBytes))

	var hdr strings.Builder
	hdr.WriteString(fmt.Sprintf("-H 'Authorization: Bearer %s' ", redactAPIKey(c.apiKey)))
	for k, v := range c.headers {
		hdr.WriteString(fmt.Sprintf("-H '%s: %s' ", escapeForCurl(k), escapeForCurl(v)))
	}
	return fmt.Sprintf("curl -X POST '%s' %s-H 'Content-Type: application/json' -d '%s'", url, hdr.String(), bodyEscaped)
}
