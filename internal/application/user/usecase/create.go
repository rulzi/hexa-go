package usecase

import (
	"context"
	"time"

	"github.com/rulzi/hexa-go/internal/application/user/dto"
	userentity "github.com/rulzi/hexa-go/internal/domain/user/entity"
	userport "github.com/rulzi/hexa-go/internal/domain/user/port"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
)

type createUserDeps struct {
	userRepo            userport.Repository
	passwordHasher      userport.PasswordHasher
	notificationService userport.NotificationService
	logger              logger.Logger
}

// CreateUser handles user creation
type CreateUser struct {
	deps createUserDeps
}

// NewCreateUser creates a new CreateUser use case.
func NewCreateUser(
	userRepo userport.Repository,
	passwordHasher userport.PasswordHasher,
	notificationService userport.NotificationService,
	appLogger logger.Logger,
) *CreateUser {
	return &CreateUser{deps: createUserDeps{
		userRepo:            userRepo,
		passwordHasher:      passwordHasher,
		notificationService: notificationService,
		logger:              appLogger,
	}}
}

// Execute creates a new user
func (uc *CreateUser) Execute(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	if err := userentity.ValidateRegistration(req.Name, req.Email, req.Password); err != nil {
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

	if err := uc.deps.notificationService.SendWelcomeEmail(ctx, createdUser.Email, createdUser.Name); err != nil {
		uc.deps.logger.WarnWithFields("failed to send welcome email", map[string]interface{}{
			"email": createdUser.Email,
			"error": err.Error(),
		})
	}

	return &dto.UserResponse{
		ID:        createdUser.ID,
		Name:      createdUser.Name,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	}, nil
}
