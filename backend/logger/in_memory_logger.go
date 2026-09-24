package logger

import (
	"context"
	"sync"

	"backend/domain"
)

// InMemoryLogger хранит логи в памяти с ограниченным размером буфера.
type InMemoryLogger struct {
	mu   sync.RWMutex
	buf  []domain.AuditEvent
	max  int
}

// NewInMemoryLogger создаёт новый логгер с указанным максимальным числом событий.
func NewInMemoryLogger(maxEvents int) *InMemoryLogger {
	if maxEvents <= 0 {
		maxEvents = 1024
	}
	return &InMemoryLogger{
		buf: make([]domain.AuditEvent, 0, maxEvents),
		max: maxEvents,
	}
}

// Append добавляет событие в буфер без блокировки caller.
func (l *InMemoryLogger) Append(ctx context.Context, event domain.AuditEvent) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.buf = append(l.buf, event)

	if len(l.buf) > l.max {
		overflow := len(l.buf) - l.max
		l.buf = l.buf[overflow:]
	}
}

// GetLogs возвращает последние N событий из буфера.
func (l *InMemoryLogger) GetLogs(limit int) ([]domain.AuditEvent, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if len(l.buf) == 0 {
		return nil, nil
	}

	if limit <= 0 {
		result := make([]domain.AuditEvent, len(l.buf))
		copy(result, l.buf)
		return result, nil
	}

	start := len(l.buf) - limit
	if start < 0 {
		start = 0
	}
	result := make([]domain.AuditEvent, len(l.buf)-start)
	copy(result, l.buf[start:])
	return result, nil
}
