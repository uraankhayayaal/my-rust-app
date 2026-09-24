package domain

import "time"

// LogLevel represents the severity level of a log entry.
type LogLevel string

const (
	LogInfo  LogLevel = "info"
	LogError LogLevel = "error"
)

// AuditEvent is an immutable value object that records the fact of sending a greeting.
type AuditEvent struct {
	Timestamp   time.Time `json:"timestamp"`
	Method      string    `json:"method"`       // HTTP method: "GET", "POST", etc.
	Endpoint    string    `json:"endpoint"`     // URL path: "/hello"
	SessionID   string    `json:"session_id,omitempty"` // client session identifier, if available
	Result      string    `json:"result"`       // outcome: "success", "error", "skipped"
	Level       LogLevel  `json:"level"`        // "info" or "error"
	Message     string    `json:"message"`      // brief description of the event
	HelloText   string    `json:"hello_text,omitempty"` // text that was sent in the greeting, if applicable
}

// NewSuccessEvent creates a new AuditEvent with level info and result success.
func NewSuccessEvent(method, endpoint, message string) AuditEvent {
	return AuditEvent{
		Timestamp: time.Now(),
		Method:    method,
		Endpoint:  endpoint,
		Result:    "success",
		Level:     LogInfo,
		Message:   message,
	}
}

// NewErrorEvent creates a new AuditEvent with level error and result error.
func NewErrorEvent(method, endpoint, errorMessage string) AuditEvent {
	return AuditEvent{
		Timestamp: time.Now(),
		Method:    method,
		Endpoint:  endpoint,
		Result:    "error",
		Level:     LogError,
		Message:   errorMessage,
	}
}
