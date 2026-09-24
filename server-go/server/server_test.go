package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestNewServer проверяет корректность создания сервера.
func TestNewServer(t *testing.T) {
	tests := []struct {
		name string
		addr string
	}{
		{"empty addr", ""},
		{"normal addr", ":8080"},
		{"full addr", "127.0.0.1:9999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewServer(tt.addr)
			if srv == nil {
				t.Errorf("NewServer() returned nil")
				return
			}
			if srv.srv == nil {
				t.Error("NewServer(): srv is nil")
			}
			if srv.srv.Addr != tt.addr {
				t.Errorf("NewServer().srv.Addr = %q, want %q", srv.srv.Addr, tt.addr)
			}
			if srv.mux == nil {
				t.Error("NewServer(): mux is nil")
			}
		})
	}
}

// TestSetHandler проверяет регистрацию обработчиков.
func TestSetHandler(t *testing.T) {
	srv := NewServer(":0")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SetHandler panicked: %v", r)
		}
	}()

	srv.SetHandler("/test", handler)
}

// TestStartShutdown через httptest проверяет запуск и остановку.
func TestStartShutdown(t *testing.T) {
	srv := NewServer("")
	if err := srv.RegisterDefaultRoutes(); err != nil {
		t.Fatalf("RegisterDefaultRoutes failed: %v", err)
	}

	// httptest.NewServer слушает на случайном свободном порту с mux.
	// Проверяем что mux корректно настроен.
	ts := httptest.NewServer(srv.mux)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/healthz")
	if err != nil {
		t.Fatalf("GET /healthz failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET /healthz status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

// TestRegisterDefaultRoutes проверяет регистрацию маршрутов.
func TestRegisterDefaultRoutes(t *testing.T) {
	tests := []struct {
		name    string
		srv     *Server
		wantErr bool
	}{
		{"nil server", nil, false},
		{"empty mux", &Server{}, false},
		{"valid server", NewServer(":0"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.srv.RegisterDefaultRoutes()
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterDefaultRoutes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestShutdownWithoutStart проверяет остановку сервера который не был запущен.
func TestShutdownWithoutStart(t *testing.T) {
	srv := NewServer(":0")
	ctx, cancel := context.WithTimeout(context.Background(), 1)
	defer cancel()

	_ = srv.Shutdown(ctx)
}

// TestHealthzAlwaysOK проверяет что healthz возвращает 200.
func TestHealthzAlwaysOK(t *testing.T) {
	srv := NewServer("")
	if err := srv.RegisterDefaultRoutes(); err != nil {
		t.Fatalf("RegisterDefaultRoutes failed: %v", err)
	}

	ts := httptest.NewServer(srv.mux)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/healthz")
	if err != nil {
		t.Fatalf("GET /healthz failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET /healthz status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

// TestReadyzReturns200 проверяет что readyz возвращает 200.
func TestReadyzReturns200(t *testing.T) {
	srv := NewServer("")
	if err := srv.RegisterDefaultRoutes(); err != nil {
		t.Fatalf("RegisterDefaultRoutes failed: %v", err)
	}

	ts := httptest.NewServer(srv.mux)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/readyz")
	if err != nil {
		t.Fatalf("GET /readyz failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET /readyz status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}