package logger

import (
	"log/slog"
	"os"
)

// SetupLogger возвращает slog.Logger под нужный env.
func SetupLogger(debug bool) *slog.Logger {
	var handler slog.Handler

	switch debug {
	case false:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	return slog.New(handler)
}
