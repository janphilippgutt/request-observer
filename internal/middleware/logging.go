package middleware

import (
	"fmt"
	"net/http"
	"time"
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

		fmt.Printf(
			"method=%s path=%s status=%d duration=%s\n",
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			duration,
		)
	})
}
