package report_test

import (
	"strings"
	"testing"

	"github.com/canpok1/github-analyzer/internal/domain"
	"github.com/canpok1/github-analyzer/internal/domain/entity"
	"github.com/canpok1/github-analyzer/internal/infra/report"
)

// MarkdownRendererがReportRendererインターフェースを満たすことを確認
var _ domain.ReportRenderer = (*report.MarkdownRenderer)(nil)

func TestMarkdownRenderer_Render_MultilineContent(t *testing.T) {
	r := report.NewMarkdownRenderer()

	input := &entity.Report{
		Overview:        "1行目\n2行目\n3行目",
		ProcessInsights: "- ポイントA\n- ポイントB",
		PotentialRisks:  "リスク1\n\nリスク2",
		ManagersHint:    "アクション1\nアクション2",
	}

	got, err := r.Render(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 改行を含む本文がそのまま含まれること
	if !strings.Contains(got, "1行目\n2行目\n3行目") {
		t.Errorf("multiline overview not preserved in output:\n%s", got)
	}
	if !strings.Contains(got, "- ポイントA\n- ポイントB") {
		t.Errorf("multiline process insights not preserved in output:\n%s", got)
	}
}

func TestMarkdownRenderer_Render_NilReport(t *testing.T) {
	r := report.NewMarkdownRenderer()

	_, err := r.Render(nil)
	if err == nil {
		t.Fatal("expected error for nil report, got nil")
	}
}

func TestMarkdownRenderer_Render_EmptySections(t *testing.T) {
	r := report.NewMarkdownRenderer()

	input := &entity.Report{}

	got, err := r.Render(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 見出しが含まれていることを確認
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
}

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
