package di

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
	"github.com/rulzi/hexa-go/internal/adapters/http"
	diarticle "github.com/rulzi/hexa-go/internal/infrastructure/di/article"
	dimedia "github.com/rulzi/hexa-go/internal/infrastructure/di/media"
	diuser "github.com/rulzi/hexa-go/internal/infrastructure/di/user"
)

// Container holds all dependencies
type Container struct {
	DB      *sql.DB
	Redis   *redis.Client
	User    *diuser.Container
	Article *diarticle.Container
	Media   *dimedia.Container
	Router  *http.Router
}

// NewContainer creates a new dependency injection container
func NewContainer(database *sql.DB, redisClient *redis.Client, jwtSecret string, jwtExpiration int, storageBasePath string, storageBaseURL string) (*Container, error) {
	// Initialize domain containers
	userContainer := diuser.NewContainer(database, jwtSecret, jwtExpiration)
	articleContainer := diarticle.NewContainer(database, redisClient)
	mediaContainer, err := dimedia.NewContainer(database, storageBasePath, storageBaseURL)
	if err != nil {
		return nil, err
	}

	// Initialize router
	router := http.NewRouter(userContainer.Handler, articleContainer.Handler, mediaContainer.Handler, userContainer.TokenValidator, storageBasePath)

	return &Container{
		DB:      database,
		Redis:   redisClient,
		User:    userContainer,
		Article: articleContainer,
		Media:   mediaContainer,
		Router:  router,
	}, nil
}
