package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/canpok1/github-analyzer/internal/domain"
)

const (
	// maxResponseSize はレスポンスボディの最大読み取りサイズ（10MB）。
	maxResponseSize = 10 * 1024 * 1024
	// defaultTimeout はHTTPクライアントのデフォルトタイムアウト。
	defaultTimeout = 60 * time.Second
)

const (
	// DefaultModel はデフォルトで使用するGeminiモデル。
	DefaultModel = "gemini-1.5-flash"
)

// Client はGemini APIクライアント。
type Client struct {
	apiKey     string
	model      string
	httpClient *http.Client
	baseURL    string
}

// NewClient は新しいGemini APIクライアントを生成する。
// apiKeyが空の場合はエラーを返す。
func NewClient(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is required")
	}

	return &Client{
		apiKey:     apiKey,
		model:      DefaultModel,
		httpClient: &http.Client{Timeout: defaultTimeout},
		baseURL:    "https://generativelanguage.googleapis.com/v1beta",
	}, nil
}

// Analyze はプロンプトとデータを元にGemini APIで分析を実行する。
func (c *Client) Analyze(ctx context.Context, req domain.AnalysisRequest) (*domain.AnalysisResponse, error) {
	model := c.model
	if req.Model != "" {
		model = req.Model
	}

	promptText := req.Prompt
	if req.Data != "" {
		promptText = req.Prompt + "\n\n" + req.Data
	}

	geminiReq := geminiRequest{
		Contents: []content{
			{
				Parts: []part{
					{Text: promptText},
				},
			},
		},
	}

	body, err := json.Marshal(geminiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:generateContent", c.baseURL, model)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-goog-api-key", c.apiKey)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() { _ = httpResp.Body.Close() }()

	respBody, err := io.ReadAll(io.LimitReader(httpResp.Body, maxResponseSize))
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		var errResp geminiErrorResponse
		if jsonErr := json.Unmarshal(respBody, &errResp); jsonErr == nil && errResp.Error.Message != "" {
			if httpResp.StatusCode == http.StatusTooManyRequests {
				return nil, fmt.Errorf("rate limit exceeded: %s", errResp.Error.Message)
			}
			return nil, fmt.Errorf("API error (status %d): %s", httpResp.StatusCode, errResp.Error.Message)
		}
		bodyStr := string(respBody)
		if len(bodyStr) > 500 {
			bodyStr = bodyStr[:500] + "..."
		}
		return nil, fmt.Errorf("API error (status %d): %s", httpResp.StatusCode, bodyStr)
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from Gemini API")
	}

	return &domain.AnalysisResponse{
		Content: geminiResp.Candidates[0].Content.Parts[0].Text,
	}, nil
}
