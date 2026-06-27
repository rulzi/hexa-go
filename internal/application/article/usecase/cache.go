package usecase

import (
	"context"

	articleentity "github.com/rulzi/hexa-go/internal/domain/article/entity"
)

// ArticleListPage holds a paginated slice of article entities for caching.
type ArticleListPage struct {
	Articles []*articleentity.Article
	Total    int64
}

// ArticleCache is an application-layer port for caching single article entities.
type ArticleCache interface {
	Get(ctx context.Context, id int64) (*articleentity.Article, error)
	Set(ctx context.Context, id int64, article *articleentity.Article) error
	Delete(ctx context.Context, id int64) error
	InvalidateList(ctx context.Context) error
}

// ArticleListCache defines the interface for article list caching.
type ArticleListCache interface {
	GetList(ctx context.Context, limit, offset int) (*ArticleListPage, error)
	SetList(ctx context.Context, limit, offset int, page *ArticleListPage) error
	InvalidateList(ctx context.Context) error
}
