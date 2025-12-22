package middleware

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/janphilippgutt/request-observer/internal/logging"
	"github.com/janphilippgutt/request-observer/internal/observability"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code                // record status code
	r.ResponseWriter.WriteHeader(code) // call the real writer
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode:     200, // default to 200 OK: If the handler never calls WriteHeader, Go assumes 200
		}

		// Call the next handler
		next.ServeHTTP(recorder, r)

		duration := time.Since(start)

		reqID := FromContext(r.Context())

		if recorder.statusCode >= 500 {
			logging.Logger.Error(
				"http request failed",
				slog.String("request_id", reqID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", recorder.statusCode),
				slog.Duration("duration", duration),
			)
		} else {
			logging.Logger.Info(
				"http request completed",
				slog.String("request_id", reqID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", recorder.statusCode),
				slog.Duration("duration", duration),
			)
		}

		observability.HTTPRequestsTotal.WithLabelValues(
			r.Method,
			r.URL.Path,
			strconv.Itoa(recorder.statusCode),
		).Inc()

		observability.HTTPRequestDuration.WithLabelValues(
			r.Method,
			r.URL.Path,
		).Observe(time.Since(start).Seconds())

	})
}
