package logging

import (
	"log/slog"
	"os"
)

// package-level variable - one logger per service, no reconfiguration per request
var Logger *slog.Logger

func Init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	// make the handler globally accessible
	Logger = slog.New(handler)
}
