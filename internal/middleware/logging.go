package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		duration := time.Since(start)

		fmt.Printf(
			"method=%s path=%s duration=%s\n",
			r.Method,
			r.URL.Path,
			duration,
		)
	})
}
