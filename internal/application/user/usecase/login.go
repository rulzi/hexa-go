package usecase

import (
	"context"

	"github.com/rulzi/hexa-go/internal/application/user/dto"
	userentity "github.com/rulzi/hexa-go/internal/domain/user/entity"
)

// LoginUser handles user authentication
type LoginUser struct {
	deps userDeps
}

func newLoginUser(deps userDeps) *LoginUser {
	return &LoginUser{deps: deps}
}

// Execute authenticates user and returns token
func (uc *LoginUser) Execute(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	userEntity, err := uc.deps.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, userentity.NewInvalidCredentials()
	}

	if userEntity == nil {
		return nil, userentity.NewInvalidCredentials()
	}

	if !uc.deps.passwordHasher.Verify(userEntity.Password, req.Password) {
		return nil, userentity.NewInvalidCredentials()
	}

	token, err := uc.deps.tokenGen.Generate(userEntity.ID, userEntity.Email)
	if err != nil {
		return nil, err
	}

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
