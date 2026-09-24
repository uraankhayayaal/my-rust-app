// domain
package domain

import "time"

// AuditEvent представляет событие аудита.
type AuditEvent struct {
	ID        string
	Timestamp time.Time
	EventType string
	Message   string
	UserID    string
}
