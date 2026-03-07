package domain

import "context"

// AnalysisRequest はAI分析リクエストを表す。
type AnalysisRequest struct {
	Prompt string
	Data   string
	Model  string
}

// AnalysisResponse はAI分析レスポンスを表す。
type AnalysisResponse struct {
	Content string
}

// Analyzer はAI分析サービスを抽象化するインターフェース。
type Analyzer interface {
	// Analyze はプロンプトとデータを元にAI分析を実行する。
	Analyze(ctx context.Context, req AnalysisRequest) (*AnalysisResponse, error)
}
