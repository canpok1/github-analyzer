package cmd

import (
	"github.com/spf13/cobra"
)

var version = "dev"

func makeRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "github-analyzer",
		Short:   "GitHub上の活動プロセスを可視化・分析するCLI診断ツール",
		Version: version,
	}
	cmd.SilenceUsage = true
	return cmd
}

func Execute() error {
	return makeRootCmd().Execute()
}
