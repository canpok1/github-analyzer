package report_test

import (
	"strings"
	"testing"

	"github.com/canpok1/github-analyzer/internal/domain/entity"
	"github.com/canpok1/github-analyzer/internal/infra/report"
)

// TODO: 正常系: 各セクションに改行を含むテキストが正しくフォーマットされる
// TODO: 境界値: 全セクションが空文字列のReportでもエラーにならない

func TestMarkdownRenderer_Render(t *testing.T) {
	r := report.NewMarkdownRenderer()

	input := &entity.Report{
		Overview:        "活発な開発が行われています",
		ProcessInsights: "レビュープロセスが効率的に機能しています",
		PotentialRisks:  "特に大きなリスクは見られません",
		ManagersHint:    "現在の体制を維持してください",
	}

	got, err := r.Render(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 各セクション見出しが含まれていることを確認
	expectedHeaders := []string{
		"## Overview",
		"## Process Insights",
		"## Potential Risks",
		"## Manager's Hint",
	}
	for _, header := range expectedHeaders {
		if !strings.Contains(got, header) {
			t.Errorf("expected header %q not found in output:\n%s", header, got)
		}
	}

	// 各セクションの本文が含まれていることを確認
	expectedBodies := []string{
		input.Overview,
		input.ProcessInsights,
		input.PotentialRisks,
		input.ManagersHint,
	}
	for _, body := range expectedBodies {
		if !strings.Contains(got, body) {
			t.Errorf("expected body %q not found in output:\n%s", body, got)
		}
	}
}
