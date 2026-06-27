package usecase

import (
	"context"
	"time"

	"github.com/rulzi/hexa-go/internal/application/user/dto"
	userentity "github.com/rulzi/hexa-go/internal/domain/user/entity"
)

// CreateUser handles user creation
type CreateUser struct {
	deps userDeps
}

func newCreateUser(deps userDeps) *CreateUser {
	return &CreateUser{deps: deps}
}

func validateCreateUserRequest(req dto.CreateUserRequest) error {
	if req.Name == "" {
		return userentity.NewNameRequired()
	}
	if req.Email == "" {
		return userentity.NewEmailRequired()
	}
	if req.Password == "" {
		return userentity.NewPasswordRequired()
	}
	if len(req.Password) < 6 {
		return userentity.NewPasswordTooShort()
	}
	return nil
}

// Execute creates a new user
func (uc *CreateUser) Execute(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	if err := validateCreateUserRequest(req); err != nil {
		return nil, err
	}

	existingUser, err := uc.deps.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, userentity.NewEmailExists()
	}

	hashedPassword, err := uc.deps.passwordHasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	newUser := &userentity.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := newUser.Validate(); err != nil {
		return nil, err
	}

	createdUser, err := uc.deps.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	_ = uc.deps.notificationService.SendWelcomeEmail(ctx, createdUser.Email, createdUser.Name)

	return &dto.UserResponse{
		ID:        createdUser.ID,
		Name:      createdUser.Name,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	}, nil
}
