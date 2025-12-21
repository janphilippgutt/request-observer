package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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

		fmt.Printf(
			"request_id%s method=%s path=%s status=%d duration=%s\n",
			reqID,
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			duration,
		)

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
