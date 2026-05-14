package config

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	openai "github.com/sashabaranov/go-openai"
)

type deepseekClient struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func newDeepseekClient(baseURL, apiKey string) *deepseekClient {
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
