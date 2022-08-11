package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/app"
)

func loggingMiddleware(l app.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			l.Info(
				fmt.Sprintf(
					"Got request: ip %s, start time %s, duration %s, user-agent: %s, method: %s, path: %s, version: %s",
					r.RemoteAddr,
					start,
					duration,
					r.Header.Get("user-agent"),
					r.Method,
					r.RequestURI,
					r.Proto,
				),
			)
		})
	}
}
