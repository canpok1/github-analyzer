package log

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileWriter_Write(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	w, err := NewFileWriter(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer w.Close()

	err = w.Write("test message")
	if err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unexpected read error: %v", err)
	}
	if !strings.Contains(string(data), "test message") {
		t.Errorf("log file should contain 'test message', got %q", string(data))
	}
}

func TestFileWriter_WriteMultiple(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	w, err := NewFileWriter(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer w.Close()

	_ = w.Write("first")
	_ = w.Write("second")

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unexpected read error: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "first") || !strings.Contains(content, "second") {
		t.Errorf("log file should contain both messages, got %q", content)
	}
}
