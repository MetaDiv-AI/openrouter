// Package errors provides OpenRouter-specific error types.
//
// Use the standard library's errors.As to check for OpenRouterError:
//
//	import stderrors "errors"
//
//	var e *errors.OpenRouterError
//	if stderrors.As(err, &e) {
//	    switch e.Code {
//	    case 429:
//	        // rate limited
//	    case 401:
//	        // auth failed
//	    }
//	}
package errors

import "fmt"

// OpenRouterError represents an error returned by the OpenRouter API.
type OpenRouterError struct {
	HTTPStatus int
	Code       int
	Message    string
	Metadata   map[string]any
}

// Error implements the error interface.
func (e *OpenRouterError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("openrouter: %s (code: %d)", e.Message, e.Code)
	}
	return fmt.Sprintf("openrouter: API error (code: %d)", e.Code)
}

// Retryable returns true if the error is transient and the request can be retried.
func (e *OpenRouterError) Retryable() bool {
	switch e.Code {
	case 429, 503, 408:
		return true
	default:
		return false
	}
}

// Is reports whether the error matches the target.
func (e *OpenRouterError) Is(target error) bool {
	if t, ok := target.(*OpenRouterError); ok {
		return e.Code == t.Code
	}
	return false
}

// Common error codes for type checking.
var (
	ErrRateLimit           = &OpenRouterError{Code: 429}
	ErrAuth                = &OpenRouterError{Code: 401}
	ErrInsufficientCredits = &OpenRouterError{Code: 402}
	ErrModelNotFound       = &OpenRouterError{Code: 404}
	ErrProvider            = &OpenRouterError{Code: 502}
	ErrBadRequest          = &OpenRouterError{Code: 400}
	ErrModeration          = &OpenRouterError{Code: 403}
	ErrTimeout             = &OpenRouterError{Code: 408}
	ErrServiceUnavailable  = &OpenRouterError{Code: 503}
	ErrPricingUnavailable  = &OpenRouterError{Code: 404, Message: "pricing not available for model"}
)
