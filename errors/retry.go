package errors

import (
	"context"
	stderrors "errors"
)

// Retryable returns true if the given error indicates a transient failure
// that may succeed on retry (429, 503, 408, or context.DeadlineExceeded).
func Retryable(err error) bool {
	if err == nil {
		return false
	}
	var oerr *OpenRouterError
	if stderrors.As(err, &oerr) {
		return oerr.Retryable()
	}
	return stderrors.Is(err, context.DeadlineExceeded)
}
