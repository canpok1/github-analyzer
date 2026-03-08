package mock

import (
	"context"

	"github.com/canpok1/github-analyzer/internal/domain"
)

const dummyReport = `## Overview
[Mock] チーム全体として安定した開発ペースを維持しており、活発な議論と着実なマージが行われています。

## Process Insights
モックモードで実行されているため、実際のデータに基づく分析は行われていません。実データでの分析を行うには、設定ファイルの mock.ai を false に変更してください。

## Potential Risks
モックモードのため、実際のリスク分析は行われていません。

## Manager's Hint
モックモードでの動作確認が完了したら、実際のAPIキーを設定して本番の分析を実行してください。`

// Analyzer は domain.Analyzer のモック実装。
// 固定のダミーレポートを返す。
type Analyzer struct{}

func (a *Analyzer) Analyze(_ context.Context, _ domain.AnalysisRequest) (*domain.AnalysisResponse, error) {
	return &domain.AnalysisResponse{
		Content: dummyReport,
	}, nil
}
