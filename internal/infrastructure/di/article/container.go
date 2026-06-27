package article

import (
	"database/sql"
	"time"

	"github.com/redis/go-redis/v9"
	articlecache "github.com/rulzi/hexa-go/internal/adapters/cache/article"
	httparticle "github.com/rulzi/hexa-go/internal/adapters/http/article"
	articledb "github.com/rulzi/hexa-go/internal/adapters/repository/article"
	"github.com/rulzi/hexa-go/internal/application/article/usecase"
	articleport "github.com/rulzi/hexa-go/internal/domain/article/port"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
)

// Container holds all article domain dependencies
type Container struct {
	Repo    articleport.Repository
	Handler *httparticle.Handler
}

// NewContainer creates a new article domain container
func NewContainer(database *sql.DB, redisClient *redis.Client, appLogger logger.Logger) *Container {
	articleRepo := articledb.NewMySQLRepository(database, appLogger)

	var articleCache usecase.ArticleCache = usecase.NoopCache{}
	var listCache usecase.ArticleListCache = usecase.NoopCache{}
	if redisClient != nil {
		redisCache := articlecache.NewRedisCache(redisClient, 5*time.Minute)
		articleCache = redisCache
		listCache = redisCache
	}

	articleHandler := httparticle.NewHandler(httparticle.Deps{
		Create: usecase.NewCreateArticle(articleRepo, articleCache),
		Get:    usecase.NewGetArticle(articleRepo, articleCache),
		List:   usecase.NewListArticle(articleRepo, listCache),
		Update: usecase.NewUpdateArticle(articleRepo, articleCache, listCache),
		Delete: usecase.NewDeleteArticle(articleRepo, articleCache, listCache),
	}, appLogger)

	return &Container{
		Repo:    articleRepo,
		Handler: articleHandler,
	}
}
