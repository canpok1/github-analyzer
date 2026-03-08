package cmd

import (
	"testing"

	"github.com/canpok1/github-analyzer/internal/domain/entity"
)

// DONE: CLIフラグのモデルが設定されている場合、設定ファイルよりCLIフラグを優先
// DONE: CLIフラグのモデルが未設定で設定ファイルに値がある場合、設定ファイルの値を使用
// DONE: モデルが両方未設定の場合、空文字列（デフォルト値）を使用
// DONE: applyConfigの戻り値からmodel stringを除去しQuery.Modelに統合

func TestApplyConfig_CLIFlagOverridesConfig(t *testing.T) {
	query := entity.Query{
		Repo:   "cli/repo",
		Prompt: "cli prompt",
		Model:  "cli-model",
	}
	cfg := entity.Config{
		Repo:          "config/repo",
		DefaultPrompt: "config prompt",
		Model:         "config-model",
	}

	q := applyConfig(query, cfg)

	if q.Repo != "cli/repo" {
		t.Errorf("Repo = %q, want %q", q.Repo, "cli/repo")
	}
	if q.Prompt != "cli prompt" {
		t.Errorf("Prompt = %q, want %q", q.Prompt, "cli prompt")
	}
	if q.Model != "cli-model" {
		t.Errorf("Model = %q, want %q", q.Model, "cli-model")
	}
}

func TestApplyConfig_FallbackToConfig(t *testing.T) {
	query := entity.Query{
		Repo:   "",
		Prompt: "",
		Model:  "",
	}
	cfg := entity.Config{
		Repo:          "config/repo",
		DefaultPrompt: "config prompt",
		Model:         "config-model",
	}

	q := applyConfig(query, cfg)

	if q.Repo != "config/repo" {
		t.Errorf("Repo = %q, want %q", q.Repo, "config/repo")
	}
	if q.Prompt != "config prompt" {
		t.Errorf("Prompt = %q, want %q", q.Prompt, "config prompt")
	}
	if q.Model != "config-model" {
		t.Errorf("Model = %q, want %q", q.Model, "config-model")
	}
}

func TestApplyConfig_DefaultValues(t *testing.T) {
	query := entity.Query{
		Repo:   "",
		Prompt: "",
		Model:  "",
	}
	cfg := entity.Config{}

	q := applyConfig(query, cfg)

	if q.Repo != "" {
		t.Errorf("Repo = %q, want empty", q.Repo)
	}
	if q.Prompt != "" {
		t.Errorf("Prompt = %q, want empty", q.Prompt)
	}
	if q.Model != "" {
		t.Errorf("Model = %q, want empty", q.Model)
	}
}
