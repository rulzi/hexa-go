package user

import (
	"context"

	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
)

// GetUserUseCase handles retrieving a user by ID
type GetUserUseCase struct {
	userRepo domainuser.Repository
}

// NewGetUserUseCase creates a new GetUserUseCase
func NewGetUserUseCase(userRepo domainuser.Repository) *GetUserUseCase {
	return &GetUserUseCase{
		userRepo: userRepo,
	}
}

// Execute executes the get user use case
func (uc *GetUserUseCase) Execute(ctx context.Context, id int64) (*UserResponse, error) {
	userEntity, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if userEntity == nil {
		return nil, domainuser.ErrUserNotFound
	}

	return &UserResponse{
		ID:        userEntity.ID,
		Name:      userEntity.Name,
		Email:     userEntity.Email,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
	}, nil
}
