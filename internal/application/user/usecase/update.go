package usecase

import (
	"context"
	"time"

	"github.com/rulzi/hexa-go/internal/application/user/dto"
	userentity "github.com/rulzi/hexa-go/internal/domain/user/entity"
)

// UpdateUser handles user updates
type UpdateUser struct {
	deps userDeps
}

func newUpdateUser(deps userDeps) *UpdateUser {
	return &UpdateUser{deps: deps}
}

// Execute updates a user
func (uc *UpdateUser) Execute(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	existingUser, err := uc.deps.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, userentity.NewUserNotFound()
	}

	if req.Email != existingUser.Email {
		emailUser, err := uc.deps.userRepo.GetByEmail(ctx, req.Email)
		if err == nil && emailUser != nil {
			return nil, userentity.NewEmailExists()
		}
	}

	existingUser.Name = req.Name
	existingUser.Email = req.Email
	existingUser.UpdatedAt = time.Now()

	if req.Password != "" {
		if err := userentity.ValidatePassword(req.Password); err != nil {
			return nil, err
		}

		hashedPassword, err := uc.deps.passwordHasher.Hash(req.Password)
		if err != nil {
			return nil, err
		}
		existingUser.Password = hashedPassword
	}

	if err := existingUser.Validate(); err != nil {
		return nil, err
	}

	updatedUser, err := uc.deps.userRepo.Update(ctx, existingUser)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        updatedUser.ID,
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	}, nil
}
