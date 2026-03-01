package llmrouter

import (
	"errors"
	"fmt"
)

// Sentinel errors for common API failure modes.
var (
	ErrUnauthorized  = errors.New("llmrouter: unauthorized")
	ErrForbidden     = errors.New("llmrouter: forbidden")
	ErrNotFound      = errors.New("llmrouter: not found")
	ErrRateLimited   = errors.New("llmrouter: rate limited")
	ErrQuotaExceeded = errors.New("llmrouter: quota exceeded")
	ErrBadRequest    = errors.New("llmrouter: bad request")
	ErrServerError   = errors.New("llmrouter: server error")
)

// APIError is returned when the LLMRouter API responds with a non-2xx status.
type APIError struct {
	StatusCode int    `json:"status_code"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("llmrouter: %d %s: %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("llmrouter: %d: %s", e.StatusCode, e.Message)
}

// Unwrap maps the APIError to its corresponding sentinel error so callers can
// use errors.Is for common cases.
func (e *APIError) Unwrap() error {
	switch e.StatusCode {
	case 400:
		return ErrBadRequest
	case 401:
		return ErrUnauthorized
	case 403:
		return ErrForbidden
	case 404:
		return ErrNotFound
	case 429:
		if e.Code == "quota_exceeded" {
			return ErrQuotaExceeded
		}
		return ErrRateLimited
	default:
		if e.StatusCode >= 500 {
			return ErrServerError
		}
		return nil
	}
}
