package cmd

import (
	"github.com/canpok1/github-analyzer/internal/domain/entity"
)

// applyConfig は設定ファイルの値をQueryに反映する。
// CLIフラグ（空でない値） > 設定ファイル > デフォルト値 の優先順位。
// 戻り値はマージ後のQueryと設定ファイルで指定されたモデル名（未指定時は空文字列）。
func applyConfig(query entity.Query, cfg entity.Config) (entity.Query, string) {
	if query.Repo == "" && cfg.Repo != "" {
		query.Repo = cfg.Repo
	}

	if query.Prompt == "" && cfg.DefaultPrompt != "" {
		query.Prompt = cfg.DefaultPrompt
	}

	return query, cfg.Model
}
