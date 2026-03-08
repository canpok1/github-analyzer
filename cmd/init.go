package cmd

import (
	"fmt"
	"os"

	"github.com/canpok1/github-analyzer/templates"
	"github.com/spf13/cobra"
)

func makeInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "設定ファイルのテンプレートを生成する",
		RunE: func(cmd *cobra.Command, args []string) error {
			const fileName = ".github-analyzer.yaml"
			f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
			if err != nil {
				if os.IsExist(err) {
					return fmt.Errorf("%s は既に存在します", fileName)
				}
				return err
			}
			defer func() { _ = f.Close() }()

			if _, err := f.Write(templates.ConfigTemplate); err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s を作成しました\n", fileName)
			return nil
		},
	}
}
