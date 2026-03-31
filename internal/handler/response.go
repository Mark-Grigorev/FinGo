package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
)

func writeError(c *gin.Context, log *slog.Logger, err error) {
	var status int
	var message string

	switch {
	case errors.Is(err, domain.ErrNotFound):
		status, message = http.StatusNotFound, "не найдено"
	case errors.Is(err, domain.ErrUnauthorized):
		status, message = http.StatusUnauthorized, "необходима авторизация"
	case errors.Is(err, domain.ErrForbidden):
		status, message = http.StatusForbidden, "доступ запрещён"
	case errors.Is(err, domain.ErrAlreadyExists):
		status, message = http.StatusConflict, "уже существует"
	case errors.Is(err, domain.ErrInvalidInput):
		status, message = http.StatusBadRequest, "неверные данные"
	default:
		log.Error("internal error", "err", err)
		status, message = http.StatusInternalServerError, "внутренняя ошибка"
	}

	c.JSON(status, gin.H{"message": message})
}
