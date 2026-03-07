package gemini

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/canpok1/github-analyzer/internal/app"
	"github.com/canpok1/github-analyzer/internal/domain"
	"github.com/canpok1/github-analyzer/internal/domain/entity"
)

// stubAnalyzer はテスト用のAnalyzerスタブ。
type stubAnalyzer struct {
	response *domain.AnalysisResponse
	err      error
}

func (s *stubAnalyzer) Analyze(_ context.Context, _ domain.AnalysisRequest) (*domain.AnalysisResponse, error) {
	return s.response, s.err
}

// validReportContent はテスト用の有効なレポートコンテンツ。
const validReportContent = `## Overview
活発な開発が行われています。

## Process Insights
レビュープロセスが効率的です。

## Potential Risks
特定メンバーへの負荷集中が見られます。

## Manager's Hint
1on1でメンバーの状況を確認しましょう。`

func newEmptyCollectedData() *app.CollectedData {
	return &app.CollectedData{
		PullRequests: []entity.PullRequest{},
		Issues:       []entity.Issue{},
		Comments:     make(map[int][]entity.Comment),
		Timeline:     make(map[int][]entity.TimelineEvent),
	}
}

func TestAnalyzeReport_Success(t *testing.T) {
	t.Parallel()

	analyzer := &stubAnalyzer{
		response: &domain.AnalysisResponse{Content: validReportContent},
	}
	data := newEmptyCollectedData()

	report, err := AnalyzeReport(context.Background(), analyzer, data, "チームの分析をお願いします")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if report.Overview != "活発な開発が行われています。" {
		t.Errorf("Overview = %q, want %q", report.Overview, "活発な開発が行われています。")
	}
	if report.ProcessInsights != "レビュープロセスが効率的です。" {
		t.Errorf("ProcessInsights = %q, want %q", report.ProcessInsights, "レビュープロセスが効率的です。")
	}
	if report.PotentialRisks != "特定メンバーへの負荷集中が見られます。" {
		t.Errorf("PotentialRisks = %q, want %q", report.PotentialRisks, "特定メンバーへの負荷集中が見られます。")
	}
	if report.ManagersHint != "1on1でメンバーの状況を確認しましょう。" {
		t.Errorf("ManagersHint = %q, want %q", report.ManagersHint, "1on1でメンバーの状況を確認しましょう。")
	}
}

func TestAnalyzeReport_EmptyUserPrompt(t *testing.T) {
	t.Parallel()

	analyzer := &stubAnalyzer{
		response: &domain.AnalysisResponse{Content: validReportContent},
	}
	data := newEmptyCollectedData()

	report, err := AnalyzeReport(context.Background(), analyzer, data, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if report.Overview != "活発な開発が行われています。" {
		t.Errorf("Overview = %q, want %q", report.Overview, "活発な開発が行われています。")
	}
}

func TestAnalyzeReport_AnalyzerError(t *testing.T) {
	t.Parallel()

	analyzer := &stubAnalyzer{
		err: fmt.Errorf("API connection failed"),
	}
	data := newEmptyCollectedData()

	_, err := AnalyzeReport(context.Background(), analyzer, data, "分析してください")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "analysis failed") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "analysis failed")
	}
	if !strings.Contains(err.Error(), "API connection failed") {
		t.Errorf("error = %q, want it to contain original error", err.Error())
	}
}

func TestAnalyzeReport_ParseError(t *testing.T) {
	t.Parallel()

	analyzer := &stubAnalyzer{
		response: &domain.AnalysisResponse{Content: "invalid content without sections"},
	}
	data := newEmptyCollectedData()

	_, err := AnalyzeReport(context.Background(), analyzer, data, "分析してください")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to parse report") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "failed to parse report")
	}
}

func TestAnalyzeReport_NilData(t *testing.T) {
	t.Parallel()

	analyzer := &stubAnalyzer{
		response: &domain.AnalysisResponse{Content: validReportContent},
	}

	_, err := AnalyzeReport(context.Background(), analyzer, nil, "分析してください")
	if err == nil {
		t.Fatal("expected error for nil data, got nil")
	}

	if !strings.Contains(err.Error(), "data") {
		t.Errorf("error = %q, want it to mention data", err.Error())
	}
}

func TestAnalyzeReport_NilAnalyzer(t *testing.T) {
	t.Parallel()

	data := newEmptyCollectedData()

	_, err := AnalyzeReport(context.Background(), nil, data, "分析してください")
	if err == nil {
		t.Fatal("expected error for nil analyzer, got nil")
	}

	if !strings.Contains(err.Error(), "analyzer") {
		t.Errorf("error = %q, want it to mention analyzer", err.Error())
	}
}

// AnalyzeReport テストリスト
// DONE: 正常系: BuildPrompt → Analyze → ParseReport の流れで正しい Report を返す
// DONE: 正常系: userPromptが空の場合でもデフォルトプロンプトで動作する
// DONE: 異常系: Analyzerがエラーを返した場合にエラーを返す
// DONE: 異常系: Analyzerの応答がパースできない場合にエラーを返す
// DONE: 異常系: dataがnilの場合にエラーを返す
// DONE: 異常系: analyzerがnilの場合にエラーを返す

func TestParseReport_AllSections(t *testing.T) {
	t.Parallel()

	content := `## Overview
活発な開発が行われています。

## Process Insights
レビュープロセスが効率的です。

## Potential Risks
特定メンバーへの負荷集中が見られます。

## Manager's Hint
1on1でメンバーの状況を確認しましょう。`

	report, err := ParseReport(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if report.Overview != "活発な開発が行われています。" {
		t.Errorf("Overview = %q, want %q", report.Overview, "活発な開発が行われています。")
	}
	if report.ProcessInsights != "レビュープロセスが効率的です。" {
		t.Errorf("ProcessInsights = %q, want %q", report.ProcessInsights, "レビュープロセスが効率的です。")
	}
	if report.PotentialRisks != "特定メンバーへの負荷集中が見られます。" {
		t.Errorf("PotentialRisks = %q, want %q", report.PotentialRisks, "特定メンバーへの負荷集中が見られます。")
	}
	if report.ManagersHint != "1on1でメンバーの状況を確認しましょう。" {
		t.Errorf("ManagersHint = %q, want %q", report.ManagersHint, "1on1でメンバーの状況を確認しましょう。")
	}
}

func TestParseReport_MultilineContent(t *testing.T) {
	t.Parallel()

	content := `## Overview
活発な開発が行われています。

## Process Insights
レビュープロセスが効率的です。
平均レビュー時間は2時間以内です。
チーム全体でコードレビューに参加しています。

## Potential Risks
特定メンバーへの負荷集中が見られます。

## Manager's Hint
1on1でメンバーの状況を確認しましょう。`

	report, err := ParseReport(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "レビュープロセスが効率的です。\n平均レビュー時間は2時間以内です。\nチーム全体でコードレビューに参加しています。"
	if report.ProcessInsights != want {
		t.Errorf("ProcessInsights = %q, want %q", report.ProcessInsights, want)
	}
}

func TestParseReport_EmptyContent(t *testing.T) {
	t.Parallel()

	_, err := ParseReport("")
	if err == nil {
		t.Fatal("expected error for empty content, got nil")
	}

	want := "empty content"
	if err.Error() != want {
		t.Errorf("error = %q, want %q", err.Error(), want)
	}
}

func TestParseReport_MissingSection(t *testing.T) {
	t.Parallel()

	content := `## Overview
活発な開発が行われています。

## Process Insights
レビュープロセスが効率的です。

## Manager's Hint
1on1でメンバーの状況を確認しましょう。`

	_, err := ParseReport(content)
	if err == nil {
		t.Fatal("expected error for missing section, got nil")
	}

	if !strings.Contains(err.Error(), "Potential Risks") {
		t.Errorf("error = %q, want it to mention missing section", err.Error())
	}
}

func TestParseReport_ExtraWhitespace(t *testing.T) {
	t.Parallel()

	content := `
前置きテキスト

## Overview

活発な開発が行われています。


## Process Insights

レビュープロセスが効率的です。

## Potential Risks

特定メンバーへの負荷集中が見られます。

## Manager's Hint

1on1でメンバーの状況を確認しましょう。

`

	report, err := ParseReport(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if report.Overview != "活発な開発が行われています。" {
		t.Errorf("Overview = %q, want %q", report.Overview, "活発な開発が行われています。")
	}
	if report.ManagersHint != "1on1でメンバーの状況を確認しましょう。" {
		t.Errorf("ManagersHint = %q, want %q", report.ManagersHint, "1on1でメンバーの状況を確認しましょう。")
	}
}
