package logger

import (
	"os"
	"testing"
	"strings"
)

func TestLogger_Log(t *testing.T) {
	logFile := "test.log"
	// Создаем логгер с маленьким лимитом для проверки ротации
	l, err := NewLogger(logFile, 10) 
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}
	defer l.Close()

	testMsg := "Test log message"
	l.Log(testMsg)

	content := readFile(logFile)
	if !strings.Contains(content, testMsg) {
		t.Errorf("Log file does not contain message: %s", testMsg)
	}
}

func TestLogger_Rotation(t *testing.T) {
	logFile := "rotate.log"
	// Очень маленький размер, чтобы вызвать ротацию почти сразу
	l, err := NewLogger(logFile, 5)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}
	defer l.Close()

	// Записываем достаточно данных, чтобы размер превысил 5 байт
	l.Log("First log line") 
	l.Log("Second log line triggers rotation")

	// После второй записи должен быть создан rotate.log.bak
	if _, err := os.Stat(logFile + ".bak"); os.IsNotExist(err) {
		t.Errorf("Rotation failed: .bak file not found")
	}

	content := readFile(logFile)
	if !strings.Contains(content, "Second log line") {
		t.Errorf("New log file does not contain the latest message")
	}
}

func readFile(filename string) string {
	b, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}
	return string(b)
}
