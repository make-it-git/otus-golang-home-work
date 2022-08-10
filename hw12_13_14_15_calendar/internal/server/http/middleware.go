package internalhttp

import (
	"net/http"
)

func loggingMiddleware(l Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Info("test")
			next.ServeHTTP(w, r)
		})
	}
}
