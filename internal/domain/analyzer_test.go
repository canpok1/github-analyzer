package domain

import (
	"context"
	"testing"
)

// TODO: 正常系: Analyzerインターフェースをモック実装で満たせること
// TODO: AnalysisRequestの構造体フィールド確認
// TODO: AnalysisResponseの構造体フィールド確認

// mockAnalyzer はテスト用のモック実装。
type mockAnalyzer struct{}

func (m *mockAnalyzer) Analyze(ctx context.Context, req AnalysisRequest) (*AnalysisResponse, error) {
	return &AnalysisResponse{
		Content: "test analysis result",
	}, nil
}

func TestAnalyzer_InterfaceImplementation(t *testing.T) {
	var _ Analyzer = &mockAnalyzer{}
}

func TestAnalysisRequest_HasExpectedFields(t *testing.T) {
	req := AnalysisRequest{
		Prompt: "analyze this",
		Data:   "some data",
		Model:  "gemini-1.5-pro",
	}

	if req.Prompt != "analyze this" {
		t.Errorf("Prompt = %q, want %q", req.Prompt, "analyze this")
	}
	if req.Data != "some data" {
		t.Errorf("Data = %q, want %q", req.Data, "some data")
	}
	if req.Model != "gemini-1.5-pro" {
		t.Errorf("Model = %q, want %q", req.Model, "gemini-1.5-pro")
	}
}

func TestAnalysisResponse_HasExpectedFields(t *testing.T) {
	resp := AnalysisResponse{
		Content: "analysis result",
	}

	if resp.Content != "analysis result" {
		t.Errorf("Content = %q, want %q", resp.Content, "analysis result")
	}
}
