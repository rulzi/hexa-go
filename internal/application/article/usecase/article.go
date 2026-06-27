package usecase

import (
	"context"

	"github.com/rulzi/hexa-go/internal/application/article/dto"
	articleport "github.com/rulzi/hexa-go/internal/domain/article/port"
	articleservice "github.com/rulzi/hexa-go/internal/domain/article/service"
)

// ArticleListCache defines the interface for article list caching (DTO-based for performance)
type ArticleListCache interface {
	GetArticleList(ctx context.Context, limit, offset int) (*dto.ListArticlesResponse, error)
	SetArticleList(ctx context.Context, limit, offset int, listResp *dto.ListArticlesResponse) error
	InvalidateArticleList(ctx context.Context) error
}

type articleDeps struct {
	articleRepo    articleport.Repository
	articleService *articleservice.Service
	cache          articleport.Cache
	listCache      ArticleListCache
}

// ArticleUseCase groups all article use case operations
type ArticleUseCase struct {
	Create *CreateArticle
	Get    *GetArticle
	List   *ListArticle
	Update *UpdateArticle
	Delete *DeleteArticle
}

// NewArticleUseCase creates a new ArticleUseCase
func NewArticleUseCase(
	articleRepo articleport.Repository,
	articleService *articleservice.Service,
	cache articleport.Cache,
	listCache ArticleListCache,
) *ArticleUseCase {
	deps := articleDeps{
		articleRepo:    articleRepo,
		articleService: articleService,
		cache:          cache,
		listCache:      listCache,
	}

	return &ArticleUseCase{
		Create: newCreateArticle(deps),
		Get:    newGetArticle(deps),
		List:   newListArticle(deps),
		Update: newUpdateArticle(deps),
		Delete: newDeleteArticle(deps),
	}
}
