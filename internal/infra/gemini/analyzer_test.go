package gemini

import (
	"strings"
	"testing"
)

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
