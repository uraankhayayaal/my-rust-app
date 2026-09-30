package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Logger struct {
	mu       sync.Mutex
	file     *os.File
	filename string
	maxSize  int64
}

func NewLogger(filename string, maxSize int64) (*Logger, error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	return &Logger{
		file:     f,
		filename: filename,
		maxSize:  maxSize,
	}, nil
}

func (l *Logger) Log(message string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Проверка размера файла перед записью
	info, err := l.file.Stat()
	if err == nil && info.Size() > l.maxSize {
		l.rotate()
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("[%s] %s\n", timestamp, message)
	l.file.WriteString(line)
}

func (l *Logger) rotate() {
	l.file.Close()

	// Переименовываем текущий файл в .bak
	bakName := l.filename + ".bak"
	os.Rename(l.filename, bakName)

	// Открываем новый файл
	f, err := os.OpenFile(l.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("failed to rotate log file: %v\n", err)
		return
	}
	l.file = f
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.file.Close()
}
