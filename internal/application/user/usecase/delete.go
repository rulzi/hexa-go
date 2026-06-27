package usecase

import (
	"context"

	userentity "github.com/rulzi/hexa-go/internal/domain/user/entity"
)

// DeleteUser handles user deletion
type DeleteUser struct {
	deps userDeps
}

func newDeleteUser(deps userDeps) *DeleteUser {
	return &DeleteUser{deps: deps}
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
