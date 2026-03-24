package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Mark-Grigorev/FinGo/internal/service"
)

const ctxUserKey = "current_user_id"

func authMiddleware(auth *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Попробуем сначала cookie, затем Authorization header
		tokenStr := ""

		if cookie, err := c.Cookie("session_token"); err == nil {
			tokenStr = cookie
		} else if header := c.GetHeader("Authorization"); strings.HasPrefix(header, "Bearer ") {
			tokenStr = strings.TrimPrefix(header, "Bearer ")
		}

		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "необходима авторизация"})
			return
		}

		payload, err := auth.VerifyToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "необходима авторизация"})
			return
		}

		c.Set(ctxUserKey, payload.UserID)
		c.Next()
	}
}

func currentUserID(c *gin.Context) int64 {
	id, _ := c.Get(ctxUserKey)
	userID, _ := id.(int64)
	return userID
}
