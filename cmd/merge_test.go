package cmd

import (
	"testing"

	"github.com/canpok1/github-analyzer/internal/domain/entity"
	"github.com/canpok1/github-analyzer/internal/infra/gemini"
)

// TODO: CLIフラグが設定されている場合、設定ファイルよりCLIフラグを優先
// TODO: CLIフラグが未設定で設定ファイルに値がある場合、設定ファイルの値を使用
// TODO: 両方未設定の場合、デフォルト値を使用
// TODO: 全フィールド（repo, prompt, model）をそれぞれテスト

func TestApplyConfig_CLIFlagOverridesConfig(t *testing.T) {
	query := entity.Query{
		Repo:   "cli/repo",
		Prompt: "cli prompt",
	}
	cfg := entity.Config{
		Repo:          "config/repo",
		DefaultPrompt: "config prompt",
		Model:         "gemini-2.0-flash",
	}

	q, model := applyConfig(query, cfg)

	if q.Repo != "cli/repo" {
		t.Errorf("Repo = %q, want %q", q.Repo, "cli/repo")
	}
	if q.Prompt != "cli prompt" {
		t.Errorf("Prompt = %q, want %q", q.Prompt, "cli prompt")
	}
	if model != "gemini-2.0-flash" {
		t.Errorf("Model = %q, want %q", model, "gemini-2.0-flash")
	}
}

func TestApplyConfig_FallbackToConfig(t *testing.T) {
	query := entity.Query{
		Repo:   "",
		Prompt: "",
	}
	cfg := entity.Config{
		Repo:          "config/repo",
		DefaultPrompt: "config prompt",
		Model:         "gemini-2.0-flash",
	}

	q, model := applyConfig(query, cfg)

	if q.Repo != "config/repo" {
		t.Errorf("Repo = %q, want %q", q.Repo, "config/repo")
	}
	if q.Prompt != "config prompt" {
		t.Errorf("Prompt = %q, want %q", q.Prompt, "config prompt")
	}
	if model != "gemini-2.0-flash" {
		t.Errorf("Model = %q, want %q", model, "gemini-2.0-flash")
	}
}

func TestApplyConfig_DefaultValues(t *testing.T) {
	query := entity.Query{
		Repo:   "",
		Prompt: "",
	}
	cfg := entity.Config{}

	q, model := applyConfig(query, cfg)

	if q.Repo != "" {
		t.Errorf("Repo = %q, want empty", q.Repo)
	}
	if q.Prompt != "" {
		t.Errorf("Prompt = %q, want empty", q.Prompt)
	}
	if model != gemini.DefaultModel {
		t.Errorf("Model = %q, want %q", model, gemini.DefaultModel)
	}
}
