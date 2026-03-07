package gemini

import (
	"context"
	"fmt"
	"net/http"

	"github.com/canpok1/github-analyzer/internal/domain"
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
		httpClient: http.DefaultClient,
		baseURL:    "https://generativelanguage.googleapis.com/v1beta",
	}, nil
}

// Analyze はプロンプトとデータを元にGemini APIで分析を実行する。
func (c *Client) Analyze(_ context.Context, _ domain.AnalysisRequest) (*domain.AnalysisResponse, error) {
	return nil, fmt.Errorf("not implemented")
}
