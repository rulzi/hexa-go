package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
)

// SetupDefaultMiddlewares applies default middlewares to the router
func SetupDefaultMiddlewares(engine *gin.Engine, appLogger logger.Logger, corsCfg CORSConfig) {
	// Recovery middleware - catches panics and returns proper error responses
	engine.Use(RecoveryMiddleware(appLogger))

	// Request ID middleware - propagates request_id through context
	engine.Use(RequestIDMiddleware())

	// Structured logging middleware - logs every HTTP request
	engine.Use(LoggingMiddleware(appLogger))

	// CORS middleware - handles cross-origin requests
	engine.Use(CORSMiddleware(corsCfg))
}
