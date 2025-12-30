package httpapi

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)

	mux.Handle("/metrics", promhttp.Handler())

	// http.HandlerFunc turns a function into an object that speaks HTTP
	mux.Handle("/observe", http.HandlerFunc(ObserveHandler))

	return mux
}
