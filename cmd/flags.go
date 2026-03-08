package cmd

import (
	"fmt"

	"github.com/canpok1/github-analyzer/internal/domain/entity"
	"github.com/spf13/cobra"
)

func defineFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("today", false, "本日の活動を分析")
	cmd.Flags().String("since", "", "指定期間の活動を分析（例: 7d, 2w, 1m）")
	cmd.Flags().Int("pr", 0, "特定のPR番号を指定")
	cmd.Flags().Int("issue", 0, "特定のIssue番号を指定")
	cmd.Flags().String("status", "", "ステータスでフィルタ（open / merged / closed）")
	cmd.Flags().String("prompt", "", "分析の切り口を自由記述")
	cmd.Flags().String("repo", "", "分析対象リポジトリ（owner/name）")
	cmd.Flags().StringP("output", "o", "", "レポート出力先ファイルパス（未指定時は標準出力）")
	cmd.Flags().String("model", "", "使用するGeminiモデル（例: gemini-2.5-flash）")
}

func noFlagsSpecified(cmd *cobra.Command) bool {
	return cmd.Flags().NFlag() == 0
}

func validateFlags(cmd *cobra.Command) error {
	today, _ := cmd.Flags().GetBool("today")
	since, _ := cmd.Flags().GetString("since")
	pr, _ := cmd.Flags().GetInt("pr")
	issue, _ := cmd.Flags().GetInt("issue")

	// --today と --since の同時指定はエラー
	if today && since != "" {
		return fmt.Errorf("--today と --since は同時に指定できません")
	}

	// --pr と --issue の同時指定はエラー
	if pr != 0 && issue != 0 {
		return fmt.Errorf("--pr と --issue は同時に指定できません")
	}

	// --since の値が不正な場合エラー
	if since != "" {
		if _, err := entity.ParseDuration(since); err != nil {
			return fmt.Errorf("--since の値が不正です: %w", err)
		}
	}

	return nil
}
