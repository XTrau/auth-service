package middlewares

import (
	"log/slog"
	"net/http"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic on http handler", slog.Any("err", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
