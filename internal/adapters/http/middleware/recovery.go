package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rulzi/hexa-go/internal/adapters/contextkey"
	"github.com/rulzi/hexa-go/internal/adapters/http/response"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
)

// RecoveryMiddleware creates a middleware for panic recovery
func RecoveryMiddleware(appLogger logger.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		fields := map[string]interface{}{
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
			"panic":  fmt.Sprintf("%v", recovered),
		}
		if requestID := contextkey.RequestID(c.Request.Context()); requestID != "" {
			fields["request_id"] = requestID
		}
		appLogger.ErrorWithFields("panic recovered", fields)
		response.ErrorResponseInternalServerError(c, "internal server error")
		c.Abort()
	})
}
