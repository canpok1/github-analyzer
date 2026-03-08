package log

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
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

func TestFileWriter_WriteConcurrent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	w, err := NewFileWriter(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer w.Close()

	const goroutines = 10
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := range goroutines {
		go func(n int) {
			defer wg.Done()
			if err := w.Write(strings.Repeat("x", 100)); err != nil {
				t.Errorf("goroutine %d: unexpected write error: %v", n, err)
			}
		}(i)
	}
	wg.Wait()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unexpected read error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != goroutines {
		t.Errorf("expected %d lines, got %d", goroutines, len(lines))
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

	if err := w.Write("first"); err != nil {
		t.Fatalf("unexpected write error for 'first': %v", err)
	}
	if err := w.Write("second"); err != nil {
		t.Fatalf("unexpected write error for 'second': %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unexpected read error: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "first") || !strings.Contains(content, "second") {
		t.Errorf("log file should contain both messages, got %q", content)
	}
}
