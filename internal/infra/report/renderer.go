package report

import (
	"fmt"

	"github.com/canpok1/github-analyzer/internal/domain/entity"
)

// MarkdownRenderer はReportをMarkdown形式にレンダリングする。
type MarkdownRenderer struct{}

// NewMarkdownRenderer はMarkdownRendererを生成する。
func NewMarkdownRenderer() *MarkdownRenderer {
	return &MarkdownRenderer{}
}

// Render はReportをMarkdown文字列に変換する。
func (r *MarkdownRenderer) Render(report *entity.Report) (string, error) {
	if report == nil {
		return "", fmt.Errorf("report must not be nil")
	}

	return fmt.Sprintf(`## Overview

%s

## Process Insights

%s

## Potential Risks

%s

## Manager's Hint

%s
`, report.Overview, report.ProcessInsights, report.PotentialRisks, report.ManagersHint), nil
}
