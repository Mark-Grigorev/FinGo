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

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com" binding:"required,email"`
	Password string `json:"password" example:"password123" binding:"required"`
}

// RegisterRequest represents registration data
type RegisterRequest struct {
	Email    string `json:"email" example:"user@example.com" binding:"required,email"`
	Name     string `json:"name" example:"John Doe" binding:"required"`
	Password string `json:"password" example:"password123" binding:"required,min=6"`
}

// UpdateProfileRequest represents profile update data
type UpdateProfileRequest struct {
	Name  string `json:"name" example:"John Doe" binding:"required"`
	Email string `json:"email" example:"user@example.com" binding:"required,email"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" example:"oldpass123" binding:"required"`
	NewPassword string `json:"new_password" example:"newpass123" binding:"required,min=6"`
}

// UserResponse represents user data response
type UserResponse struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email     string    `json:"email" example:"user@example.com"`
	Name      string    `json:"name" example:"John Doe"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Message string `json:"message" example:"invalid credentials"`
}

// login godoc
// @Summary Login user
// @Description Authenticate user with email and password, returns user data and sets session cookie
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} UserResponse "User authenticated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *authHandler) login(c *gin.Context) {
	var in LoginRequest
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

// register godoc
// @Summary Register new user
// @Description Create a new user account, returns user data and sets session cookie
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration data"
// @Success 201 {object} UserResponse "User created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request format or validation failed"
// @Failure 409 {object} ErrorResponse "Email already exists"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/register [post]
func (h *authHandler) register(c *gin.Context) {
	var in RegisterRequest
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

// logout godoc
// @Summary Logout user
// @Description Clear session cookie and logout user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 204 "Successfully logged out"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /auth/logout [post]
func (h *authHandler) logout(c *gin.Context) {
	clearTokenCookie(c)
	c.Status(http.StatusNoContent)
}

// me godoc
// @Summary Get current user info
// @Description Returns information about the authenticated user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserResponse "User information retrieved successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/me [get]
func (h *authHandler) me(c *gin.Context) {
	userID := currentUserID(c)
	user, err := h.svc.GetUser(c.Request.Context(), userID)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// updateProfile godoc
// @Summary Update user profile
// @Description Update authenticated user's name and email
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "Profile update data"
// @Success 200 {object} UserResponse "Profile updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 409 {object} ErrorResponse "Email already exists"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /user/profile [put]
func (h *authHandler) updateProfile(c *gin.Context) {
	userID := currentUserID(c)
	var in UpdateProfileRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный формат запроса"})
		return
	}
	user, err := h.svc.UpdateProfile(c.Request.Context(), userID, in.Name, in.Email)
	if err != nil {
		writeError(c, h.log, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// changePassword godoc
// @Summary Change user password
// @Description Change authenticated user's password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Password change data"
// @Success 204 "Password changed successfully"
// @Failure 400 {object} ErrorResponse "Invalid request format or password too short"
// @Failure 401 {object} ErrorResponse "Unauthorized or invalid old password"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /user/password [put]
func (h *authHandler) changePassword(c *gin.Context) {
	userID := currentUserID(c)
	var in ChangePasswordRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "новый пароль должен содержать минимум 6 символов"})
		return
	}
	if err := h.svc.ChangePassword(c.Request.Context(), userID, in.OldPassword, in.NewPassword); err != nil {
		writeError(c, h.log, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func setTokenCookie(c *gin.Context, tokenStr string, expires time.Time) {
	maxAge := int(time.Until(expires).Seconds())
	c.SetCookie("session_token", tokenStr, maxAge, "/", "", false, true)
}

func clearTokenCookie(c *gin.Context) {
	c.SetCookie("session_token", "", -1, "/", "", false, true)
}
