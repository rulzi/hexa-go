package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rulzi/hexa-go/internal/infrastructure/config"
	"github.com/rulzi/hexa-go/internal/infrastructure/database"
	"github.com/rulzi/hexa-go/internal/infrastructure/di"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	appLogger := logger.NewSimpleLogger()
	appLogger.Info("Starting application...")

	// Connect to database
	db, err := database.NewMySQLConnection(cfg.Database.GetDSN())
	if err != nil {
		appLogger.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	defer db.Close()
	appLogger.Info("Database connected successfully")

	// Connect to Redis
	var redisClient *redis.Client
	if cfg.Redis.Host != "" {
		redisClient, err = database.NewRedisConnection(
			cfg.Redis.GetAddr(),
			cfg.Redis.Password,
			cfg.Redis.DB,
		)
		if err != nil {
			appLogger.Error(fmt.Sprintf("Failed to connect to Redis: %v. Continuing without cache.", err))
		} else {
			defer redisClient.Close()
			appLogger.Info("Redis connected successfully")
		}
	}

	// Initialize dependency injection container
	container := di.NewContainer(db, redisClient, cfg.JWT.Secret, cfg.JWT.Expiration)

	// Initialize database tables
	ctx := context.Background()
	if err := container.InitializeDatabase(ctx); err != nil {
		appLogger.Fatal(fmt.Sprintf("Failed to initialize database: %v", err))
	}
	appLogger.Info("Database tables initialized")

	// Setup Gin router
	if cfg.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// Setup routes
	container.Router.SetupRoutes(router, cfg.Server.Debug)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	appLogger.Info(fmt.Sprintf("Server starting on %s", addr))

	if err := http.ListenAndServe(addr, router); err != nil {
		appLogger.Fatal(fmt.Sprintf("Failed to start server: %v", err))
	}
}
