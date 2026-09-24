package domain

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewSuccessEvent(t *testing.T) {
	method := "POST"
	endpoint := "/hello"
	message := "Greeting sent successfully"
	helloText := "Hello, World!"

	evt := NewSuccessEvent(method, endpoint, message)
	_ = helloText

	if evt.Method != method {
		t.Errorf("Method = %q, want %q", evt.Method, method)
	}
	if evt.Endpoint != endpoint {
		t.Errorf("Endpoint = %q, want %q", evt.Endpoint, endpoint)
	}
	if evt.Message != message {
		t.Errorf("Message = %q, want %q", evt.Message, message)
	}
	if evt.Result != "success" {
		t.Errorf("Result = %q, want %q", evt.Result, "success")
	}
	if evt.Level != LogInfo {
		t.Errorf("Level = %q, want %q", evt.Level, LogInfo)
	}
	if evt.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
	// HelloText must not be set by NewSuccessEvent
	if evt.HelloText != "" {
		t.Errorf("HelloText = %q, want empty string", evt.HelloText)
	}
	// SessionID must be empty by default
	if evt.SessionID != "" {
		t.Errorf("SessionID = %q, want empty string", evt.SessionID)
	}
}

func TestNewErrorEvent(t *testing.T) {
	method := "GET"
	endpoint := "/api/data"
	errorMessage := "Connection timeout"

	evt := NewErrorEvent(method, endpoint, errorMessage)

	if evt.Method != method {
		t.Errorf("Method = %q, want %q", evt.Method, method)
	}
	if evt.Endpoint != endpoint {
		t.Errorf("Endpoint = %q, want %q", evt.Endpoint, endpoint)
	}
	if evt.Message != errorMessage {
		t.Errorf("Message = %q, want %q", evt.Message, errorMessage)
	}
	if evt.Result != "error" {
		t.Errorf("Result = %q, want %q", evt.Result, "error")
	}
	if evt.Level != LogError {
		t.Errorf("Level = %q, want %q", evt.Level, LogError)
	}
	if evt.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
	// HelloText must be empty for error events
	if evt.HelloText != "" {
		t.Errorf("HelloText = %q, want empty string", evt.HelloText)
	}
	if evt.SessionID != "" {
		t.Errorf("SessionID = %q, want empty string", evt.SessionID)
	}
}

func TestAuditEvent_JSONSerialization(t *testing.T) {
	// Create event with all fields set
	expected := AuditEvent{
		Timestamp: time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC),
		Method:    "GET",
		Endpoint:  "/hello",
		SessionID: "abc-123",
		Result:    "success",
		Level:     LogInfo,
		Message:   "Greeting sent successfully",
		HelloText: "Hello, World!",
	}

	// Marshal to JSON
	data, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Unmarshal back
	var actual AuditEvent
	err = json.Unmarshal(data, &actual)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Compare fields
	if !actual.Timestamp.Equal(expected.Timestamp) {
		t.Errorf("Timestamp = %v, want %v", actual.Timestamp, expected.Timestamp)
	}
	if actual.Method != expected.Method {
		t.Errorf("Method = %q, want %q", actual.Method, expected.Method)
	}
	if actual.Endpoint != expected.Endpoint {
		t.Errorf("Endpoint = %q, want %q", actual.Endpoint, expected.Endpoint)
	}
	if actual.SessionID != expected.SessionID {
		t.Errorf("SessionID = %q, want %q", actual.SessionID, expected.SessionID)
	}
	if actual.Result != expected.Result {
		t.Errorf("Result = %q, want %q", actual.Result, expected.Result)
	}
	if actual.Level != expected.Level {
		t.Errorf("Level = %q, want %q", actual.Level, expected.Level)
	}
	if actual.Message != expected.Message {
		t.Errorf("Message = %q, want %q", actual.Message, expected.Message)
	}
	if actual.HelloText != expected.HelloText {
		t.Errorf("HelloText = %q, want %q", actual.HelloText, expected.HelloText)
	}
}

func TestAuditEvent_JSON_Omitempty(t *testing.T) {
	// Event without SessionID and HelloText - these should be omitted
	evt := AuditEvent{
		Timestamp: time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC),
		Method:    "POST",
		Endpoint:  "/api/submit",
		Result:    "success",
		Level:     LogInfo,
		Message:   "Data submitted",
	}

	data, err := json.Marshal(evt)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Verify session_id and hello_text are omitted (not present in JSON)
	var raw map[string]interface{}
	err = json.Unmarshal(data, &raw)
	if err != nil {
		t.Fatalf("Unmarshal to raw map failed: %v", err)
	}

	if _, exists := raw["session_id"]; exists {
		t.Error("session_id should be omitted when empty")
	}
	if _, exists := raw["hello_text"]; exists {
		t.Error("hello_text should be omitted when empty")
	}

	// Check that non-empty fields are present
	for _, field := range []string{"timestamp", "method", "endpoint", "result", "level", "message"} {
		if _, exists := raw[field]; !exists {
			t.Errorf("field %q should be present in JSON", field)
		}
	}
}

func TestAuditEvent_JSON_WithEmptyHelloTextAndSessionID(t *testing.T) {
	evt := AuditEvent{
		Timestamp: time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC),
		Method:    "GET",
		Endpoint:  "/hello",
		Result:    "skipped",
		Level:     LogInfo,
		Message:   "Skipped greeting",
		SessionID: "",
		HelloText: "",
	}

	data, err := json.Marshal(evt)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var actual AuditEvent
	err = json.Unmarshal(data, &actual)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if actual.Result != "skipped" {
		t.Errorf("Result = %q, want %q", actual.Result, "skipped")
	}
	// After unmarshalling with omitempty, the fields should be empty strings
	if actual.SessionID != "" {
		t.Errorf("SessionID = %q, want empty", actual.SessionID)
	}
	if actual.HelloText != "" {
		t.Errorf("HelloText = %q, want empty", actual.HelloText)
	}
}

func TestLogLevelConstants(t *testing.T) {
	if LogInfo != "info" {
		t.Errorf("LogInfo = %q, want 'info'", LogInfo)
	}
	if LogLevelStr := LogError; LogLevelStr != "error" {
		t.Errorf("LogError = %q, want 'error'", LogLevelStr)
	}
}
