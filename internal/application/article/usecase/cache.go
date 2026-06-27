package usecase

import (
	"context"

	"github.com/rulzi/hexa-go/internal/application/article/dto"
	articleentity "github.com/rulzi/hexa-go/internal/domain/article/entity"
)

// ArticleCache is an application-layer port for caching single article entities.
type ArticleCache interface {
	Get(ctx context.Context, id int64) (*articleentity.Article, error)
	Set(ctx context.Context, id int64, article *articleentity.Article) error
	Delete(ctx context.Context, id int64) error
	InvalidateList(ctx context.Context) error
}

// ArticleListCache defines the interface for article list caching.
type ArticleListCache interface {
	GetArticleList(ctx context.Context, limit, offset int) (*dto.ListArticlesResponse, error)
	SetArticleList(ctx context.Context, limit, offset int, listResp *dto.ListArticlesResponse) error
	InvalidateArticleList(ctx context.Context) error
}
