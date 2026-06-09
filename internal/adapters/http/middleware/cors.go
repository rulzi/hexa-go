package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	corsAllowHeaders = "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"
	corsAllowMethods = "POST, OPTIONS, GET, PUT, DELETE, PATCH"
)

// CORSConfig holds CORS middleware configuration.
type CORSConfig struct {
	AllowedOrigins []string
}

// CORSMiddleware creates a middleware for CORS handling.
// When AllowedOrigins contains "*", all origins are permitted without credentials.
// Otherwise only listed origins are reflected and credentials are enabled.
func CORSMiddleware(cfg CORSConfig) gin.HandlerFunc {
	allowAll := false
	allowed := make(map[string]struct{}, len(cfg.AllowedOrigins))

	for _, origin := range cfg.AllowedOrigins {
		origin = strings.TrimSpace(origin)
		if origin == "" {
			continue
		}
		if origin == "*" {
			allowAll = true
			break
		}
		allowed[origin] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if allowAll {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else if origin != "" {
			if _, ok := allowed[origin]; ok {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
				c.Writer.Header().Set("Vary", "Origin")
			}
		}

		if allowAll || c.Writer.Header().Get("Access-Control-Allow-Origin") != "" {
			c.Writer.Header().Set("Access-Control-Allow-Headers", corsAllowHeaders)
			c.Writer.Header().Set("Access-Control-Allow-Methods", corsAllowMethods)
		}

		if c.Request.Method == http.MethodOptions {
			if c.Writer.Header().Get("Access-Control-Allow-Origin") != "" {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
