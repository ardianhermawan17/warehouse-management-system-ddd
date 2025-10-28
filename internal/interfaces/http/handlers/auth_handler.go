package handlers

import (
	"net/http"

	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/infrastructure/auth"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/interfaces/http/response"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	jwtManager *auth.JWTManager
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		jwtManager: jwtManager,
	}
}

// LoginRequest is the DTO for login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is the DTO for login response
type LoginResponse struct {
	Token string `json:"token"`
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("invalid request"))
		return
	}

	// Simple authentication (in production, validate against database)
	// For demo purposes, accept any username/password
	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("invalid credentials"))
		return
	}

	// Generate token (user ID is hardcoded for demo)
	token, err := h.jwtManager.GenerateToken(1, req.Username, 24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("failed to generate token"))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("login successful", &LoginResponse{
		Token: token,
	}))
}
