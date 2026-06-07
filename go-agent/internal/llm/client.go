package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

// Client is the minimal chat-completion interface the runtime depends on.
type Client interface {
	CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
	CreateChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (ChatCompletionStream, error)
}

// ChatCompletionStream is a streamed chat-completion response.
type ChatCompletionStream interface {
	Recv() (openai.ChatCompletionStreamResponse, error)
	Close() error
}

type deepseekClient struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

// NewDeepseekClient returns a Client backed by the DeepSeek chat-completions API.
func NewDeepseekClient(baseURL, apiKey string) Client {
	return &deepseekClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		http:    &http.Client{},
	}
}

func (c *deepseekClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	body := map[string]any{
		"model":            req.Model,
		"messages":         req.Messages,
		"thinking":         map[string]string{"type": "enabled"},
		"reasoning_effort": "high",
		"stream":           false,
	}

	if len(req.Tools) > 0 {
		body["tools"] = req.Tools
	}

	jsonBody, err := json.Marshal(body)
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
		return openai.ChatCompletionResponse{}, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result openai.ChatCompletionResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("unmarshal response: %w", err)
	}

	return result, nil
}

func (c *deepseekClient) CreateChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (ChatCompletionStream, error) {
	body := map[string]any{
		"model":            req.Model,
		"messages":         req.Messages,
		"thinking":         map[string]string{"type": "enabled"},
		"reasoning_effort": "high",
		"stream":           true,
	}

	if len(req.Tools) > 0 {
		body["tools"] = req.Tools
	}

	jsonBody, err := json.Marshal(body)
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
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
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
