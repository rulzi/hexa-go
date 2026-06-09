package port

import (
	"context"

	"github.com/rulzi/hexa-go/internal/domain/article/entity"
)

// Cache is a port for caching article data.
type Cache interface {
	Get(ctx context.Context, id int64) (*entity.Article, error)
	Set(ctx context.Context, id int64, article *entity.Article) error
	Delete(ctx context.Context, id int64) error
	InvalidateList(ctx context.Context) error
}
