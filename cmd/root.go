package cmd

import (
	"fmt"
	"os"

	"github.com/canpok1/github-analyzer/internal/app"
	"github.com/canpok1/github-analyzer/internal/domain"
	"github.com/canpok1/github-analyzer/internal/infra/config"
	"github.com/canpok1/github-analyzer/internal/infra/gemini"
	ghclient "github.com/canpok1/github-analyzer/internal/infra/github"
	applog "github.com/canpok1/github-analyzer/internal/infra/log"
	"github.com/canpok1/github-analyzer/internal/infra/mock"
	"github.com/canpok1/github-analyzer/internal/infra/report"
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
	cmd.AddCommand(makeInitCmd())
	defineFlags(cmd)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if noFlagsSpecified(cmd) {
			return cmd.Help()
		}
		if err := validateFlags(cmd); err != nil {
			return err
		}
		return runAnalyze(cmd)
	}
	return cmd
}

func Execute() error {
	return makeRootCmd().Execute()
}

// runAnalyze はDIを行い、app.Analyzeを実行する。
func runAnalyze(cmd *cobra.Command) error {
	query := buildQuery(cmd)

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
	}

	query = applyConfig(query, cfg)

	// GitHubRepository の解決
	var ghRepo domain.GitHubRepository
	if cfg.Mock.Repository {
		ghRepo = &mock.GitHubRepository{}
	} else {
		token, err := resolveToken()
		if err != nil {
			return err
		}
		ghRepo = ghclient.NewClient(token)
	}

	// Analyzer の解決
	var analyzer domain.Analyzer
	if cfg.Mock.AI {
		analyzer = &mock.Analyzer{}
	} else {
		geminiAPIKey := os.Getenv("GEMINI_API_KEY")
		if geminiAPIKey == "" {
			return fmt.Errorf("GEMINI_API_KEY 環境変数を設定してください")
		}
		geminiClient, err := gemini.NewClient(geminiAPIKey)
		if err != nil {
			return err
		}
		geminiClient.SetModel(query.Model)
		analyzer = geminiClient
	}

	// Logger の解決
	var logger app.LogFunc
	if cfg.LogFile != "" {
		logWriter, err := applog.NewFileWriter(cfg.LogFile)
		if err != nil {
			return fmt.Errorf("ログファイルの作成に失敗しました: %w", err)
		}
		defer func() { _ = logWriter.Close() }()
		logger = applog.NewWarnOnErrorFunc(logWriter.Write, os.Stderr)
	}

	outputPath, _ := cmd.Flags().GetString("output")
	renderer := report.NewMarkdownRenderer()
	writer := report.NewWriter(outputPath, cmd.OutOrStdout())

	deps := app.AnalyzeDeps{
		GitHubRepo:    ghRepo,
		Analyzer:      analyzer,
		PromptBuilder: gemini.BuildPrompt,
		ReportParser:  gemini.ParseReport,
		Renderer:      renderer,
		Writer:        writer,
		Stderr:        os.Stderr,
		Logger:        logger,
	}

	return app.Analyze(cmd.Context(), deps, query)
}
