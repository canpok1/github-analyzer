package gemini

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/canpok1/github-analyzer/internal/domain"
)

// TODO: NewClient: APIキーが空の場合エラー
// TODO: NewClient: 正しいAPIキーで正常にクライアント生成
// TODO: NewClient: domain.Analyzerインターフェースの実装確認
// TODO: NewClient: デフォルトモデル設定確認
// TODO: NewClient: WithModelオプションでモデル変更
// TODO: Analyze: 正常系 - 成功レスポンスをパース
// TODO: Analyze: APIエラー（400等）でエラーを返す
// TODO: Analyze: レートリミット（429）で専用エラーを返す
// TODO: Analyze: タイムアウトでエラーを返す
// TODO: Analyze: レスポンスのJSONパースエラー
// TODO: Analyze: 空コンテンツのレスポンス

func TestNewClient_EmptyAPIKey_ReturnsError(t *testing.T) {
	_, err := NewClient("")
	if err == nil {
		t.Fatal("NewClient with empty API key should return error")
	}
}

func TestNewClient_ValidAPIKey_ReturnsClient(t *testing.T) {
	c, err := NewClient("test-api-key")
	if err != nil {
		t.Fatalf("NewClient returned unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("NewClient returned nil client")
	}
}

func TestNewClient_ImplementsAnalyzer(t *testing.T) {
	c, err := NewClient("test-api-key")
	if err != nil {
		t.Fatalf("NewClient returned unexpected error: %v", err)
	}
	var _ domain.Analyzer = c
}

func TestNewClient_DefaultModel(t *testing.T) {
	c, err := NewClient("test-api-key")
	if err != nil {
		t.Fatalf("NewClient returned unexpected error: %v", err)
	}
	if c.model != DefaultModel {
		t.Errorf("model = %q, want %q", c.model, DefaultModel)
	}
}

// newTestClient はテスト用のクライアントを生成するヘルパー。
func newTestClient(t *testing.T, server *httptest.Server) *Client {
	t.Helper()
	c, err := NewClient("test-api-key")
	if err != nil {
		t.Fatalf("NewClient returned unexpected error: %v", err)
	}
	c.baseURL = server.URL
	c.httpClient = server.Client()
	return c
}

func TestAnalyze_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"candidates": [
				{
					"content": {
						"parts": [
							{"text": "Analysis result here"}
						]
					}
				}
			]
		}`)
	}))
	defer server.Close()

	c := newTestClient(t, server)
	resp, err := c.Analyze(context.Background(), domain.AnalysisRequest{
		Prompt: "analyze this",
		Data:   "some data",
	})

	if err != nil {
		t.Fatalf("Analyze returned unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("Analyze returned nil response")
	}
	if resp.Content != "Analysis result here" {
		t.Errorf("Content = %q, want %q", resp.Content, "Analysis result here")
	}
}
