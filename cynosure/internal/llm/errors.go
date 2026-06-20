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
	// statusOverloaded 是部分 LLM 网关使用的非标准"服务器过载"状态码
	// （net/http 中没有对应 529 的常量）。
	statusOverloaded = 529

	// retryBaseDelay 是首次瞬时重试前的等待时间。
	retryBaseDelay = 500 * time.Millisecond
	// retryMaxDelay 限制指数退避的最大上限。
	retryMaxDelay = 30 * time.Second
	// maxTransientRetries 是针对 429/529 的最大重试次数。
	maxTransientRetries = 10
)

// APIError 是携带 LLM API 返回的 HTTP 状态码的类型化错误，
// 便于调用方区分 413（上下文溢出）、429/529（瞬时错误）
// 与其他失败情况。
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Body)
}

// IsContextOverflow 报告 err 是否为 413 上下文溢出错误。
// HTTP 413 状态码是主要判断信号；当网关重写状态码时，
// 会以响应体中是否包含 "prompt_too_long" 作为兜底判断。
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

// isRetryableStatus 报告某状态码是否为应当以退避方式重试的瞬时失败
// （429 限流、529 过载）。
func isRetryableStatus(code int) bool {
	return code == http.StatusTooManyRequests || code == statusOverloaded
}

// backoffDelay 返回第 n 次重试前的等待时间（n 从 0 开始）：
// retryBaseDelay * 2^n，并附加最多 50% 的抖动，上限为 retryMaxDelay。
func backoffDelay(attempt int) time.Duration {
	delay := retryBaseDelay << attempt
	if delay <= 0 || delay > retryMaxDelay {
		delay = retryMaxDelay
	}
	jitter := time.Duration(rand.Int63n(int64(delay)/2 + 1))
	return delay + jitter
}
