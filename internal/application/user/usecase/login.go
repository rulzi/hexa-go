package usecase

import (
	"context"

	"github.com/rulzi/hexa-go/internal/application/user/dto"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
)

// LoginUseCase handles user login
type LoginUseCase struct {
	userRepo    domainuser.Repository
	userService *domainuser.Service
}

// NewLoginUseCase creates a new LoginUseCase
func NewLoginUseCase(
	userRepo domainuser.Repository,
	userService *domainuser.Service,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:    userRepo,
		userService: userService,
	}
}

// Execute executes the login use case
func (uc *LoginUseCase) Execute(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Get user by email
	userEntity, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, domainuser.ErrInvalidCredentials
	}

	// Verify password
	if !uc.userService.VerifyPassword(userEntity.Password, req.Password) {
		return nil, domainuser.ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := uc.userService.GenerateJWT(userEntity.ID, userEntity.Email)
	if err != nil {
		return nil, err
	}

	// Return response
	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        userEntity.ID,
			Name:      userEntity.Name,
			Email:     userEntity.Email,
			CreatedAt: userEntity.CreatedAt,
			UpdatedAt: userEntity.UpdatedAt,
		},
	}, nil
}
