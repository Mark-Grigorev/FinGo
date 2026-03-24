package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Mark-Grigorev/FinGo/internal/service"
)

type authHandler struct {
	svc *service.AuthService
	log *slog.Logger
}

func (h *authHandler) login(c *gin.Context) {
	var in struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный формат запроса"})
		return
	}

	user, tokenStr, payload, err := h.svc.Login(c.Request.Context(), in.Email, in.Password)
	if err != nil {
		writeError(c, h.log, err)
		return
	}

	setTokenCookie(c, tokenStr, payload.ExpiredAt)
	c.JSON(http.StatusOK, user)
}

func (h *authHandler) register(c *gin.Context) {
	var in struct {
		Email    string `json:"email" binding:"required,email"`
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный формат запроса"})
		return
	}

	user, tokenStr, payload, err := h.svc.Register(c.Request.Context(), in.Email, in.Name, in.Password)
	if err != nil {
		writeError(c, h.log, err)
		return
	}

	setTokenCookie(c, tokenStr, payload.ExpiredAt)
	c.JSON(http.StatusCreated, user)
}

func (h *authHandler) logout(c *gin.Context) {
	clearTokenCookie(c)
	c.Status(http.StatusNoContent)
}

func (h *authHandler) me(c *gin.Context) {
	userID := currentUserID(c)
	user, err := h.svc.GetUser(c.Request.Context(), userID)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func setTokenCookie(c *gin.Context, tokenStr string, expires time.Time) {
	maxAge := int(time.Until(expires).Seconds())
	c.SetCookie("session_token", tokenStr, maxAge, "/", "", false, true)
}

func clearTokenCookie(c *gin.Context) {
	c.SetCookie("session_token", "", -1, "/", "", false, true)
}
