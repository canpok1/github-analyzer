package cmd

import (
	"strings"
	"testing"
	"time"
)

// DONE: 正常系: GH_TOKENが優先される
// DONE: 正常系: GH_TOKEN未設定時にGITHUB_TOKENが使われる
// DONE: 異常系: GH_TOKEN/GITHUB_TOKEN両方未設定の場合エラーを返す
// DONE: 正常系: buildQueryがフラグからQueryを正しく構築する
// DONE: 正常系: --modelフラグでQuery.Modelが設定される
// DONE: 正常系: --todayフラグでSinceが当日0時に設定される
// DONE: 正常系: --sinceフラグでSinceが正しく計算される
// DONE: 異常系: GEMINI_API_KEY未設定の場合エラーを返す
// DONE: 異常系: GH_TOKEN/GITHUB_TOKEN未設定の場合エラーを返す

func TestResolveToken_GHTokenPriority(t *testing.T) {
	t.Setenv("GH_TOKEN", "gh-token-value")
	t.Setenv("GITHUB_TOKEN", "github-token-value")

	token, err := resolveToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "gh-token-value" {
		t.Errorf("expected 'gh-token-value', got %q", token)
	}
}

func TestResolveToken_FallbackToGitHubToken(t *testing.T) {
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "github-token-value")

	token, err := resolveToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "github-token-value" {
		t.Errorf("expected 'github-token-value', got %q", token)
	}
}

func TestResolveToken_BothUnset(t *testing.T) {
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")

	_, err := resolveToken()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "GH_TOKEN") {
		t.Errorf("error should mention GH_TOKEN: %v", err)
	}
}

func TestBuildQuery_WithPR(t *testing.T) {
	cmd := makeRootCmd()
	if err := cmd.ParseFlags([]string{"--pr", "42", "--repo", "owner/repo"}); err != nil {
		t.Fatalf("failed to parse flags: %v", err)
	}

	query := buildQuery(cmd)
	if query.PR == nil || *query.PR != 42 {
		t.Errorf("expected PR=42, got %v", query.PR)
	}
	if query.Repo != "owner/repo" {
		t.Errorf("expected Repo='owner/repo', got %q", query.Repo)
	}
}

func TestBuildQuery_WithIssue(t *testing.T) {
	cmd := makeRootCmd()
	if err := cmd.ParseFlags([]string{"--issue", "10", "--repo", "o/r"}); err != nil {
		t.Fatalf("failed to parse flags: %v", err)
	}

	query := buildQuery(cmd)
	if query.Issue == nil || *query.Issue != 10 {
		t.Errorf("expected Issue=10, got %v", query.Issue)
	}
}

func TestBuildQuery_WithToday(t *testing.T) {
	cmd := makeRootCmd()
	if err := cmd.ParseFlags([]string{"--today", "--repo", "o/r"}); err != nil {
		t.Fatalf("failed to parse flags: %v", err)
	}

	query := buildQuery(cmd)
	if query.Since == nil {
		t.Fatal("expected Since to be set")
	}

	now := time.Now()
	expectedDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if !query.Since.Equal(expectedDate) {
		t.Errorf("expected Since=%v, got %v", expectedDate, *query.Since)
	}
}

func TestBuildQuery_WithSince(t *testing.T) {
	cmd := makeRootCmd()
	if err := cmd.ParseFlags([]string{"--since", "7d", "--repo", "o/r"}); err != nil {
		t.Fatalf("failed to parse flags: %v", err)
	}

	query := buildQuery(cmd)
	if query.Since == nil {
		t.Fatal("expected Since to be set")
	}

	// 7日前のおおよその時刻であること
	expected := time.Now().Add(-7 * 24 * time.Hour)
	diff := query.Since.Sub(expected)
	if diff < -time.Second || diff > time.Second {
		t.Errorf("expected Since ~= %v, got %v", expected, *query.Since)
	}
}

func TestRunAnalyze_MissingGeminiAPIKey(t *testing.T) {
	t.Setenv("GH_TOKEN", "test-token")
	t.Setenv("GEMINI_API_KEY", "")

	cmd := makeRootCmd()
	cmd.SetArgs([]string{"--pr", "1", "--repo", "owner/repo"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "GEMINI_API_KEY") {
		t.Errorf("error should mention GEMINI_API_KEY: %v", err)
	}
}

func TestRunAnalyze_MissingGHToken(t *testing.T) {
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("GEMINI_API_KEY", "test-key")

	cmd := makeRootCmd()
	cmd.SetArgs([]string{"--pr", "1", "--repo", "owner/repo"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "GH_TOKEN") {
		t.Errorf("error should mention GH_TOKEN: %v", err)
	}
}

func TestBuildQuery_WithModel(t *testing.T) {
	cmd := makeRootCmd()
	if err := cmd.ParseFlags([]string{"--pr", "1", "--repo", "o/r", "--model", "gemini-2.5-flash"}); err != nil {
		t.Fatalf("failed to parse flags: %v", err)
	}

	query := buildQuery(cmd)
	if query.Model != "gemini-2.5-flash" {
		t.Errorf("expected Model='gemini-2.5-flash', got %q", query.Model)
	}
}

func TestBuildQuery_WithPrompt(t *testing.T) {
	cmd := makeRootCmd()
	if err := cmd.ParseFlags([]string{"--pr", "1", "--repo", "o/r", "--prompt", "check quality"}); err != nil {
		t.Fatalf("failed to parse flags: %v", err)
	}

	query := buildQuery(cmd)
	if query.Prompt != "check quality" {
		t.Errorf("expected Prompt='check quality', got %q", query.Prompt)
	}
}
