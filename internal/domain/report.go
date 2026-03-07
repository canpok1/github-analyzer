package domain

import "github.com/canpok1/github-analyzer/internal/domain/entity"

// ReportRenderer はReportエンティティをフォーマット済み文字列に変換するインターフェース。
type ReportRenderer interface {
	// Render はReportをフォーマット済み文字列に変換する。
	Render(report *entity.Report) (string, error)
}

// ReportWriter はレンダリング済みレポートを出力先に書き込むインターフェース。
type ReportWriter interface {
	// Write はレンダリング済みレポートを出力先に書き込む。
	Write(content string) error
}
