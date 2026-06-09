package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rulzi/hexa-go/internal/adapters/contextkey"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
)

// LoggingMiddleware logs each HTTP request with structured fields.
func LoggingMiddleware(appLogger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		fields := map[string]interface{}{
			"method":     c.Request.Method,
			"path":       path,
			"status":     c.Writer.Status(),
			"latency_ms": latency.Milliseconds(),
			"client_ip":  c.ClientIP(),
		}

		if requestID := contextkey.RequestID(c.Request.Context()); requestID != "" {
			fields["request_id"] = requestID
		}

		status := c.Writer.Status()
		switch {
		case status >= 500:
			appLogger.ErrorWithFields("http request completed", fields)
		case status >= 400:
			appLogger.WarnWithFields("http request completed", fields)
		default:
			appLogger.InfoWithFields("http request completed", fields)
		}
	}
}
