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

func TestCLIHelp(t *testing.T) {
	cmd := exec.Command(binaryPath, "--help")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run CLI with --help: %v\n%s", err, out)
	}

	output := string(out)

	// コマンド名が含まれる
	if !strings.Contains(output, "github-analyzer") {
		t.Error("help output should contain 'github-analyzer'")
	}

	// 主要フラグが含まれる
	expectedFlags := []string{"--today", "--since", "--pr", "--issue", "--output", "--repo", "--prompt", "--status"}
	for _, flag := range expectedFlags {
		if !strings.Contains(output, flag) {
			t.Errorf("help output should contain %q", flag)
		}
	}
}

func TestCLINoFlags(t *testing.T) {
	cmd := exec.Command(binaryPath)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error when no flags specified, got nil")
	}

	output := string(out)
	if !strings.Contains(output, "--today") {
		t.Error("error output should mention available flags like --today")
	}
}

func TestCLITodayAndSinceConflict(t *testing.T) {
	cmd := exec.Command(binaryPath, "--today", "--since", "7d")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error when --today and --since are both specified")
	}

	output := string(out)
	if !strings.Contains(output, "--today") || !strings.Contains(output, "--since") {
		t.Errorf("error output should mention --today and --since conflict: %s", output)
	}
}

func TestCLIPRAndIssueConflict(t *testing.T) {
	cmd := exec.Command(binaryPath, "--pr", "123", "--issue", "456")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error when --pr and --issue are both specified")
	}

	output := string(out)
	if !strings.Contains(output, "--pr") || !strings.Contains(output, "--issue") {
		t.Errorf("error output should mention --pr and --issue conflict: %s", output)
	}
}

func TestCLISinceInvalidValue(t *testing.T) {
	cmd := exec.Command(binaryPath, "--since", "invalid")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error for invalid --since value")
	}

	output := string(out)
	if !strings.Contains(output, "--since") {
		t.Errorf("error output should mention --since: %s", output)
	}
}

func TestCLIMissingGHToken(t *testing.T) {
	cmd := exec.Command(binaryPath, "--pr", "1", "--repo", "owner/repo")
	// 環境変数をクリアした状態で実行
	cmd.Env = filterEnv(os.Environ(), "GH_TOKEN", "GITHUB_TOKEN")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error when GH_TOKEN/GITHUB_TOKEN are not set")
	}

	output := string(out)
	if !strings.Contains(output, "GH_TOKEN") && !strings.Contains(output, "GITHUB_TOKEN") {
		t.Errorf("error output should mention GH_TOKEN or GITHUB_TOKEN: %s", output)
	}
}

func TestCLIMissingGeminiAPIKey(t *testing.T) {
	cmd := exec.Command(binaryPath, "--pr", "1", "--repo", "owner/repo")
	// GH_TOKENは設定し、GEMINI_API_KEYをクリア
	env := filterEnv(os.Environ(), "GEMINI_API_KEY")
	env = append(env, "GH_TOKEN=dummy-token")
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error when GEMINI_API_KEY is not set")
	}

	output := string(out)
	if !strings.Contains(output, "GEMINI_API_KEY") {
		t.Errorf("error output should mention GEMINI_API_KEY: %s", output)
	}
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
func runCLI(args []string, env []string) (string, error) {
	cmd := exec.Command(binaryPath, args...)
	if env != nil {
		cmd.Env = env
	}
	out, err := cmd.CombinedOutput()
	return string(out), err
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

func TestCLIOutputFlag(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "report.md")
	assertTokenError(t, []string{"--pr", "1", "--repo", "owner/repo", "--output", outputPath})
}

func TestCLIUnknownFlag(t *testing.T) {
	cmd := exec.Command(binaryPath, "--unknown-flag")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}

	output := string(out)
	if !strings.Contains(output, "unknown flag") {
		t.Errorf("error output should mention 'unknown flag': %s", output)
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

func TestCLIPRZeroValue(t *testing.T) {
	// --pr 0 は対象未指定と同等
	cmd := exec.Command(binaryPath, "--pr", "0")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error when --pr 0 is specified")
	}

	output := string(out)
	if !strings.Contains(output, "いずれかを指定してください") && !strings.Contains(output, "--today") {
		t.Errorf("error output should indicate no target specified: %s", output)
	}
}

func TestCLIIssueZeroValue(t *testing.T) {
	// --issue 0 は対象未指定と同等
	cmd := exec.Command(binaryPath, "--issue", "0")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error when --issue 0 is specified")
	}

	output := string(out)
	if !strings.Contains(output, "いずれかを指定してください") && !strings.Contains(output, "--today") {
		t.Errorf("error output should indicate no target specified: %s", output)
	}
}

func TestCLITodayAndSinceAndPRConflict(t *testing.T) {
	cmd := exec.Command(binaryPath, "--today", "--since", "7d", "--pr", "1")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error when --today, --since, and --pr are all specified")
	}

	output := string(out)
	if !strings.Contains(output, "--today") || !strings.Contains(output, "--since") {
		t.Errorf("error output should mention --today and --since conflict: %s", output)
	}
}

func TestCLIOutputWithSubdirPath(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "subdir", "report.md")
	assertTokenError(t, []string{"--pr", "1", "--repo", "owner/repo", "--output", outputPath})
}

func TestCLIVersion(t *testing.T) {
	cmd := exec.Command(binaryPath, "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run CLI: %v\n%s", err, out)
	}
	if len(out) == 0 {
		t.Fatal("expected version output, got empty")
	}
}
