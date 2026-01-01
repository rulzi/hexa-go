package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rulzi/hexa-go/internal/adapters/http/response"
)

// RecoveryMiddleware creates a middleware for panic recovery
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		response.ErrorResponseInternalServerError(c, "internal server error")
		c.Abort()
	})
}
