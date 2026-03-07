package gemini

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/canpok1/github-analyzer/internal/domain"
)

// DONE: NewClient: APIキーが空の場合エラー
// DONE: NewClient: 正しいAPIキーで正常にクライアント生成
// DONE: NewClient: domain.Analyzerインターフェースの実装確認
// DONE: NewClient: デフォルトモデル設定確認
// DONE: Analyze: 正常系 - 成功レスポンスをパース
// DONE: Analyze: APIエラー（400等）でエラーを返す
// DONE: Analyze: レートリミット（429）で専用エラーを返す
// DONE: Analyze: タイムアウトでエラーを返す
// DONE: Analyze: レスポンスのJSONパースエラー
// DONE: Analyze: 空コンテンツのレスポンス
// DONE: Analyze: モデルのオーバーライド
// DONE: Analyze: リクエストフォーマット確認

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

func TestAnalyze_APIError_400(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":{"code":400,"message":"Invalid request","status":"INVALID_ARGUMENT"}}`)
	}))
	defer server.Close()

	c := newTestClient(t, server)
	_, err := c.Analyze(context.Background(), domain.AnalysisRequest{Prompt: "test"})

	if err == nil {
		t.Fatal("Analyze should return error for 400 status")
	}
	if !strings.Contains(err.Error(), "API error") {
		t.Errorf("error = %q, want to contain %q", err.Error(), "API error")
	}
	if !strings.Contains(err.Error(), "Invalid request") {
		t.Errorf("error = %q, want to contain %q", err.Error(), "Invalid request")
	}
}

func TestAnalyze_RateLimit_429(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprint(w, `{"error":{"code":429,"message":"Resource has been exhausted","status":"RESOURCE_EXHAUSTED"}}`)
	}))
	defer server.Close()

	c := newTestClient(t, server)
	_, err := c.Analyze(context.Background(), domain.AnalysisRequest{Prompt: "test"})

	if err == nil {
		t.Fatal("Analyze should return error for 429 status")
	}
	if !strings.Contains(err.Error(), "rate limit exceeded") {
		t.Errorf("error = %q, want to contain %q", err.Error(), "rate limit exceeded")
	}
}

func TestAnalyze_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := newTestClient(t, server)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := c.Analyze(ctx, domain.AnalysisRequest{Prompt: "test"})

	if err == nil {
		t.Fatal("Analyze should return error on timeout")
	}
	if !strings.Contains(err.Error(), "failed to send request") {
		t.Errorf("error = %q, want to contain %q", err.Error(), "failed to send request")
	}
}

func TestAnalyze_InvalidJSON_Response(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{invalid json}`)
	}))
	defer server.Close()

	c := newTestClient(t, server)
	_, err := c.Analyze(context.Background(), domain.AnalysisRequest{Prompt: "test"})

	if err == nil {
		t.Fatal("Analyze should return error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "failed to parse response") {
		t.Errorf("error = %q, want to contain %q", err.Error(), "failed to parse response")
	}
}

func TestAnalyze_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"candidates":[]}`)
	}))
	defer server.Close()

	c := newTestClient(t, server)
	_, err := c.Analyze(context.Background(), domain.AnalysisRequest{Prompt: "test"})

	if err == nil {
		t.Fatal("Analyze should return error for empty response")
	}
	if !strings.Contains(err.Error(), "empty response") {
		t.Errorf("error = %q, want to contain %q", err.Error(), "empty response")
	}
}

func TestAnalyze_ModelOverride(t *testing.T) {
	var receivedURL string
	var receivedAPIKey string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURL = r.URL.String()
		receivedAPIKey = r.Header.Get("x-goog-api-key")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"candidates":[{"content":{"parts":[{"text":"ok"}]}}]}`)
	}))
	defer server.Close()

	c := newTestClient(t, server)
	_, err := c.Analyze(context.Background(), domain.AnalysisRequest{
		Prompt: "test",
		Model:  "gemini-1.5-pro",
	})

	if err != nil {
		t.Fatalf("Analyze returned unexpected error: %v", err)
	}
	if !strings.Contains(receivedURL, "gemini-1.5-pro") {
		t.Errorf("URL = %q, want to contain %q", receivedURL, "gemini-1.5-pro")
	}
	if receivedAPIKey != "test-api-key" {
		t.Errorf("x-goog-api-key = %q, want %q", receivedAPIKey, "test-api-key")
	}
	if strings.Contains(receivedURL, "key=") {
		t.Errorf("URL should not contain API key in query params: %q", receivedURL)
	}
}

func TestAnalyze_RequestFormat(t *testing.T) {
	var receivedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		receivedBody = string(body)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"candidates":[{"content":{"parts":[{"text":"ok"}]}}]}`)
	}))
	defer server.Close()

	c := newTestClient(t, server)
	_, err := c.Analyze(context.Background(), domain.AnalysisRequest{
		Prompt: "analyze this",
		Data:   "some data",
	})

	if err != nil {
		t.Fatalf("Analyze returned unexpected error: %v", err)
	}
	if !strings.Contains(receivedBody, "analyze this") {
		t.Errorf("request body = %q, want to contain %q", receivedBody, "analyze this")
	}
	if !strings.Contains(receivedBody, "some data") {
		t.Errorf("request body = %q, want to contain %q", receivedBody, "some data")
	}
}
