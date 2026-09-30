package middleware

import (
	"fmt"
	"net/http"
	"time"

	"guess-number/internal/logger"
)

// LoggingMiddleware перехватывает HTTP-запросы и логирует их через переданный Logger
func LoggingMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Выполнение следующего обработчика
			next.ServeHTTP(w, r)

			// Логирование завершенного запроса
			logMsg := fmt.Sprintf("HTTP Request: method=%s path=%s duration=%s", r.Method, r.URL.Path, time.Since(start))
			log.Log(logMsg)
		})
	}
}
