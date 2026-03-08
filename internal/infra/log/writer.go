package log

import (
	"fmt"
	"os"
	"time"
)

// FileWriter はファイルにログを書き込む実装。
type FileWriter struct {
	file *os.File
}

// NewFileWriter は指定パスにログファイルを作成する。
func NewFileWriter(path string) (*FileWriter, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, fmt.Errorf("ログファイルを開けません: %w", err)
	}
	return &FileWriter{file: f}, nil
}

func (w *FileWriter) Write(message string) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	_, err := fmt.Fprintf(w.file, "[%s] %s\n", timestamp, message)
	return err
}

func (w *FileWriter) Close() error {
	return w.file.Close()
}
