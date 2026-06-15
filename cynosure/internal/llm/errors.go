package llm

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const (
	// statusOverloaded is the non-standard "server overloaded" status used by
	// some LLM gateways (no net/http constant exists for 529).
	statusOverloaded = 529

	// retryBaseDelay is the wait before the first transient retry.
	retryBaseDelay = 500 * time.Millisecond
	// retryMaxDelay caps the exponential backoff.
	retryMaxDelay = 30 * time.Second
	// maxTransientRetries is the maximum number of retries for 429/529.
	maxTransientRetries = 10
)

// APIError is a typed error carrying the HTTP status code returned by the LLM
// API, so callers can distinguish 413 (context overflow) and 429/529
// (transient) from other failures.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Body)
}

// IsContextOverflow reports whether err is a 413 context-overflow error. The
// HTTP 413 status code is the primary signal; a "prompt_too_long" body is used
// as a fallback in case a gateway rewrites the status code.
func IsContextOverflow(err error) bool {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return false
	}
	if apiErr.StatusCode == http.StatusRequestEntityTooLarge {
		return true
	}
	return strings.Contains(apiErr.Body, "prompt_too_long")
}

// isRetryableStatus reports whether a status code is a transient failure that
// should be retried with backoff (429 rate limit, 529 overloaded).
func isRetryableStatus(code int) bool {
	return code == http.StatusTooManyRequests || code == statusOverloaded
}

// backoffDelay returns the wait before retry attempt n (n starts at 0):
// retryBaseDelay * 2^n with up to 50% jitter, capped at retryMaxDelay.
func backoffDelay(attempt int) time.Duration {
	delay := retryBaseDelay << attempt
	if delay <= 0 || delay > retryMaxDelay {
		delay = retryMaxDelay
	}
	jitter := time.Duration(rand.Int63n(int64(delay)/2 + 1))
	return delay + jitter
}
