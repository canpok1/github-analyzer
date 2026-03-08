package log

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// FileWriter はファイルにログを書き込む実装。
// 複数のgoroutineから安全に使用できる。
type FileWriter struct {
	mu   sync.Mutex
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
	w.mu.Lock()
	defer w.mu.Unlock()
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	_, err := fmt.Fprintf(w.file, "[%s] %s\n", timestamp, message)
	return err
}

func (w *FileWriter) Close() error {
	return w.file.Close()
}
