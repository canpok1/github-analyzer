package gemini

import (
	"context"
	"fmt"
	"strings"

	"github.com/canpok1/github-analyzer/internal/app"
	"github.com/canpok1/github-analyzer/internal/domain"
	"github.com/canpok1/github-analyzer/internal/domain/entity"
)

// AnalyzeReport はCollectedDataとユーザープロンプトを受け取り、
// BuildPrompt → Analyze → ParseReport の流れでレポートを生成する。
func AnalyzeReport(ctx context.Context, analyzer domain.Analyzer, data *app.CollectedData, userPrompt string) (*entity.Report, error) {
	if data == nil {
		return nil, fmt.Errorf("data must not be nil")
	}
	if analyzer == nil {
		return nil, fmt.Errorf("analyzer must not be nil")
	}

	req := BuildPrompt(data, userPrompt)

	resp, err := analyzer.Analyze(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	report, err := ParseReport(resp.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse report: %w", err)
	}

	return report, nil
}

// sectionHeaders はレポートの各セクション見出し。
var sectionHeaders = []string{
	"## Overview",
	"## Process Insights",
	"## Potential Risks",
	"## Manager's Hint",
}

// ParseReport はGemini APIの生テキスト出力をReportエンティティにパースする。
func ParseReport(content string) (*entity.Report, error) {
	if content == "" {
		return nil, fmt.Errorf("empty content")
	}

	// 全セクションの開始位置を事前計算
	indices := make([]int, len(sectionHeaders))
	for i, header := range sectionHeaders {
		idx := strings.Index(content, header)
		if idx == -1 {
			return nil, fmt.Errorf("missing section: %s", header)
		}
		indices[i] = idx + len(header)
	}

	// 各セクションの本文を抽出
	bodies := make([]string, len(sectionHeaders))
	for i := range sectionHeaders {
		start := indices[i]
		end := len(content)
		if i+1 < len(indices) {
			end = indices[i+1] - len(sectionHeaders[i+1])
		}
		bodies[i] = strings.TrimSpace(content[start:end])
	}

	return &entity.Report{
		Overview:        bodies[0],
		ProcessInsights: bodies[1],
		PotentialRisks:  bodies[2],
		ManagersHint:    bodies[3],
	}, nil
}
