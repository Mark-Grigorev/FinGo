package logger

import (
	"log/slog"
	"os"
)

// SetupLogger возвращает slog.Logger под нужный env.
func SetupLogger(env string) *slog.Logger {
	var handler slog.Handler

	switch env {
	case "prod":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	case "dev":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	default: // local
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	return slog.New(handler)
}
