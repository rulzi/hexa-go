package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rulzi/hexa-go/internal/adapters/http/response"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
)

// RecoveryMiddleware creates a middleware for panic recovery
func RecoveryMiddleware(appLogger logger.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		appLogger.ErrorWithFields("panic recovered", map[string]interface{}{
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
			"panic":  fmt.Sprintf("%v", recovered),
		})
		response.ErrorResponseInternalServerError(c, "internal server error")
		c.Abort()
	})
}
