package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rulzi/hexa-go/internal/infrastructure/config"
	"github.com/rulzi/hexa-go/internal/infrastructure/database"
	"github.com/rulzi/hexa-go/internal/infrastructure/di"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		return
	}

	// Initialize logger with file support
	appLogger, err := logger.NewLogger(logger.LoggerConfig{
		Level:         cfg.Logger.Level,
		Format:        cfg.Logger.Format,
		FilePath:      cfg.Logger.FilePath,
		EnableFile:    cfg.Logger.EnableFile,
		EnableConsole: cfg.Logger.EnableConsole,
		ReportCaller:  cfg.Logger.ReportCaller,
	})
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}
	defer appLogger.Close()

	appLogger.Info("Starting application...")
	appLogger.InfoWithFields("Application configuration loaded", map[string]interface{}{
		"log_level":      cfg.Logger.Level,
		"log_format":     cfg.Logger.Format,
		"log_file":       cfg.Logger.FilePath,
		"enable_file":    cfg.Logger.EnableFile,
		"enable_console": cfg.Logger.EnableConsole,
	})

	// Connect to database
	db, err := database.NewMySQLConnection(cfg.Database.GetDSN())
	if err != nil {
		appLogger.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	defer func() {
		if err := db.Close(); err != nil {
			appLogger.Error(fmt.Sprintf("Failed to close database connection: %v", err))
		}
	}()
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
			defer func() {
				if err := redisClient.Close(); err != nil {
					appLogger.Error(fmt.Sprintf("Failed to close Redis connection: %v", err))
				}
			}()
		}
		appLogger.Info("Redis connected successfully")
	}

	// Initialize dependency injection container
	container, err := di.NewContainer(db, redisClient, appLogger, cfg.JWT.Secret, cfg.JWT.Expiration, cfg.Storage)
	if err != nil {
		appLogger.Fatal(fmt.Sprintf("Failed to initialize container: %v", err))
	}

	// Setup Gin router
	if cfg.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// Setup routes
	container.Router.SetupRoutes(router)

	// Start server with graceful shutdown
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		appLogger.Info(fmt.Sprintf("Server starting on %s", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal(fmt.Sprintf("Failed to start server: %v", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	const shutdownTimeout = 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error(fmt.Sprintf("Server forced to shutdown: %v", err))
	}

	appLogger.Info("Server exited")
}
