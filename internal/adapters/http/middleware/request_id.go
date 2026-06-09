package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"github.com/rulzi/hexa-go/internal/adapters/contextkey"
)

// RequestIDMiddleware injects a request ID into the Gin context and request context.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(contextkey.RequestIDHeader)
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Set("request_id", requestID)
		c.Header(contextkey.RequestIDHeader, requestID)
		c.Request = c.Request.WithContext(contextkey.WithRequestID(c.Request.Context(), requestID))

		c.Next()
	}
}

func generateRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}
