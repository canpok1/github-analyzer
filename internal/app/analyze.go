package app

import (
	"context"
	"fmt"
	"io"

	"github.com/canpok1/github-analyzer/internal/domain"
	"github.com/canpok1/github-analyzer/internal/domain/entity"
)

// PromptBuilderFunc はCollectedDataとユーザープロンプトからAnalysisRequestを構築する関数型。
type PromptBuilderFunc func(data *CollectedData, userPrompt string) domain.AnalysisRequest

// ReportParserFunc はAI分析結果のテキストをReportエンティティにパースする関数型。
type ReportParserFunc func(content string) (*entity.Report, error)

// AnalyzeDeps はAnalyze関数の依存をまとめた構造体。
type AnalyzeDeps struct {
	GitHubRepo    domain.GitHubRepository
	Analyzer      domain.Analyzer
	PromptBuilder PromptBuilderFunc
	ReportParser  ReportParserFunc
	Renderer      domain.ReportRenderer
	Writer        domain.ReportWriter
	Stderr        io.Writer
}

// printProgress はstderrに進捗メッセージを出力する。
func printProgress(w io.Writer, msg string) {
	_, _ = fmt.Fprintf(w, "%s\n", msg)
}

// Analyze はデータ収集→プロンプト構築→AI分析→レポート生成→出力の一連フローを実行する。
func Analyze(ctx context.Context, deps AnalyzeDeps, query entity.Query) error {
	printProgress(deps.Stderr, "データを収集中...")

	data, err := CollectData(ctx, deps.GitHubRepo, query)
	if err != nil {
		return fmt.Errorf("データ収集に失敗しました: %w", err)
	}

	printProgress(deps.Stderr, "AI分析を実行中...")

	req := deps.PromptBuilder(data, query.Prompt)
	resp, err := deps.Analyzer.Analyze(ctx, req)
	if err != nil {
		return fmt.Errorf("AI分析に失敗しました: %w", err)
	}

	printProgress(deps.Stderr, "レポートを生成中...")

	report, err := deps.ReportParser(resp.Content)
	if err != nil {
		return fmt.Errorf("レポートのパースに失敗しました: %w", err)
	}

	rendered, err := deps.Renderer.Render(report)
	if err != nil {
		return fmt.Errorf("レポートのレンダリングに失敗しました: %w", err)
	}

	if err := deps.Writer.Write(rendered); err != nil {
		return fmt.Errorf("レポートの出力に失敗しました: %w", err)
	}

	printProgress(deps.Stderr, "完了")
	return nil
}
