package http

import (
	"github.com/gin-gonic/gin"
	httparticle "github.com/rulzi/hexa-go/internal/adapters/http/article"
	httpmedia "github.com/rulzi/hexa-go/internal/adapters/http/media"
	"github.com/rulzi/hexa-go/internal/adapters/http/middleware"
	"github.com/rulzi/hexa-go/internal/adapters/http/response"
	httpuser "github.com/rulzi/hexa-go/internal/adapters/http/user"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
)

// Router sets up the HTTP routes
type Router struct {
	userHandler     *httpuser.Handler
	articleHandler  *httparticle.Handler
	mediaHandler    *httpmedia.Handler
	userService     *domainuser.Service
	storageBasePath string
}

// NewRouter creates a new router
func NewRouter(userHandler *httpuser.Handler, articleHandler *httparticle.Handler, mediaHandler *httpmedia.Handler, userService *domainuser.Service, storageBasePath string) *Router {
	return &Router{
		userHandler:     userHandler,
		articleHandler:  articleHandler,
		mediaHandler:    mediaHandler,
		userService:     userService,
		storageBasePath: storageBasePath,
	}
}

// SetupRoutes configures all HTTP routes
func (r *Router) SetupRoutes(engine *gin.Engine, debug bool) {
	// Apply default middlewares
	middleware.SetupDefaultMiddlewares(engine, debug)

	api := engine.Group("/api/v1")
	{
		// Public routes (no authentication required)
		// Media files endpoint (public access)
		api.StaticFS("/media/files", gin.Dir(r.storageBasePath, false))

		users := api.Group("/users")
		{
			users.POST("/register", r.userHandler.Register) // Register
			users.POST("/login", r.userHandler.Login)       // Login
		}

		// Protected routes (authentication required)
		authMiddleware := middleware.AuthMiddleware(r.userService)
		protected := api.Group("")
		protected.Use(authMiddleware)
		{
			usersProtected := protected.Group("/users")
			{
				usersProtected.POST("", r.userHandler.Create)
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

			mediaProtected := protected.Group("/media")
			{
				mediaProtected.POST("", r.mediaHandler.Create)
				mediaProtected.GET("", r.mediaHandler.List)
				mediaProtected.GET("/:id", r.mediaHandler.Get)
				mediaProtected.PUT("/:id", r.mediaHandler.Update)
				mediaProtected.DELETE("/:id", r.mediaHandler.Delete)
			}
		}
	}

	// Health check endpoint
	engine.GET("/health", func(c *gin.Context) {
		response.SuccessResponseOK(c, "Service is healthy", gin.H{"status": "ok"})
	})
}
