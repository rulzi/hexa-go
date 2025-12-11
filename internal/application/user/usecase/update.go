package usecase

import (
	"context"
	"time"

	"github.com/rulzi/hexa-go/internal/application/user/dto"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
)

// UpdateUserUseCase handles updating a user
type UpdateUserUseCase struct {
	userRepo    domainuser.Repository
	userService *domainuser.Service
}

// NewUpdateUserUseCase creates a new UpdateUserUseCase
func NewUpdateUserUseCase(
	userRepo domainuser.Repository,
	userService *domainuser.Service,
) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userRepo:    userRepo,
		userService: userService,
	}
}

// Execute executes the update user use case
func (uc *UpdateUserUseCase) Execute(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// Get existing user
	existingUser, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, domainuser.ErrUserNotFound
	}

	// Check if email is being changed and if it already exists
	if req.Email != existingUser.Email {
		emailUser, err := uc.userRepo.GetByEmail(ctx, req.Email)
		if err == nil && emailUser != nil {
			return nil, domainuser.ErrEmailExists
		}
	}

	// Update fields
	existingUser.Name = req.Name
	existingUser.Email = req.Email
	existingUser.UpdatedAt = time.Now()

	// Update password if provided
	if req.Password != "" {
		hashedPassword, err := uc.userService.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		existingUser.Password = hashedPassword
	}

	// Validate entity
	if err := existingUser.Validate(); err != nil {
		return nil, err
	}

	// Update in repository
	updatedUser, err := uc.userRepo.Update(ctx, existingUser)
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
