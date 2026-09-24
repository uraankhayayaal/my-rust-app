package logger

import (
	"context"
	"sync"
	"testing"
	"time"

	"backend/domain"
)

func TestAppendAndGetLogs(t *testing.T) {
	store := NewInMemoryLogger(100)
	ctx := context.Background()

	event := domain.AuditEvent{
		ID:        "1",
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EventType: "login",
		Message:   "User logged in",
		UserID:    "user-1",
	}

	store.Append(ctx, event)

	logs, err := store.GetLogs(-1)
	if err != nil {
		t.Fatalf("GetLogs returned error: %v", err)
	}
	if logs == nil {
		t.Fatal("expected non-nil logs")
	}
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}
	if logs[0].ID != event.ID {
		t.Errorf("expected event ID %q, got %q", event.ID, logs[0].ID)
	}
}

func TestGetLogsEmpty(t *testing.T) {
	store := NewInMemoryLogger(100)

	logs, err := store.GetLogs(-1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if logs != nil {
		t.Fatalf("expected nil logs for empty store, got %v", logs)
	}
}

func TestCapTrimming(t *testing.T) {
	store := NewInMemoryLogger(3)
	ctx := context.Background()

	for i := 1; i <= 5; i++ {
		event := domain.AuditEvent{
			ID:        string(rune('A' + i - 1)),
			Timestamp: time.Now(),
			EventType: "test",
			Message:   "event", 
		}
		store.Append(ctx, event)
	}

	logs, err := store.GetLogs(-1)
	if err != nil {
		t.Fatalf("GetLogs returned error: %v", err)
	}
	if len(logs) != 3 {
		t.Fatalf("expected 3 logs after overflow, got %d", len(logs))
	}
	// oldest должны быть удалены — остались события с ID C, D, E
	if logs[0].ID != "C" {
		t.Errorf("expected first remaining event 'C', got %q", logs[0].ID)
	}
	if logs[2].ID != "E" {
		t.Errorf("expected last remaining event 'E', got %q", logs[2].ID)
	}
}

func TestGetLogsWithLimit(t *testing.T) {
	store := NewInMemoryLogger(10)
	ctx := context.Background()

	for i := 1; i <= 10; i++ {
		event := domain.AuditEvent{
			ID:        string(rune('A' + i - 1)),
			Timestamp: time.Now(),
			EventType: "test",
			Message:   "event",
		}
		store.Append(ctx, event)
	}

	logs, err := store.GetLogs(3)
	if err != nil {
		t.Fatalf("GetLogs returned error: %v", err)
	}
	if len(logs) != 3 {
		t.Fatalf("expected 3 logs, got %d", len(logs))
	}
	// последние 3: H I J
	if logs[0].ID != "H" || logs[1].ID != "I" || logs[2].ID != "J" {
		t.Errorf("expected [H I J], got %v", logs)
	}

	logs, err = store.GetLogs(0)
	if err != nil {
		t.Fatalf("GetLogs with limit=0 returned error: %v", err)
	}
	if len(logs) != 10 {
		t.Fatalf("expected all 10 logs, got %d", len(logs))
	}
}

func TestConcurrentAppend(t *testing.T) {
	store := NewInMemoryLogger(10000)
	ctx := context.Background()

	var wg sync.WaitGroup
	numGoroutines := 100
	eventsPerGoroutine := 100

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for i := 0; i < eventsPerGoroutine; i++ {
				event := domain.AuditEvent{
					ID:        string(rune('A' + (base+i)%26)),
					Timestamp: time.Now(),
					EventType: "concurrent",
					Message:   "goroutine", 
				}
				store.Append(ctx, event)
			}
		}(g)
	}

	wg.Wait()

	logs, err := store.GetLogs(-1)
	if err != nil {
		t.Fatalf("GetLogs returned error: %v", err)
	}
	if len(logs) != 10000 {
		t.Fatalf("expected 10000 logs, got %d", len(logs))
	}
}
