package article

import (
	"database/sql"
	"time"

	"github.com/redis/go-redis/v9"
	articlecache "github.com/rulzi/hexa-go/internal/adapters/cache/article"
	articledb "github.com/rulzi/hexa-go/internal/adapters/db/article"
	httparticle "github.com/rulzi/hexa-go/internal/adapters/http/article"
	"github.com/rulzi/hexa-go/internal/application/article/usecase"
	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
)

// Container holds all article domain dependencies
type Container struct {
	Repo          domainarticle.Repository
	Service       *domainarticle.Service
	CreateUseCase *usecase.CreateArticleUseCase
	GetUseCase    *usecase.GetArticleUseCase
	ListUseCase   *usecase.ListArticlesUseCase
	UpdateUseCase *usecase.UpdateArticleUseCase
	DeleteUseCase *usecase.DeleteArticleUseCase
	Handler       *httparticle.Handler
}

// NewContainer creates a new article domain container
func NewContainer(database *sql.DB, redisClient *redis.Client) *Container {
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

	// Initialize use cases (application layer)
	createArticleUseCase := usecase.NewCreateArticleUseCase(articleRepo, articleService, domainCache)
	getArticleUseCase := usecase.NewGetArticleUseCase(articleRepo, domainCache)
	listArticlesUseCase := usecase.NewListArticlesUseCase(articleRepo, domainCache, dtoCache)
	updateArticleUseCase := usecase.NewUpdateArticleUseCase(articleRepo, articleService, domainCache, dtoCache)
	deleteArticleUseCase := usecase.NewDeleteArticleUseCase(articleRepo, domainCache, dtoCache)

	// Initialize HTTP handler (driving adapter)
	articleHandler := httparticle.NewHandler(
		createArticleUseCase,
		getArticleUseCase,
		listArticlesUseCase,
		updateArticleUseCase,
		deleteArticleUseCase,
	)

	return &Container{
		Repo:          articleRepo,
		Service:       articleService,
		CreateUseCase: createArticleUseCase,
		GetUseCase:    getArticleUseCase,
		ListUseCase:   listArticlesUseCase,
		UpdateUseCase: updateArticleUseCase,
		DeleteUseCase: deleteArticleUseCase,
		Handler:       articleHandler,
	}
}
