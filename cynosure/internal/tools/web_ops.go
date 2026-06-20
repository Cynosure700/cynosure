package tools

import (
	"context"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	webFetchTimeout     = 30 * time.Second
	webFetchMaxBodySize = 2 << 20 // 2 MB
)

var (
	scriptStyleRe = regexp.MustCompile(`(?is)<(script|style)[^>]*>.*?</(script|style)>`)
	htmlTagRe     = regexp.MustCompile(`(?s)<[^>]+>`)
	whitespaceRe  = regexp.MustCompile(`[ \t\r\f\v]+`)
	blankLinesRe  = regexp.MustCompile(`\n{3,}`)
)

// RunWebFetch 获取一个 URL，将 HTML 正文转换为纯文本；当 ctx 中存在
// WebProcessor 时，通过 LLM 对该文本运行 prompt。若没有处理器，则直接返回
// 清洗后的文本。
func RunWebFetch(ctx context.Context, rawURL, prompt string) (string, error) {
	if strings.TrimSpace(rawURL) == "" {
		return "", fmt.Errorf("url is required")
	}
	rawURL = upgradeToHTTPS(rawURL)

	text, err := fetchAndCleanText(ctx, rawURL)
	if err != nil {
		return "", err
	}

	if proc, ok := webProcessorFromContext(ctx); ok {
		out, err := proc(ctx, prompt, text)
		if err != nil {
			return "", fmt.Errorf("process content: %w", err)
		}
		return out, nil
	}
	return text, nil
}

func upgradeToHTTPS(rawURL string) string {
	if strings.HasPrefix(rawURL, "http://") {
		return "https://" + strings.TrimPrefix(rawURL, "http://")
	}
	if !strings.Contains(rawURL, "://") {
		return "https://" + rawURL
	}
	return rawURL
}

func fetchAndCleanText(ctx context.Context, rawURL string) (string, error) {
	reqCtx, cancel := context.WithTimeout(ctx, webFetchTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "cynosure-web-fetch/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d fetching %s", resp.StatusCode, rawURL)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, webFetchMaxBodySize))
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	text := htmlToText(string(body))
	return text, nil
}

// htmlToText 去除脚本/样式与标签，解码实体并将空白折叠成可读的纯文本。
func htmlToText(s string) string {
	s = scriptStyleRe.ReplaceAllString(s, " ")
	s = htmlTagRe.ReplaceAllString(s, " ")
	s = html.UnescapeString(s)
	s = whitespaceRe.ReplaceAllString(s, " ")
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	s = strings.Join(lines, "\n")
	s = blankLinesRe.ReplaceAllString(s, "\n\n")
	return strings.TrimSpace(s)
}

// RunWebSearch 在配置搜索后端之前只是一个占位实现。它会明确报告该能力不可用，
// 以便模型据此调整。
func RunWebSearch(ctx context.Context, query string) (string, error) {
	if strings.TrimSpace(query) == "" {
		return "", fmt.Errorf("query is required")
	}
	return "web_search is not configured: no search backend available in this deployment.", nil
}
