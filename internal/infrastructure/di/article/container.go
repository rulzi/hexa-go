package article

import (
	"database/sql"
	"time"

	"github.com/redis/go-redis/v9"
	articlecache "github.com/rulzi/hexa-go/internal/adapters/cache/article"
	httparticle "github.com/rulzi/hexa-go/internal/adapters/http/article"
	articledb "github.com/rulzi/hexa-go/internal/adapters/repository/article"
	"github.com/rulzi/hexa-go/internal/application/article/usecase"
	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
)

// Container holds all article domain dependencies
type Container struct {
	Repo        domainarticle.Repository
	Service     *domainarticle.Service
	ArticleUC   *usecase.ArticleUseCase
	Handler     *httparticle.Handler
}

// NewContainer creates a new article domain container
func NewContainer(database *sql.DB, redisClient *redis.Client, appLogger logger.Logger) *Container {
	// Initialize repository (driven adapter)
	articleRepo := articledb.NewMySQLRepository(database)

	// Initialize cache (driven adapter)
	var domainCache domainarticle.Cache
	var dtoCache usecase.ArticleListCache
	if redisClient != nil {
		dtoCacheAdapter := articlecache.NewRedisCache(redisClient, 5*time.Minute)
		domainCache = articlecache.NewDomainCacheAdapter(dtoCacheAdapter)
		dtoCache = dtoCacheAdapter
	}

	// Initialize domain service
	articleService := domainarticle.NewService(articleRepo)

	// Initialize use case (application layer)
	articleUseCase := usecase.NewArticleUseCase(articleRepo, articleService, domainCache, dtoCache)

	// Initialize HTTP handler (driving adapter)
	articleHandler := httparticle.NewHandler(articleUseCase, appLogger)

	return &Container{
		Repo:      articleRepo,
		Service:   articleService,
		ArticleUC: articleUseCase,
		Handler:   articleHandler,
	}
}
