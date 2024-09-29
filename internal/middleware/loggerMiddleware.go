package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func LoggerMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			method := r.Method
			uri := r.URL.Path
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			logger.Infoln("method", method, "uri", uri, "duration", duration)
		})
	}
}
