package di

import (
	"context"
	"database/sql"

	"github.com/redis/go-redis/v9"
	"github.com/rulzi/hexa-go/internal/adapters/db"
	"github.com/rulzi/hexa-go/internal/adapters/http"
)

// Container holds all dependencies
type Container struct {
	DB      *sql.DB
	Redis   *redis.Client
	User    *UserContainer
	Article *ArticleContainer
	Router  *http.Router
}

// NewContainer creates a new dependency injection container
func NewContainer(database *sql.DB, redisClient *redis.Client, jwtSecret string, jwtExpiration int) *Container {
	// Initialize domain containers
	userContainer := NewUserContainer(database, jwtSecret, jwtExpiration)
	articleContainer := NewArticleContainer(database, redisClient)

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

// InitializeDatabase creates the necessary database tables
func (c *Container) InitializeDatabase(ctx context.Context) error {
	// Create users table
	userRepo := c.User.Repo.(*db.UserMySQLRepository)
	if err := userRepo.CreateTable(ctx); err != nil {
		return err
	}

	// Create articles table
	articleRepo := c.Article.Repo.(*db.ArticleMySQLRepository)
	if err := articleRepo.CreateTable(ctx); err != nil {
		return err
	}

	return nil
}
