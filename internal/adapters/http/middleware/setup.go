package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
)

// SetupDefaultMiddlewares applies default middlewares to the router
func SetupDefaultMiddlewares(engine *gin.Engine, debug bool, appLogger logger.Logger) {
	// Recovery middleware - catches panics and returns proper error responses
	engine.Use(RecoveryMiddleware(appLogger))

	// Logger middleware - logs HTTP requests (only in debug mode)
	if debug {
		engine.Use(gin.Logger())
	}

	// CORS middleware - handles cross-origin requests
	engine.Use(CORSMiddleware())
}
