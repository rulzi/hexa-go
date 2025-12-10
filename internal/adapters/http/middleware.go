package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
)

// AuthMiddleware creates a middleware for JWT authentication
func AuthMiddleware(userService *domainuser.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := userService.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}

// RecoveryMiddleware creates a middleware for panic recovery
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		c.Abort()
	})
}

// SetupDefaultMiddlewares applies default middlewares to the router
func SetupDefaultMiddlewares(engine *gin.Engine, debug bool) {
	// Recovery middleware - catches panics and returns proper error responses
	engine.Use(RecoveryMiddleware())

	// Logger middleware - logs HTTP requests (only in debug mode)
	if debug {
		engine.Use(gin.Logger())
	}

	// CORS middleware - handles cross-origin requests
	engine.Use(CORSMiddleware())
}

// CORSMiddleware creates a middleware for CORS handling
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
