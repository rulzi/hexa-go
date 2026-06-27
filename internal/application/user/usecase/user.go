package usecase

import (
	userport "github.com/rulzi/hexa-go/internal/domain/user/port"
)

type userDeps struct {
	userRepo            userport.Repository
	passwordHasher      userport.PasswordHasher
	notificationService userport.NotificationService
	tokenGen            userport.TokenGenerator
}

// UserUseCase groups all user use case operations
type UserUseCase struct {
	Create *CreateUser
	Get    *GetUser
	List   *ListUser
	Update *UpdateUser
	Delete *DeleteUser
	Login  *LoginUser
}

// NewUserUseCase creates a new UserUseCase
func NewUserUseCase(
	userRepo userport.Repository,
	passwordHasher userport.PasswordHasher,
	notificationService userport.NotificationService,
	tokenGen userport.TokenGenerator,
) *UserUseCase {
	deps := userDeps{
		userRepo:            userRepo,
		passwordHasher:      passwordHasher,
		notificationService: notificationService,
		tokenGen:            tokenGen,
	}

	return &UserUseCase{
		Create: newCreateUser(deps),
		Get:    newGetUser(deps),
		List:   newListUser(deps),
		Update: newUpdateUser(deps),
		Delete: newDeleteUser(deps),
		Login:  newLoginUser(deps),
	}
}
