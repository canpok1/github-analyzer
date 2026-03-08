package config

import (
	"os"
	"path/filepath"
	"testing"
)

// DONE: 正常系: 全フィールドが設定されたYAMLを読み込める
// DONE: 正常系: 一部フィールドのみ設定されたYAMLを読み込める
// DONE: 正常系: 空のYAMLファイルはゼロ値のConfigを返す
// DONE: 異常系: ファイルが存在しない場合はゼロ値のConfigを返す（エラーにしない）
// DONE: 異常系: 不正なYAMLの場合エラーを返す

func TestLoadFromPath_AllFields(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".github-analyzer.yaml")

	content := `repo: owner/repo
tone: friendly
default_prompt: チームの活動を分析してください
model: gemini-2.0-flash
`
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadFromPath(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Repo != "owner/repo" {
		t.Errorf("Repo = %q, want %q", cfg.Repo, "owner/repo")
	}
	if cfg.Tone != "friendly" {
		t.Errorf("Tone = %q, want %q", cfg.Tone, "friendly")
	}
	if cfg.DefaultPrompt != "チームの活動を分析してください" {
		t.Errorf("DefaultPrompt = %q, want %q", cfg.DefaultPrompt, "チームの活動を分析してください")
	}
	if cfg.Model != "gemini-2.0-flash" {
		t.Errorf("Model = %q, want %q", cfg.Model, "gemini-2.0-flash")
	}
}

func TestLoadFromPath_PartialFields(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".github-analyzer.yaml")

	content := `repo: owner/repo
model: gemini-2.0-flash
`
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadFromPath(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Repo != "owner/repo" {
		t.Errorf("Repo = %q, want %q", cfg.Repo, "owner/repo")
	}
	if cfg.Tone != "" {
		t.Errorf("Tone = %q, want empty", cfg.Tone)
	}
	if cfg.DefaultPrompt != "" {
		t.Errorf("DefaultPrompt = %q, want empty", cfg.DefaultPrompt)
	}
	if cfg.Model != "gemini-2.0-flash" {
		t.Errorf("Model = %q, want %q", cfg.Model, "gemini-2.0-flash")
	}
}

func TestLoadFromPath_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".github-analyzer.yaml")

	if err := os.WriteFile(configPath, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadFromPath(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Repo != "" {
		t.Errorf("Repo = %q, want empty", cfg.Repo)
	}
	if cfg.Tone != "" {
		t.Errorf("Tone = %q, want empty", cfg.Tone)
	}
	if cfg.DefaultPrompt != "" {
		t.Errorf("DefaultPrompt = %q, want empty", cfg.DefaultPrompt)
	}
	if cfg.Model != "" {
		t.Errorf("Model = %q, want empty", cfg.Model)
	}
}

func TestLoadFromPath_FileNotExist(t *testing.T) {
	cfg, err := LoadFromPath("/nonexistent/path/.github-analyzer.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Repo != "" {
		t.Errorf("Repo = %q, want empty", cfg.Repo)
	}
}

func TestLoadFromPath_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".github-analyzer.yaml")

	content := `{invalid: yaml: [broken`
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadFromPath(configPath)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestLoad_UsesHomeDir(t *testing.T) {
	// HOMEを一時ディレクトリに差し替えてテスト
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Chdir(t.TempDir()) // カレントには設定ファイルなし

	configPath := filepath.Join(dir, ".github-analyzer.yaml")
	content := `repo: test/repo
`
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Repo != "test/repo" {
		t.Errorf("Repo = %q, want %q", cfg.Repo, "test/repo")
	}
}

func TestLoadFromPath_MockConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".github-analyzer.yaml")

	content := `mock:
  ai: true
  repository: true
log_file: /tmp/test.log
`
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadFromPath(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !cfg.Mock.AI {
		t.Error("Mock.AI should be true")
	}
	if !cfg.Mock.Repository {
		t.Error("Mock.Repository should be true")
	}
	if cfg.LogFile != "/tmp/test.log" {
		t.Errorf("LogFile = %q, want %q", cfg.LogFile, "/tmp/test.log")
	}
}

func TestLoadFromPath_MockDefaults(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".github-analyzer.yaml")

	content := `repo: owner/repo`
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadFromPath(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Mock.AI {
		t.Error("Mock.AI should default to false")
	}
	if cfg.Mock.Repository {
		t.Error("Mock.Repository should default to false")
	}
	if cfg.LogFile != "" {
		t.Errorf("LogFile should default to empty, got %q", cfg.LogFile)
	}
}

func TestLoadFromPath_PartialMock(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".github-analyzer.yaml")

	content := `mock:
  ai: true
`
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadFromPath(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !cfg.Mock.AI {
		t.Error("Mock.AI should be true")
	}
	if cfg.Mock.Repository {
		t.Error("Mock.Repository should be false when not specified")
	}
}

func TestLoad_NoConfigFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Chdir(t.TempDir())

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Repo != "" {
		t.Errorf("Repo = %q, want empty", cfg.Repo)
	}
}

func TestLoad_CurrentDirOnly(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	cwdDir := t.TempDir()
	t.Chdir(cwdDir)

	content := `repo: current/repo
tone: casual
`
	if err := os.WriteFile(filepath.Join(cwdDir, ".github-analyzer.yaml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Repo != "current/repo" {
		t.Errorf("Repo = %q, want %q", cfg.Repo, "current/repo")
	}
	if cfg.Tone != "casual" {
		t.Errorf("Tone = %q, want %q", cfg.Tone, "casual")
	}
}

func TestLoad_CurrentDirOverridesHome(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	homeContent := `repo: home/repo
tone: friendly
model: gemini-2.0-flash
`
	if err := os.WriteFile(filepath.Join(homeDir, ".github-analyzer.yaml"), []byte(homeContent), 0o644); err != nil {
		t.Fatal(err)
	}

	cwdDir := t.TempDir()
	t.Chdir(cwdDir)

	cwdContent := `repo: current/repo
`
	if err := os.WriteFile(filepath.Join(cwdDir, ".github-analyzer.yaml"), []byte(cwdContent), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// カレントで設定されたrepoはカレントの値
	if cfg.Repo != "current/repo" {
		t.Errorf("Repo = %q, want %q", cfg.Repo, "current/repo")
	}
	// カレントで未設定のtoneはホームの値を維持
	if cfg.Tone != "friendly" {
		t.Errorf("Tone = %q, want %q", cfg.Tone, "friendly")
	}
	// カレントで未設定のmodelはホームの値を維持
	if cfg.Model != "gemini-2.0-flash" {
		t.Errorf("Model = %q, want %q", cfg.Model, "gemini-2.0-flash")
	}
}

func TestLoad_NestedMockFieldMerge(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	homeContent := `mock:
  ai: true
  repository: true
`
	if err := os.WriteFile(filepath.Join(homeDir, ".github-analyzer.yaml"), []byte(homeContent), 0o644); err != nil {
		t.Fatal(err)
	}

	cwdDir := t.TempDir()
	t.Chdir(cwdDir)

	cwdContent := `mock:
  ai: false
`
	if err := os.WriteFile(filepath.Join(cwdDir, ".github-analyzer.yaml"), []byte(cwdContent), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// カレントで設定されたmock.aiはカレントの値(false)
	if cfg.Mock.AI {
		t.Error("Mock.AI should be false (overridden by current dir)")
	}
	// カレントで未設定のmock.repositoryはホームの値(true)を維持
	if !cfg.Mock.Repository {
		t.Error("Mock.Repository should be true (kept from home dir)")
	}
}

func TestLoad_HomeDirOnly(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	homeContent := `repo: home/repo
tone: friendly
`
	if err := os.WriteFile(filepath.Join(homeDir, ".github-analyzer.yaml"), []byte(homeContent), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Chdir(t.TempDir()) // カレントには設定ファイルなし

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Repo != "home/repo" {
		t.Errorf("Repo = %q, want %q", cfg.Repo, "home/repo")
	}
	if cfg.Tone != "friendly" {
		t.Errorf("Tone = %q, want %q", cfg.Tone, "friendly")
	}
}
