package logger

import (
	"os"
	"testing"
)

func TestLogger_Log(t *testing.T) {
	logFile := "test.log"
	l, err := NewLogger(logFile)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}
	defer l.Close()

	testMsg := "Test log message"
	l.Log(testMsg)

	// Проверка записи в файл
	content := readFile(logFile)

	if !contains(content, testMsg) {
		t.Errorf("Log file does not contain message: %s", testMsg)
	}
}
func readFile(filename string) string {
	b, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}
	return string(b)
}

func contains(s, substr string) bool {
	// Вспомогательная функция для теста
	return len(s) > 0 && (func() bool {
		for i := 0; i < len(s); i++ {
			if s[i:i+len(substr)] == substr {
				return true
			}
		}
		return false
	})()
}
