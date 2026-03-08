package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// DONE: 正常系: テンプレートファイルがカレントディレクトリに生成される
// DONE: 異常系: 既にファイルが存在する場合はエラーになり上書きされない
// DONE: 正常系: テンプレートに各フィールドのコメント付きサンプルが含まれている

func TestInitCmd_CreatesTemplateFile(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(dir)

	cmd := makeRootCmd()
	cmd.SetArgs([]string{"init"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	configPath := filepath.Join(dir, ".github-analyzer.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("expected .github-analyzer.yaml to be created, but it does not exist")
	}
}

func TestInitCmd_FileAlreadyExists(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(dir)

	configPath := filepath.Join(dir, ".github-analyzer.yaml")
	originalContent := "original content"
	if err := os.WriteFile(configPath, []byte(originalContent), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := makeRootCmd()
	cmd.SetArgs([]string{"init"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when file already exists, got nil")
	}

	data, _ := os.ReadFile(configPath)
	if string(data) != originalContent {
		t.Errorf("file was overwritten: got %q, want %q", string(data), originalContent)
	}
}

func TestInitCmd_TemplateContainsFieldComments(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(dir)

	cmd := makeRootCmd()
	cmd.SetArgs([]string{"init"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	configPath := filepath.Join(dir, ".github-analyzer.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}

	content := string(data)
	expectedFields := []string{"repo:", "tone:", "default_prompt:", "model:"}
	for _, field := range expectedFields {
		if !strings.Contains(content, field) {
			t.Errorf("template does not contain field %q", field)
		}
	}
}
