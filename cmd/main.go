package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Mark-Grigorev/FinGo/internal/config"
	"github.com/Mark-Grigorev/FinGo/internal/logger"
)

func main() {
	if code := run(); code != 0 {
		os.Exit(code)
	}
}

func run() int {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: %v\n", err)
		return 1
	}

	log := logger.SetupLogger(cfg.App.Env)
	log.Info("starting FinGo",
		slog.String("env", cfg.App.Env),
		slog.Int("port", cfg.App.Port),
		slog.String("db_host", cfg.DB.Host),
		slog.String("db_name", cfg.DB.Name),
	)

	// TODO: init db, router, server
	_ = cfg

	return 0
}
