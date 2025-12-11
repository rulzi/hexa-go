package di

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
	"github.com/rulzi/hexa-go/internal/adapters/http"
	diarticle "github.com/rulzi/hexa-go/internal/infrastructure/di/article"
	diuser "github.com/rulzi/hexa-go/internal/infrastructure/di/user"
)

// Container holds all dependencies
type Container struct {
	DB      *sql.DB
	Redis   *redis.Client
	User    *diuser.Container
	Article *diarticle.Container
	Router  *http.Router
}

// NewContainer creates a new dependency injection container
func NewContainer(database *sql.DB, redisClient *redis.Client, jwtSecret string, jwtExpiration int) *Container {
	// Initialize domain containers
	userContainer := diuser.NewContainer(database, jwtSecret, jwtExpiration)
	articleContainer := diarticle.NewContainer(database, redisClient)

	// Initialize router
	router := http.NewRouter(userContainer.Handler, articleContainer.Handler, userContainer.Service)

	return &Container{
		DB:      database,
		Redis:   redisClient,
		User:    userContainer,
		Article: articleContainer,
		Router:  router,
	}
}
