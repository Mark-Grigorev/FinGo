package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Mark-Grigorev/FinGo/internal/service"
)

type categoryHandler struct {
	svc *service.CategoryService
	log *slog.Logger
}

func (h *categoryHandler) list(c *gin.Context) {
	userID := currentUserID(c)
	categories, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": categories})
}

func (h *categoryHandler) create(c *gin.Context) {
	var in struct {
		Name  string `json:"name"  binding:"required"`
		Icon  string `json:"icon"`
		Color string `json:"color"`
		Type  string `json:"type"  binding:"required,oneof=income expense"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный формат запроса"})
		return
	}

	userID := currentUserID(c)
	cat, err := h.svc.Create(c.Request.Context(), userID, in.Name, in.Icon, in.Color, in.Type)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func (h *categoryHandler) update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный id"})
		return
	}

	var in struct {
		Name  string `json:"name"  binding:"required"`
		Icon  string `json:"icon"`
		Color string `json:"color"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный формат запроса"})
		return
	}

	userID := currentUserID(c)
	cat, err := h.svc.Update(c.Request.Context(), id, userID, in.Name, in.Icon, in.Color)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (h *categoryHandler) delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный id"})
		return
	}

	userID := currentUserID(c)
	if err := h.svc.Delete(c.Request.Context(), id, userID); err != nil {
		writeError(c, h.log, err)
		return
	}
	c.Status(http.StatusNoContent)
}
