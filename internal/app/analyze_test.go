package app

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/canpok1/github-analyzer/internal/domain"
	"github.com/canpok1/github-analyzer/internal/domain/entity"
)

// DONE: 正常系: 進捗メッセージがstderrに出力される
// DONE: 異常系: CollectDataでエラーが発生した場合、エラーを返す
// DONE: 異常系: Analyzerでエラーが発生した場合、エラーを返す
// DONE: 異常系: ReportParserでエラーが発生した場合、エラーを返す
// DONE: 異常系: Rendererでエラーが発生した場合、エラーを返す
// DONE: 異常系: Writerでエラーが発生した場合、エラーを返す

// analyzeTestAnalyzer はテスト用のAnalyzer実装。
type analyzeTestAnalyzer struct {
	resp *domain.AnalysisResponse
	err  error
}

func (m *analyzeTestAnalyzer) Analyze(_ context.Context, _ domain.AnalysisRequest) (*domain.AnalysisResponse, error) {
	return m.resp, m.err
}

// analyzeTestRenderer はテスト用のReportRenderer実装。
type analyzeTestRenderer struct {
	output string
	err    error
}

func (m *analyzeTestRenderer) Render(_ *entity.Report) (string, error) {
	return m.output, m.err
}

// analyzeTestWriter はテスト用のReportWriter実装。
type analyzeTestWriter struct {
	written string
	err     error
}

func (m *analyzeTestWriter) Write(content string) error {
	m.written = content
	return m.err
}

func TestAnalyze_Success(t *testing.T) {
	pr1 := 1
	query := entity.Query{
		Repo: "owner/repo",
		PR:   &pr1,
	}

	gh := &mockGitHubRepository{}
	analyzer := &analyzeTestAnalyzer{
		resp: &domain.AnalysisResponse{Content: "test content"},
	}
	renderer := &analyzeTestRenderer{output: "rendered report"}
	writer := &analyzeTestWriter{}
	stderr := &bytes.Buffer{}

	deps := AnalyzeDeps{
		GitHubRepo: gh,
		Analyzer:   analyzer,
		PromptBuilder: func(data *CollectedData, userPrompt string) domain.AnalysisRequest {
			return domain.AnalysisRequest{Prompt: "test", Data: "data"}
		},
		ReportParser: func(content string) (*entity.Report, error) {
			return &entity.Report{Overview: "overview"}, nil
		},
		Renderer: renderer,
		Writer:   writer,
		Stderr:   stderr,
	}

	err := Analyze(context.Background(), deps, query)
	if err != nil {
		t.Fatalf("Analyze() returned error: %v", err)
	}

	if writer.written != "rendered report" {
		t.Errorf("expected writer to receive 'rendered report', got %q", writer.written)
	}

	// 進捗メッセージがstderrに出力されていること
	stderrOutput := stderr.String()
	expectedMessages := []string{"データを収集中", "AI分析を実行中", "レポートを生成中", "完了"}
	for _, msg := range expectedMessages {
		if !strings.Contains(stderrOutput, msg) {
			t.Errorf("stderr should contain %q, got %q", msg, stderrOutput)
		}
	}
}

// newSuccessDeps はテスト用の正常系依存を生成するヘルパー。
func newSuccessDeps() (AnalyzeDeps, *analyzeTestWriter) {
	writer := &analyzeTestWriter{}
	deps := AnalyzeDeps{
		GitHubRepo: &mockGitHubRepository{},
		Analyzer: &analyzeTestAnalyzer{
			resp: &domain.AnalysisResponse{Content: "content"},
		},
		PromptBuilder: func(_ *CollectedData, _ string) domain.AnalysisRequest {
			return domain.AnalysisRequest{Prompt: "p", Data: "d"}
		},
		ReportParser: func(_ string) (*entity.Report, error) {
			return &entity.Report{Overview: "o"}, nil
		},
		Renderer: &analyzeTestRenderer{output: "rendered"},
		Writer:   writer,
		Stderr:   &bytes.Buffer{},
	}
	return deps, writer
}

func TestAnalyze_CollectDataError(t *testing.T) {
	deps, _ := newSuccessDeps()
	// 不正なrepo形式でCollectDataをエラーにする
	query := entity.Query{
		Repo: "invalid",
		PR:   intPtr(1),
	}

	err := Analyze(context.Background(), deps, query)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "データ収集に失敗しました") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestAnalyze_AnalyzerError(t *testing.T) {
	deps, _ := newSuccessDeps()
	deps.Analyzer = &analyzeTestAnalyzer{err: fmt.Errorf("API unavailable")}
	query := entity.Query{
		Repo: "owner/repo",
		PR:   intPtr(1),
	}

	err := Analyze(context.Background(), deps, query)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "AI分析に失敗しました") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestAnalyze_ReportParserError(t *testing.T) {
	deps, _ := newSuccessDeps()
	deps.ReportParser = func(_ string) (*entity.Report, error) {
		return nil, fmt.Errorf("parse error")
	}
	query := entity.Query{
		Repo: "owner/repo",
		PR:   intPtr(1),
	}

	err := Analyze(context.Background(), deps, query)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "レポートのパースに失敗しました") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestAnalyze_RendererError(t *testing.T) {
	deps, _ := newSuccessDeps()
	deps.Renderer = &analyzeTestRenderer{err: fmt.Errorf("render error")}
	query := entity.Query{
		Repo: "owner/repo",
		PR:   intPtr(1),
	}

	err := Analyze(context.Background(), deps, query)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "レポートのレンダリングに失敗しました") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestAnalyze_WriterError(t *testing.T) {
	deps, _ := newSuccessDeps()
	deps.Writer = &analyzeTestWriter{err: fmt.Errorf("write error")}
	query := entity.Query{
		Repo: "owner/repo",
		PR:   intPtr(1),
	}

	err := Analyze(context.Background(), deps, query)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "レポートの出力に失敗しました") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestAnalyze_LoggerReceivesPromptAndData(t *testing.T) {
	deps, _ := newSuccessDeps()
	var logged []string
	deps.Logger = func(msg string) {
		logged = append(logged, msg)
	}
	deps.PromptBuilder = func(_ *CollectedData, _ string) domain.AnalysisRequest {
		return domain.AnalysisRequest{Prompt: "test-prompt", Data: "test-data"}
	}
	query := entity.Query{
		Repo: "owner/repo",
		PR:   intPtr(1),
	}

	err := Analyze(context.Background(), deps, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(logged) < 2 {
		t.Fatalf("expected at least 2 log entries, got %d", len(logged))
	}
	if !strings.Contains(logged[0], "test-prompt") {
		t.Errorf("first log should contain prompt, got %q", logged[0])
	}
	if !strings.Contains(logged[1], "test-data") {
		t.Errorf("second log should contain data, got %q", logged[1])
	}
}

func TestAnalyze_NilLoggerDoesNotPanic(t *testing.T) {
	deps, _ := newSuccessDeps()
	deps.Logger = nil
	query := entity.Query{
		Repo: "owner/repo",
		PR:   intPtr(1),
	}

	err := Analyze(context.Background(), deps, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func intPtr(v int) *int {
	return &v
}
