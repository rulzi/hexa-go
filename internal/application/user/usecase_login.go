package user

import (
	"context"

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

// LoginRequest represents the request DTO for login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response DTO for login
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// Execute executes the login use case
func (uc *LoginUseCase) Execute(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
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
	return &LoginResponse{
		Token: token,
		User: UserResponse{
			ID:        userEntity.ID,
			Name:      userEntity.Name,
			Email:     userEntity.Email,
			CreatedAt: userEntity.CreatedAt,
			UpdatedAt: userEntity.UpdatedAt,
		},
	}, nil
}

