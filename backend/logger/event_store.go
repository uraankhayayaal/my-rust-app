package logger

import (
	"context"

	"backend/domain"
)

// EventStore интерфейс для всех имплементаций хранения логов.
type EventStore interface {
	// Append добавляет событие в буфер. Не блокирует caller (async).
	Append(ctx context.Context, event domain.AuditEvent)

	// GetLogs возвращает последние N событий из буфера.
	// limit=0 или отрицательный = вернуть все.
	// Если events=nil и err=nil — пустой список.
	GetLogs(limit int) ([]domain.AuditEvent, error)
}
