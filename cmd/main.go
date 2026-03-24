package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mark-Grigorev/FinGo/internal/config"
	"github.com/Mark-Grigorev/FinGo/internal/handler"
	"github.com/Mark-Grigorev/FinGo/internal/logger"
	"github.com/Mark-Grigorev/FinGo/internal/repository"
	"github.com/Mark-Grigorev/FinGo/internal/service"
	"github.com/Mark-Grigorev/FinGo/pkg/token"
)

// Exit codes — каждый сценарий ошибки имеет уникальный код.
// Описание кодов: см. README.md → "Exit Codes / Коды завершения".
const (
	exitOK             = 0
	exitConfigError    = 2
	exitTokenError     = 3
	exitDBConnect      = 4
	exitDBMigrate      = 5
	exitServerShutdown = 6
)

func main() {
	os.Exit(run())
}

func run() int {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "config:", err)
		return exitConfigError
	}

	log := logger.SetupLogger(cfg.App.Debug)
	log.Info("starting FinGo",
		slog.String("env", cfg.App.Env),
		slog.Int("port", cfg.App.Port),
	)

	// Token maker
	duration, err := time.ParseDuration(cfg.Token.Duration)
	if err != nil {
		duration = 24 * time.Hour
	}
	maker, err := token.New(cfg.Token.SymmetricKey, duration)
	if err != nil {
		log.Error("token maker init failed", "err", err)
		return exitTokenError
	}

	// Database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	store, err := repository.New(ctx, cfg.DB)
	if err != nil {
		log.Error("database connection failed", "err", err)
		return exitDBConnect
	}
	defer store.Close()
	log.Info("connected to database")

	if err := store.Migrate(); err != nil {
		log.Error("migration failed", "err", err)
		return exitDBMigrate
	}
	log.Info("migrations applied")

	// Services & router
	authSvc := service.NewAuth(store, maker, log)
	accountSvc := service.NewAccount(store, log)
	txSvc := service.NewTransaction(store, log)

	router := handler.NewRouter(log, authSvc, accountSvc, txSvc)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info("server listening", slog.Int("port", cfg.App.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", "err", err)
		}
	}()

	<-quit
	log.Info("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", "err", err)
		return exitServerShutdown
	}

	log.Info("server stopped")
	return exitOK
}
