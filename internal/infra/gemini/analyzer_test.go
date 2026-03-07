package gemini

import (
	"testing"
)

// TODO: 正常系: セクション内に複数行のテキストがある場合
// TODO: 異常系: 空文字列の場合エラーを返す
// TODO: 異常系: 必須セクション(Overview)が欠落している場合エラーを返す
// TODO: 境界値: セクション見出しの前後に余分な空白・改行がある場合

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
