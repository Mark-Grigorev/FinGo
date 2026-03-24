package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/Mark-Grigorev/FinGo/internal/service"
)

func NewRouter(
	log *slog.Logger,
	auth *service.AuthService,
	accounts *service.AccountService,
	transactions *service.TransactionService,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	authH := &authHandler{svc: auth, log: log}
	accH := &accountHandler{svc: accounts, log: log}
	txH := &transactionHandler{svc: transactions, log: log}

	api := r.Group("/api")
	{
		// Публичные маршруты
		api.POST("/auth/login", authH.login)
		api.POST("/auth/register", authH.register)
		api.POST("/auth/logout", authH.logout)

		// Защищённые маршруты
		protected := api.Group("/")
		protected.Use(authMiddleware(auth))
		{
			protected.GET("/auth/me", authH.me)

			protected.GET("/accounts", accH.list)
			protected.POST("/accounts", accH.create)
			protected.PUT("/accounts/:id", accH.update)
			protected.DELETE("/accounts/:id", accH.delete)

			protected.GET("/transactions", txH.list)
			protected.POST("/transactions", txH.create)
		}
	}

	return r
}
