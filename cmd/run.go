package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/canpok1/github-analyzer/internal/domain/entity"
	"github.com/spf13/cobra"
)

// resolveToken はGH_TOKEN → GITHUB_TOKENの優先順でトークンを解決する。
func resolveToken() (string, error) {
	if token := os.Getenv("GH_TOKEN"); token != "" {
		return token, nil
	}
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token, nil
	}
	return "", fmt.Errorf("GH_TOKEN または GITHUB_TOKEN 環境変数を設定してください")
}

// buildQuery はcobraコマンドのフラグからQueryを構築する。
func buildQuery(cmd *cobra.Command) (entity.Query, error) {
	query := entity.Query{}

	repo, _ := cmd.Flags().GetString("repo")
	query.Repo = repo

	prompt, _ := cmd.Flags().GetString("prompt")
	query.Prompt = prompt

	status, _ := cmd.Flags().GetString("status")
	query.Status = status

	today, _ := cmd.Flags().GetBool("today")
	sinceStr, _ := cmd.Flags().GetString("since")
	pr, _ := cmd.Flags().GetInt("pr")
	issue, _ := cmd.Flags().GetInt("issue")

	if today {
		now := time.Now()
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		query.Since = &startOfDay
	}

	if sinceStr != "" {
		dur, err := entity.ParseDuration(sinceStr)
		if err != nil {
			return entity.Query{}, fmt.Errorf("--since の値が不正です: %w", err)
		}
		since := time.Now().Add(-dur)
		query.Since = &since
	}

	if pr != 0 {
		query.PR = &pr
	}

	if issue != 0 {
		query.Issue = &issue
	}

	return query, nil
}
