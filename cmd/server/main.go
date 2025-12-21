package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/janphilippgutt/request-observer/internal/httpapi"
	"github.com/janphilippgutt/request-observer/internal/middleware"
	"github.com/janphilippgutt/request-observer/internal/observability"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {

	handler := middleware.RequestID(
		middleware.Logging(
			httpapi.NewRouter(),
		),
	)

	// Register metrics
	prometheus.MustRegister(
		observability.HTTPRequestsTotal,
		observability.HTTPRequestDuration,
	)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	// Start server in a goroutine
	go func() {
		fmt.Println("Server running on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Server error:", err)
		}
	}()

	// Listen for shutdown signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	fmt.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Shutdown error:", err)
	}

	fmt.Println("Server stopped")
}
