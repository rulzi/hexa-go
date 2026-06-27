package usecase

import (
	"context"

	userentity "github.com/rulzi/hexa-go/internal/domain/user/entity"
	userport "github.com/rulzi/hexa-go/internal/domain/user/port"
)

type deleteUserDeps struct {
	userRepo userport.Repository
}

// DeleteUser handles user deletion
type DeleteUser struct {
	deps deleteUserDeps
}

// NewDeleteUser creates a new DeleteUser use case.
func NewDeleteUser(userRepo userport.Repository) *DeleteUser {
	return &DeleteUser{deps: deleteUserDeps{userRepo: userRepo}}
}

// Execute deletes a user
func (uc *DeleteUser) Execute(ctx context.Context, id int64) error {
	existingUser, err := uc.deps.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingUser == nil {
		return userentity.NewUserNotFound()
	}

	return uc.deps.userRepo.Delete(ctx, id)
}
