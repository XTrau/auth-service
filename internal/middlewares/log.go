package middlewares

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{w, 200}
		start := time.Now()

		next.ServeHTTP(rw, r)

		elapsed := time.Since(start)
		slog.Info(
			"request",
			slog.String("Method", r.Method),
			slog.String("Path", r.URL.Path),
			slog.Int("Status", rw.statusCode),
			slog.Time("RequestTime", start),
			slog.Duration("Duration", elapsed),
		)
	})
}
