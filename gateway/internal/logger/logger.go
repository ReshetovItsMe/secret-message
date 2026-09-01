// Package logger provides a standard library slog logger pre-configured
// for the gateway.
package logger

import (
	"log/slog"
	"os"
)

// New builds a slog logger. In pretty mode it uses a human-readable text
// handler (local dev); otherwise it emits structured JSON (production,
// container logs).
func New(pretty bool) *slog.Logger {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	var handler slog.Handler
	if pretty {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	return slog.New(handler)
}
