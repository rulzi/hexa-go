package di

import (
	"database/sql"
	"time"

	"github.com/redis/go-redis/v9"
	articlecache "github.com/rulzi/hexa-go/internal/adapters/cache/article"
	articledb "github.com/rulzi/hexa-go/internal/adapters/db/article"
	httparticle "github.com/rulzi/hexa-go/internal/adapters/http/article"
	apparticle "github.com/rulzi/hexa-go/internal/application/article"
	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
)

// ArticleContainer holds all article domain dependencies
type ArticleContainer struct {
	Repo          domainarticle.Repository
	Service       *domainarticle.Service
	CreateUseCase *apparticle.CreateArticleUseCase
	GetUseCase    *apparticle.GetArticleUseCase
	ListUseCase   *apparticle.ListArticlesUseCase
	UpdateUseCase *apparticle.UpdateArticleUseCase
	DeleteUseCase *apparticle.DeleteArticleUseCase
	Handler       *httparticle.Handler
}

// NewArticleContainer creates a new article domain container
func NewArticleContainer(database *sql.DB, redisClient *redis.Client) *ArticleContainer {
	// Initialize repository (driven adapter)
	articleRepo := articledb.NewMySQLRepository(database)

	// Initialize cache (driven adapter)
	var articleCache apparticle.ArticleCache
	var articleSingleCache apparticle.ArticleSingleCache
	if redisClient != nil {
		cacheAdapter := articlecache.NewRedisCache(redisClient, 5*time.Minute)
		articleCache = cacheAdapter
		articleSingleCache = cacheAdapter
	}

	// Initialize domain service
	articleService := domainarticle.NewService(articleRepo)

	// Initialize use cases (application layer)
	createArticleUseCase := apparticle.NewCreateArticleUseCase(articleRepo, articleService, articleCache)
	getArticleUseCase := apparticle.NewGetArticleUseCase(articleRepo, articleSingleCache)
	listArticlesUseCase := apparticle.NewListArticlesUseCase(articleRepo, articleCache)
	updateArticleUseCase := apparticle.NewUpdateArticleUseCase(articleRepo, articleService, articleSingleCache, articleCache)
	deleteArticleUseCase := apparticle.NewDeleteArticleUseCase(articleRepo, articleSingleCache, articleCache)

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
