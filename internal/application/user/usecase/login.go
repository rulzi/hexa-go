package usecase

import (
	"context"

	"github.com/rulzi/hexa-go/internal/application/user/dto"
	userauth "github.com/rulzi/hexa-go/internal/domain/user/auth"
	userport "github.com/rulzi/hexa-go/internal/domain/user/port"
)

type loginUserDeps struct {
	userRepo       userport.Repository
	passwordHasher userport.PasswordHasher
	tokenGen       userport.TokenGenerator
}

// LoginUser handles user authentication
type LoginUser struct {
	deps loginUserDeps
}

// NewLoginUser creates a new LoginUser use case.
func NewLoginUser(
	userRepo userport.Repository,
	passwordHasher userport.PasswordHasher,
	tokenGen userport.TokenGenerator,
) *LoginUser {
	return &LoginUser{deps: loginUserDeps{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenGen:       tokenGen,
	}}
}

// Execute authenticates user and returns token
func (uc *LoginUser) Execute(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	result, err := userauth.Authenticate(
		ctx,
		req.Email,
		req.Password,
		uc.deps.userRepo,
		uc.deps.passwordHasher,
		uc.deps.tokenGen,
	)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: result.Token,
		User: dto.UserResponse{
			ID:        result.User.ID,
			Name:      result.User.Name,
			Email:     result.User.Email,
			CreatedAt: result.User.CreatedAt,
			UpdatedAt: result.User.UpdatedAt,
		},
	}, nil
}
