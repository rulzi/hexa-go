package port

import (
	"context"

	"github.com/rulzi/hexa-go/internal/domain/article/entity"
)

// Repository is the driven port (interface) for article persistence.
type Repository interface {
	Create(ctx context.Context, article *entity.Article) (*entity.Article, error)
	GetByID(ctx context.Context, id int64) (*entity.Article, error)
	Update(ctx context.Context, article *entity.Article) (*entity.Article, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*entity.Article, error)
	ListByAuthor(ctx context.Context, authorID int64, limit, offset int) ([]*entity.Article, error)
	Count(ctx context.Context) (int64, error)
	CountByAuthor(ctx context.Context, authorID int64) (int64, error)
}
