package cmd

import (
	"github.com/canpok1/github-analyzer/internal/domain/entity"
	"github.com/canpok1/github-analyzer/internal/infra/gemini"
)

// applyConfig は設定ファイルの値をQueryに反映する。
// CLIフラグ（空でない値） > 設定ファイル > デフォルト値 の優先順位。
// 戻り値はマージ後のQueryと使用するモデル名。
func applyConfig(query entity.Query, cfg entity.Config) (entity.Query, string) {
	if query.Repo == "" && cfg.Repo != "" {
		query.Repo = cfg.Repo
	}

	if query.Prompt == "" && cfg.DefaultPrompt != "" {
		query.Prompt = cfg.DefaultPrompt
	}

	model := gemini.DefaultModel
	if cfg.Model != "" {
		model = cfg.Model
	}

	return query, model
}
