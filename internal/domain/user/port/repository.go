package port

//go:generate mockgen -destination=mocks/mock_repository.go -package=mocks github.com/rulzi/hexa-go/internal/domain/user/port Repository

import (
	"context"

	"github.com/rulzi/hexa-go/internal/domain/user/entity"
)

// Repository is the driven port (interface) for user persistence.
// This defines what the domain needs, not how it's implemented.
type Repository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*entity.User, error)
	Count(ctx context.Context) (int64, error)
}
