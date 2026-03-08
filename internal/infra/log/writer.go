package log

import (
	"fmt"
	"io"
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
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	w.mu.Lock()
	defer w.mu.Unlock()
	_, err := fmt.Fprintf(w.file, "[%s] %s\n", timestamp, message)
	return err
}

func (w *FileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.file.Close()
}

// NewWarnOnErrorFunc はログ書き込み関数をラップし、
// 初回エラー時にstderrへ警告を出力するLogFuncを返す。
func NewWarnOnErrorFunc(write func(string) error, stderr io.Writer) func(string) {
	var once sync.Once
	return func(msg string) {
		if err := write(msg); err != nil {
			once.Do(func() {
				_, _ = fmt.Fprintf(stderr, "警告: ログ書き込みに失敗しました: %v\n", err)
			})
		}
	}
}
