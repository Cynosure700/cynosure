package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// Client 是 runtime 依赖的最小化聊天补全接口。
type Client interface {
	CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
	CreateChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (ChatCompletionStream, error)
}

// ChatCompletionStream 是流式聊天补全响应。
type ChatCompletionStream interface {
	Recv() (openai.ChatCompletionStreamResponse, error)
	Close() error
}

type deepseekClient struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

// NewDeepseekClient 返回一个基于 DeepSeek 聊天补全 API 实现的 Client。
func NewDeepseekClient(baseURL, apiKey string) Client {
	return &deepseekClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		http:    &http.Client{},
	}
}

// buildRequestBody 组装 DeepSeek 请求体，在设置了 max_tokens 时透传该值，
// 以便调用方控制/覆盖输出预算。
func buildRequestBody(req openai.ChatCompletionRequest, stream bool) map[string]any {
	body := map[string]any{
		"model":            req.Model,
		"messages":         req.Messages,
		"thinking":         map[string]string{"type": "enabled"},
		"reasoning_effort": "high",
		"stream":           stream,
	}
	if len(req.Tools) > 0 {
		body["tools"] = req.Tools
	}
	if req.MaxTokens > 0 {
		body["max_tokens"] = req.MaxTokens
	}
	return body
}

// retryTransient 反复执行 do，直到成功、返回不可重试错误，或瞬时重试预算耗尽。
// 429/529 会以指数退避加抖动的方式重试；ctx 被取消时立即中止。
func retryTransient[T any](ctx context.Context, do func() (T, error)) (T, error) {
	var zero T
	for attempt := 0; ; attempt++ {
		result, err := do()
		if err == nil {
			return result, nil
		}
		var apiErr *APIError
		if !errors.As(err, &apiErr) || !isRetryableStatus(apiErr.StatusCode) || attempt >= maxTransientRetries {
			return zero, err
		}
		select {
		case <-time.After(backoffDelay(attempt)):
		case <-ctx.Done():
			return zero, ctx.Err()
		}
	}
}

func (c *deepseekClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	return retryTransient(ctx, func() (openai.ChatCompletionResponse, error) {
		return c.createChatCompletionOnce(ctx, req)
	})
}

func (c *deepseekClient) createChatCompletionOnce(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	jsonBody, err := json.Marshal(buildRequestBody(req, false))
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("marshal request: %w", err)
	}

	url := c.baseURL + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return openai.ChatCompletionResponse{}, &APIError{StatusCode: resp.StatusCode, Body: string(respBody)}
	}

	var result openai.ChatCompletionResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("unmarshal response: %w", err)
	}

	return result, nil
}

func (c *deepseekClient) CreateChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (ChatCompletionStream, error) {
	return retryTransient(ctx, func() (ChatCompletionStream, error) {
		return c.createChatCompletionStreamOnce(ctx, req)
	})
}

func (c *deepseekClient) createChatCompletionStreamOnce(ctx context.Context, req openai.ChatCompletionRequest) (ChatCompletionStream, error) {
	jsonBody, err := json.Marshal(buildRequestBody(req, true))
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := c.baseURL + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		respBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("read response: %w", readErr)
		}
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(respBody)}
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	return &deepseekChatStream{body: resp.Body, scanner: scanner}, nil
}

type deepseekChatStream struct {
	body    io.Closer
	scanner *bufio.Scanner
}

func (s *deepseekChatStream) Recv() (openai.ChatCompletionStreamResponse, error) {
	for s.scanner.Scan() {
		line := strings.TrimSpace(s.scanner.Text())
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			return openai.ChatCompletionStreamResponse{}, io.EOF
		}
		var chunk openai.ChatCompletionStreamResponse
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			return openai.ChatCompletionStreamResponse{}, fmt.Errorf("unmarshal stream response: %w", err)
		}
		return chunk, nil
	}
	if err := s.scanner.Err(); err != nil {
		return openai.ChatCompletionStreamResponse{}, err
	}
	return openai.ChatCompletionStreamResponse{}, io.EOF
}

func (s *deepseekChatStream) Close() error {
	return s.body.Close()
}
