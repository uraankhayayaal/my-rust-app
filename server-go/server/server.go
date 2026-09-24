package server

import (
	"context"
	"net/http"
)

// Server — базовая структура HTTP-сервиса.
type Server struct {
	mux *http.ServeMux
	srv *http.Server
}

// NewServer создаёт новый сервер с указанным адресом.
func NewServer(addr string) *Server {
	return &Server{
		mux: http.NewServeMux(),
		srv: &http.Server{
			Addr: addr,
		},
	}
}

// SetHandler регистрирует handler по указанному пути.
func (s *Server) SetHandler(path string, handler http.Handler) {
	s.mux.Handle(path, handler)
}

// Start запускает HTTP-сервер и блокируется до остановки.
func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

// Shutdown gracefully останавливает сервер.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

// healthzHandler возвращает 200 OK при любом состоянии.
type healthzHandler struct{}

func (h *healthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// readyzHandler возвращает 200 OK когда сервер готов принимать запросы.
type readyzHandler struct{}

func (h *readyzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// RegisterDefaultRoutes регистрирует стандартные маршруты: /healthz и /readyz.
func (s *Server) RegisterDefaultRoutes() error {
	if s == nil || s.mux == nil {
		return nil
	}
	s.SetHandler("/healthz", &healthzHandler{})
	s.SetHandler("/readyz", &readyzHandler{})
	return nil
}