package report

import (
	"fmt"
	"io"
	"os"
)

// Writer はレンダリング済みレポートを出力先に書き込む。
type Writer struct {
	outputPath string
	writer     io.Writer
}

// NewWriter はWriterを生成する。
// outputPathが空文字列の場合はwriterに書き込む。
// outputPathが指定された場合はファイルに書き込む。
func NewWriter(outputPath string, writer io.Writer) *Writer {
	return &Writer{
		outputPath: outputPath,
		writer:     writer,
	}
}

// Write はレンダリング済みレポートを出力先に書き込む。
func (w *Writer) Write(content string) error {
	if w.outputPath != "" {
		return w.writeToFile(content)
	}
	_, err := fmt.Fprint(w.writer, content)
	return err
}

// writeToFile はファイルにレポートを書き込む。
func (w *Writer) writeToFile(content string) error {
	if err := os.WriteFile(w.outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}
	return nil
}
