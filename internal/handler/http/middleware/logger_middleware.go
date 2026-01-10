package middleware

import (
	"net/http"
	"time"

	"github.com/Soujuruya/01_SPEC/internal/pkg/logger"
	"github.com/google/uuid"
)

func Logger(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			traceID := uuid.NewString()

			log.Info("request started",
				"trace_id", traceID,
				"method", r.Method,
				"path", r.URL.Path,
				"remote", r.RemoteAddr,
			)

			next.ServeHTTP(w, r)

			log.Info("request finished",
				"trace_id", traceID,
				"method", r.Method,
				"path", r.URL.Path,
				"duration_ms", time.Since(start).Milliseconds(),
			)
		})
	}
}
