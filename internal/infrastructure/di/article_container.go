package di

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

// ArticleContainer holds all article domain dependencies
type ArticleContainer struct {
	Repo          domainarticle.Repository
	Service       *domainarticle.Service
	CreateUseCase *usecase.CreateArticleUseCase
	GetUseCase    *usecase.GetArticleUseCase
	ListUseCase   *usecase.ListArticlesUseCase
	UpdateUseCase *usecase.UpdateArticleUseCase
	DeleteUseCase *usecase.DeleteArticleUseCase
	Handler       *httparticle.Handler
}

// NewArticleContainer creates a new article domain container
func NewArticleContainer(database *sql.DB, redisClient *redis.Client) *ArticleContainer {
	// Initialize repository (driven adapter)
	articleRepo := articledb.NewMySQLRepository(database)

	// Initialize cache (driven adapter)
	var articleCache usecase.ArticleCache
	var articleSingleCache usecase.ArticleSingleCache
	if redisClient != nil {
		cacheAdapter := articlecache.NewRedisCache(redisClient, 5*time.Minute)
		articleCache = cacheAdapter
		articleSingleCache = cacheAdapter
	}

	// Initialize domain service
	articleService := domainarticle.NewService(articleRepo)

	// Initialize use cases (application layer)
	createArticleUseCase := usecase.NewCreateArticleUseCase(articleRepo, articleService, articleCache)
	getArticleUseCase := usecase.NewGetArticleUseCase(articleRepo, articleSingleCache)
	listArticlesUseCase := usecase.NewListArticlesUseCase(articleRepo, articleCache)
	updateArticleUseCase := usecase.NewUpdateArticleUseCase(articleRepo, articleService, articleSingleCache, articleCache)
	deleteArticleUseCase := usecase.NewDeleteArticleUseCase(articleRepo, articleSingleCache, articleCache)

	// Initialize HTTP handler (driving adapter)
	articleHandler := httparticle.NewHandler(
		createArticleUseCase,
		getArticleUseCase,
		listArticlesUseCase,
		updateArticleUseCase,
		deleteArticleUseCase,
	)

	return &ArticleContainer{
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
