package logger

import (
	"fmt"
	"os"
	"time"
)

type Logger struct {
	file *os.File
}

func NewLogger(filename string) (*Logger, error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	return &Logger{file: f}, nil
}

func (l *Logger) Log(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("[%s] %s\n", timestamp, message)
	l.file.WriteString(line)
}

func (l *Logger) Close() error {
	return l.file.Close()
}
