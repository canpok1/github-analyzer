//go:build e2e

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "github-analyzer-e2e-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	binaryPath = filepath.Join(dir, "github-analyzer")
	cmd := exec.Command("go", "build", "-o", binaryPath, "github.com/canpok1/github-analyzer")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic("failed to build binary: " + err.Error())
	}

	os.Exit(m.Run())
}

// filterEnv は環境変数リストから指定キーを除外する。
func filterEnv(env []string, keys ...string) []string {
	filtered := make([]string, 0, len(env))
	for _, e := range env {
		exclude := false
		for _, key := range keys {
			if strings.HasPrefix(e, key+"=") {
				exclude = true
				break
			}
		}
		if !exclude {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

// envWithoutTokens はトークン系環境変数を全て除外した環境変数リストを返す。
func envWithoutTokens() []string {
	return filterEnv(os.Environ(), "GH_TOKEN", "GITHUB_TOKEN", "GEMINI_API_KEY")
}

// runCLI はCLIバイナリを指定引数・環境変数で実行し、出力とエラーを返す。
// env が nil の場合は現在の環境変数をそのまま使用する。
func runCLI(args []string, env []string) (string, error) {
	cmd := exec.Command(binaryPath, args...)
	if env != nil {
		cmd.Env = env
	}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// assertCLIError はCLIがエラー終了し、出力に期待文字列が含まれることを確認するヘルパー。
func assertCLIError(t *testing.T, args []string, env []string, expectedSubstrings ...string) {
	t.Helper()
	output, err := runCLI(args, env)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	for _, s := range expectedSubstrings {
		if !strings.Contains(output, s) {
			t.Errorf("expected output to contain %q, got: %s", s, output)
		}
	}
}

// assertTokenError はバリデーション通過後にトークンエラーで失敗することを確認するヘルパー。
func assertTokenError(t *testing.T, args []string) {
	t.Helper()
	output, err := runCLI(args, envWithoutTokens())
	if err == nil {
		t.Fatal("expected error (missing token), got nil")
	}
	if strings.Contains(output, "いずれかを指定してください") {
		t.Errorf("should not be a validation error, expected token error: %s", output)
	}
}

func TestCLIHelp(t *testing.T) {
	output, err := runCLI([]string{"--help"}, nil)
	if err != nil {
		t.Fatalf("failed to run CLI with --help: %v\n%s", err, output)
	}

	if !strings.Contains(output, "github-analyzer") {
		t.Error("help output should contain 'github-analyzer'")
	}

	expectedFlags := []string{"--today", "--since", "--pr", "--issue", "--output", "--repo", "--prompt", "--status"}
	for _, flag := range expectedFlags {
		if !strings.Contains(output, flag) {
			t.Errorf("help output should contain %q", flag)
		}
	}
}

func TestCLIVersion(t *testing.T) {
	output, err := runCLI([]string{"--version"}, nil)
	if err != nil {
		t.Fatalf("failed to run CLI: %v\n%s", err, output)
	}
	if len(output) == 0 {
		t.Fatal("expected version output, got empty")
	}
}

func TestCLINoFlags(t *testing.T) {
	assertCLIError(t, nil, nil, "--today")
}

func TestCLIUnknownFlag(t *testing.T) {
	assertCLIError(t, []string{"--unknown-flag"}, nil, "unknown flag")
}

func TestCLITodayAndSinceConflict(t *testing.T) {
	assertCLIError(t, []string{"--today", "--since", "7d"}, nil, "--today", "--since")
}

func TestCLIPRAndIssueConflict(t *testing.T) {
	assertCLIError(t, []string{"--pr", "123", "--issue", "456"}, nil, "--pr", "--issue")
}

func TestCLITodayAndSinceAndPRConflict(t *testing.T) {
	assertCLIError(t, []string{"--today", "--since", "7d", "--pr", "1"}, nil, "--today", "--since")
}

func TestCLISinceInvalidValue(t *testing.T) {
	assertCLIError(t, []string{"--since", "invalid"}, nil, "--since")
}

func TestCLISinceValidValues(t *testing.T) {
	for _, value := range []string{"7d", "2w", "1m"} {
		t.Run(value, func(t *testing.T) {
			output, err := runCLI([]string{"--since", value}, envWithoutTokens())
			if err == nil {
				t.Fatal("expected error (missing token), got nil")
			}
			if strings.Contains(output, "--since の値が不正") {
				t.Errorf("should not be a since validation error for %s: %s", value, output)
			}
		})
	}
}

func TestCLIZeroValueTargets(t *testing.T) {
	for _, flag := range []string{"--pr", "--issue"} {
		t.Run(flag, func(t *testing.T) {
			output, err := runCLI([]string{flag, "0"}, nil)
			if err == nil {
				t.Fatal("expected error when " + flag + " 0 is specified")
			}
			if !strings.Contains(output, "いずれかを指定してください") && !strings.Contains(output, "--today") {
				t.Errorf("error output should indicate no target specified: %s", output)
			}
		})
	}
}

func TestCLIMissingGHToken(t *testing.T) {
	output, err := runCLI(
		[]string{"--pr", "1", "--repo", "owner/repo"},
		filterEnv(os.Environ(), "GH_TOKEN", "GITHUB_TOKEN"),
	)
	if err == nil {
		t.Fatal("expected error when GH_TOKEN/GITHUB_TOKEN are not set")
	}
	if !strings.Contains(output, "GH_TOKEN") && !strings.Contains(output, "GITHUB_TOKEN") {
		t.Errorf("error output should mention GH_TOKEN or GITHUB_TOKEN: %s", output)
	}
}

func TestCLIMissingGeminiAPIKey(t *testing.T) {
	env := filterEnv(os.Environ(), "GEMINI_API_KEY")
	env = append(env, "GH_TOKEN=dummy-token")
	output, err := runCLI([]string{"--pr", "1", "--repo", "owner/repo"}, env)
	if err == nil {
		t.Fatal("expected error when GEMINI_API_KEY is not set")
	}
	if !strings.Contains(output, "GEMINI_API_KEY") {
		t.Errorf("error output should mention GEMINI_API_KEY: %s", output)
	}
}

func TestCLIStatusFlag(t *testing.T) {
	for _, status := range []string{"open", "merged", "closed"} {
		t.Run(status, func(t *testing.T) {
			assertTokenError(t, []string{"--today", "--status", status})
		})
	}
}

func TestCLIPromptFlag(t *testing.T) {
	assertTokenError(t, []string{"--pr", "1", "--repo", "owner/repo", "--prompt", "コードの品質を分析してください"})
}

func TestCLIIssueFlag(t *testing.T) {
	assertTokenError(t, []string{"--issue", "42", "--repo", "owner/repo"})
}

func TestCLIValidFlagCombinations(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"today+pr", []string{"--today", "--pr", "10", "--repo", "owner/repo"}},
		{"today+issue", []string{"--today", "--issue", "5", "--repo", "owner/repo"}},
		{"since+pr", []string{"--since", "7d", "--pr", "10", "--repo", "owner/repo"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertTokenError(t, tt.args)
		})
	}
}

func TestCLIOutputFlag(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{"simple", "report.md"},
		{"subdir", filepath.Join("subdir", "report.md")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			outputPath := filepath.Join(tmpDir, tt.path)
			assertTokenError(t, []string{"--pr", "1", "--repo", "owner/repo", "--output", outputPath})
		})
	}
}
