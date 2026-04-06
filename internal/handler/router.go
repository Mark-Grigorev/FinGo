package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/Mark-Grigorev/FinGo/internal/service"
)

// @title FinGo API
// @version 1.0
// @description Personal Finance Management API
// @termsOfService http://swagger.io/terms/

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
type Router struct {
	router *gin.Engine
	port   int
	log    *slog.Logger
}

type RouterCfg struct {
	Port      int
	Debug     bool
	UploadDir string
}

func NewRouter(
	log *slog.Logger,
	auth *service.AuthService,
	accounts *service.AccountService,
	transactions *service.TransactionService,
	dashboard *service.DashboardService,
	categories *service.CategoryService,
	budgets *service.BudgetService,
	recurring *service.RecurringService,
	currencies *service.CurrencyService,
	routerCfg RouterCfg,
) *Router {
	// if routerCfg.Debug {
	gin.SetMode(gin.DebugMode)
	// } else {
	// gin.SetMode(gin.ReleaseMode)
	// }

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	authH := &authHandler{svc: auth, log: log}
	accH := &accountHandler{svc: accounts, log: log}
	txH := &transactionHandler{svc: transactions, log: log}
	dashboardH := &dashboardHandler{svc: dashboard, log: log}
	catH := &categoryHandler{svc: categories, log: log}
	uploadH := &uploadHandler{uploadDir: routerCfg.UploadDir, log: log}
	budgetH := &budgetHandler{svc: budgets, log: log}
	recurringH := &recurringHandler{svc: recurring, log: log}
	currencyH := &currencyHandler{svc: currencies, log: log}

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Serve uploaded files
	r.Static("/api/uploads", routerCfg.UploadDir)

	api := r.Group("/api")
	{
		// Публичные маршруты
		api.POST("/auth/login", authH.login)
		api.POST("/auth/register", authH.register)
		api.POST("/auth/logout", authH.logout)
		api.POST("/auth/forgot-password", authH.forgotPassword)
		api.POST("/auth/reset-password", authH.resetPassword)

		// Защищённые маршруты
		protected := api.Group("/")
		protected.Use(authMiddleware(auth))
		{
			protected.GET("/auth/me", authH.me)
			protected.PUT("/user/profile", authH.updateProfile)
			protected.PUT("/user/password", authH.changePassword)

			protected.GET("/accounts", accH.list)
			protected.POST("/accounts", accH.create)
			protected.PUT("/accounts/:id", accH.update)
			protected.DELETE("/accounts/:id", accH.delete)

			protected.GET("/transactions/export", txH.export)
			protected.GET("/transactions", txH.list)
			protected.POST("/transactions", txH.create)
			protected.DELETE("/transactions/:id", txH.delete)

			protected.GET("/dashboard/summary", dashboardH.summary)
			protected.GET("/dashboard/report", dashboardH.report)

			protected.GET("/budgets", budgetH.list)
			protected.POST("/budgets", budgetH.create)
			protected.PUT("/budgets/:id", budgetH.update)
			protected.DELETE("/budgets/:id", budgetH.delete)

			protected.GET("/recurring", recurringH.list)
			protected.POST("/recurring", recurringH.create)
			protected.PUT("/recurring/:id", recurringH.update)
			protected.DELETE("/recurring/:id", recurringH.delete)

			protected.GET("/categories", catH.list)
			protected.POST("/categories", catH.create)
			protected.PUT("/categories/:id", catH.update)
			protected.DELETE("/categories/:id", catH.delete)

			protected.POST("/categories/icons", uploadH.uploadIcon)

			protected.GET("/currencies/rates", currencyH.listRates)
			protected.PUT("/currencies/rates/:currency", currencyH.upsertRate)
			protected.DELETE("/currencies/rates/:currency", currencyH.deleteRate)
			protected.GET("/currencies/base", currencyH.getBaseCurrency)
			protected.PUT("/currencies/base", currencyH.setBaseCurrency)
		}
	}

	return &Router{
		router: r,
		port:   routerCfg.Port,
		log:    log,
	}
}

func (r *Router) newServer() *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", r.port),
		Handler:      r.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func (r *Router) Serve(srv *http.Server) {
	r.log.Info("server listening", slog.Int("port", r.port))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		r.log.Error("server error", "err", err)
	}
}

func (r *Router) GracefulShutdown(srv *http.Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	r.log.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		r.log.Error("graceful shutdown failed", "err", err)
		return err
	}
	return nil
}

func (r *Router) Start() error {
	srv := r.newServer()
	go r.Serve(srv)
	return r.GracefulShutdown(srv)
}
