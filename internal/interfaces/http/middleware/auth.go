package middleware

import (
	"net/http"
	"strings"

	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/infrastructure/auth"
	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/interfaces/http/response"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token
func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse("missing authorization header"))
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse("invalid authorization header format"))
			c.Abort()
			return
		}

		token := parts[1]

		// Verify token
		claims, err := jwtManager.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse("invalid or expired token"))
			c.Abort()
			return
		}

		// Store claims in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
