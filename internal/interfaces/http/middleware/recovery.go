package middleware

import (
	"log"
	"net/http"

	"github.com/ardianhermawan17/warehouse-management-system-ddd/internal/interfaces/http/response"
	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.JSON(http.StatusInternalServerError, response.ErrorResponse("internal server error"))
				c.Abort()
			}
		}()
		c.Next()
	}
}
