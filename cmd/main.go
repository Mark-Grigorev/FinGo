package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
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

	// Upload directory
	uploadDir := "/app/uploads/icons"
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		log.Error("create upload dir", "err", err)
		return exitConfigError
	}

	// Services & router
	authSvc := service.NewAuth(store, maker, log)
	accountSvc := service.NewAccount(store, log)
	txSvc := service.NewTransaction(store, log)
	dashboardSvc := service.NewDashboard(store, log)
	categorySvc := service.NewCategory(store, log)
	budgetSvc := service.NewBudget(store, log)
	recurringSvc := service.NewRecurring(store, log)

	err = handler.NewRouter(
		log,
		authSvc,
		accountSvc,
		txSvc,
		dashboardSvc,
		categorySvc,
		budgetSvc,
		recurringSvc,
		handler.RouterCfg{
			Port:      cfg.App.Port,
			Debug:     cfg.App.Debug,
			UploadDir: uploadDir,
		}).Start()
	if err != nil {
		log.Error("server shutdown error", "err", err)
		return exitServerShutdown
	}

	log.Info("server stopped")
	return exitOK
}
