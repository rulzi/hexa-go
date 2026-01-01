package middleware

import (
	"github.com/gin-gonic/gin"
)

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
