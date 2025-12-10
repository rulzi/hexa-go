package http

import (
	"github.com/gin-gonic/gin"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
)

// Router sets up the HTTP routes
type Router struct {
	userHandler    *UserHandler
	articleHandler *ArticleHandler
	userService    *domainuser.Service
}

// NewRouter creates a new router
func NewRouter(userHandler *UserHandler, articleHandler *ArticleHandler, userService *domainuser.Service) *Router {
	return &Router{
		userHandler:    userHandler,
		articleHandler: articleHandler,
		userService:    userService,
	}
}

// SetupRoutes configures all HTTP routes
func (r *Router) SetupRoutes(engine *gin.Engine, debug bool) {
	// Apply default middlewares
	SetupDefaultMiddlewares(engine, debug)

	api := engine.Group("/api/v1")
	{
		// Public routes (no authentication required)
		users := api.Group("/users")
		{
			users.POST("/register", r.userHandler.Register) // Register
			users.POST("/login", r.userHandler.Login)       // Login
		}

		// Protected routes (authentication required)
		authMiddleware := AuthMiddleware(r.userService)
		protected := api.Group("")
		protected.Use(authMiddleware)
		{
			usersProtected := protected.Group("/users")
			{
				usersProtected.GET("", r.userHandler.List)
				usersProtected.GET("/:id", r.userHandler.Get)
				usersProtected.PUT("/:id", r.userHandler.Update)
				usersProtected.DELETE("/:id", r.userHandler.Delete)
			}

			articlesProtected := protected.Group("/articles")
			{
				articlesProtected.POST("", r.articleHandler.Create)
				articlesProtected.GET("", r.articleHandler.List)
				articlesProtected.GET("/:id", r.articleHandler.Get)
				articlesProtected.PUT("/:id", r.articleHandler.Update)
				articlesProtected.DELETE("/:id", r.articleHandler.Delete)
			}
		}
	}

	// Health check endpoint
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
