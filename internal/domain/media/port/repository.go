package port

import (
	"context"

	"github.com/rulzi/hexa-go/internal/domain/media/entity"
)

// Repository is the driven port (interface) for media persistence.
type Repository interface {
	Create(ctx context.Context, media *entity.Media) (*entity.Media, error)
	GetByID(ctx context.Context, id int64) (*entity.Media, error)
	Update(ctx context.Context, media *entity.Media) (*entity.Media, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*entity.Media, error)
	Count(ctx context.Context) (int64, error)
}
