package usecase

import (
	"context"
	"time"

	"github.com/rulzi/hexa-go/internal/application/user/dto"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
)

// EmailSender defines the interface for sending emails
type EmailSender interface {
	SendWelcomeEmail(ctx context.Context, email, name string) error
}

// CreateUserUseCase handles the creation of a new user
type CreateUserUseCase struct {
	userRepo    domainuser.Repository
	userService *domainuser.Service
	emailSender EmailSender // External service adapter
}

// NewCreateUserUseCase creates a new CreateUserUseCase
func NewCreateUserUseCase(
	userRepo domainuser.Repository,
	userService *domainuser.Service,
	emailSender EmailSender,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo:    userRepo,
		userService: userService,
		emailSender: emailSender,
	}
}

// Execute executes the create user use case
func (uc *CreateUserUseCase) Execute(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check if email already exists
	existingUser, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, domainuser.ErrEmailExists
	}

	// Hash password
	hashedPassword, err := uc.userService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user entity
	newUser := &domainuser.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Validate entity
	if err := newUser.Validate(); err != nil {
		return nil, err
	}

	// Save to repository
	createdUser, err := uc.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	// Send welcome email (external service)
	_ = uc.emailSender.SendWelcomeEmail(ctx, createdUser.Email, createdUser.Name)

	// Return response DTO
	return &dto.UserResponse{
		ID:        createdUser.ID,
		Name:      createdUser.Name,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	}, nil
}
