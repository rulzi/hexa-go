package di

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
	"github.com/rulzi/hexa-go/internal/adapters/http"
	"github.com/rulzi/hexa-go/internal/infrastructure/config"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
	diarticle "github.com/rulzi/hexa-go/internal/infrastructure/di/article"
	dimedia "github.com/rulzi/hexa-go/internal/infrastructure/di/media"
	diuser "github.com/rulzi/hexa-go/internal/infrastructure/di/user"
)

// Container holds all dependencies
type Container struct {
	DB      *sql.DB
	Redis   *redis.Client
	Logger  logger.Logger
	User    *diuser.Container
	Article *diarticle.Container
	Media   *dimedia.Container
	Router  *http.Router
}

// NewContainer creates a new dependency injection container
func NewContainer(database *sql.DB, redisClient *redis.Client, appLogger logger.Logger, jwtSecret string, jwtExpiration int, storageCfg config.StorageConfig) (*Container, error) {
	// Initialize domain containers
	userContainer := diuser.NewContainer(database, appLogger, jwtSecret, jwtExpiration)
	articleContainer := diarticle.NewContainer(database, redisClient, appLogger)
	mediaContainer, err := dimedia.NewContainer(database, appLogger, storageCfg)
	if err != nil {
		return nil, err
	}

	// Initialize router
	router := http.NewRouter(userContainer.Handler, articleContainer.Handler, mediaContainer.Handler, userContainer.TokenValidator, storageCfg.BasePath, appLogger)

	return &Container{
		DB:      database,
		Redis:   redisClient,
		Logger:  appLogger,
		User:    userContainer,
		Article: articleContainer,
		Media:   mediaContainer,
		Router:  router,
	}, nil
}
