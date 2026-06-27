package usecase

import (
	"context"

	"github.com/rulzi/hexa-go/internal/application/user/dto"
	userentity "github.com/rulzi/hexa-go/internal/domain/user/entity"
	userport "github.com/rulzi/hexa-go/internal/domain/user/port"
)

type getUserDeps struct {
	userRepo userport.Repository
}

// GetUser handles user retrieval by ID
type GetUser struct {
	deps getUserDeps
}

// NewGetUser creates a new GetUser use case.
func NewGetUser(userRepo userport.Repository) *GetUser {
	return &GetUser{deps: getUserDeps{userRepo: userRepo}}
}

// Execute retrieves a user by ID
func (uc *GetUser) Execute(ctx context.Context, id int64) (*dto.UserResponse, error) {
	userEntity, err := uc.deps.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if userEntity == nil {
		return nil, userentity.NewUserNotFound()
	}

	return &dto.UserResponse{
		ID:        userEntity.ID,
		Name:      userEntity.Name,
		Email:     userEntity.Email,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
	}, nil
}
