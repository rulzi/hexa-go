package user

import (
	"context"

	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
)

// DeleteUserUseCase handles deleting a user
type DeleteUserUseCase struct {
	userRepo domainuser.Repository
}

// NewDeleteUserUseCase creates a new DeleteUserUseCase
func NewDeleteUserUseCase(userRepo domainuser.Repository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		userRepo: userRepo,
	}
}

// Execute executes the delete user use case
func (uc *DeleteUserUseCase) Execute(ctx context.Context, id int64) error {
	// Check if user exists
	existingUser, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingUser == nil {
		return domainuser.ErrUserNotFound
	}

	// Delete user
	return uc.userRepo.Delete(ctx, id)
}
