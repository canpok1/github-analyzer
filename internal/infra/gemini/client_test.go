package gemini

import (
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
