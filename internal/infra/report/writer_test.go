package report_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/canpok1/github-analyzer/internal/domain"
	"github.com/canpok1/github-analyzer/internal/infra/report"
)

// WriterがReportWriterインターフェースを満たすことを確認
var _ domain.ReportWriter = (*report.Writer)(nil)

func TestWriter_Write_Stdout(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter("", &buf)

	content := "# Report\n\nテスト出力"
	err := w.Write(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if buf.String() != content {
		t.Errorf("expected %q, got %q", content, buf.String())
	}
}

func TestWriter_Write_File(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "report.md")
	w := report.NewWriter(outputPath, nil)

	content := "# Report\n\nファイル出力テスト"
	err := w.Write(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if string(got) != content {
		t.Errorf("expected %q, got %q", content, string(got))
	}
}

func TestWriter_Write_InvalidPath(t *testing.T) {
	w := report.NewWriter("/nonexistent/dir/report.md", nil)

	err := w.Write("test")
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}
