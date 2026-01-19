package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rulzi/hexa-go/internal/application/user/dto"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
	"github.com/stretchr/testify/assert"
)

// Ensure dto is used
var _ *dto.LoginResponse

func TestNewLoginUseCase(t *testing.T) {
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	tokenGen := &mockTokenGenerator{}

	uc := NewLoginUseCase(repo, passwordHasher, tokenGen)

	assert.NotNil(t, uc)
	assert.Equal(t, repo, uc.userRepo)
	assert.Equal(t, passwordHasher, uc.passwordHasher)
	assert.Equal(t, tokenGen, uc.tokenGen)
}

func TestLoginUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	tokenGen := &mockTokenGenerator{}

	uc := NewLoginUseCase(repo, passwordHasher, tokenGen)

	req := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	userEntity := &domainuser.User{
		ID:        1,
		Name:      "Test User",
		Email:     req.Email,
		Password:  "hashed_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	token := "jwt_token_123"

	repo.On("GetByEmail", ctx, req.Email).Return(userEntity, nil)
	passwordHasher.On("Verify", userEntity.Password, req.Password).Return(true)
	tokenGen.On("Generate", userEntity.ID, userEntity.Email).Return(token, nil)

	result, err := uc.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, token, result.Token)
	assert.Equal(t, userEntity.ID, result.User.ID)
	assert.Equal(t, userEntity.Name, result.User.Name)
	assert.Equal(t, userEntity.Email, result.User.Email)

	repo.AssertExpectations(t)
	passwordHasher.AssertExpectations(t)
	tokenGen.AssertExpectations(t)
}

func TestLoginUseCase_Execute_UserNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	tokenGen := &mockTokenGenerator{}

	uc := NewLoginUseCase(repo, passwordHasher, tokenGen)

	req := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	repo.On("GetByEmail", ctx, req.Email).Return(nil, errors.New("not found"))

	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, domainuser.ErrInvalidCredentials, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
	passwordHasher.AssertNotCalled(t, "Verify")
	tokenGen.AssertNotCalled(t, "Generate")
}

func TestLoginUseCase_Execute_GetByEmailError(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	tokenGen := &mockTokenGenerator{}

	uc := NewLoginUseCase(repo, passwordHasher, tokenGen)

	req := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	repoError := errors.New("database error")
	repo.On("GetByEmail", ctx, req.Email).Return(nil, repoError)

	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, domainuser.ErrInvalidCredentials, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
	passwordHasher.AssertNotCalled(t, "Verify")
}

func TestLoginUseCase_Execute_InvalidPassword(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	tokenGen := &mockTokenGenerator{}

	uc := NewLoginUseCase(repo, passwordHasher, tokenGen)

	req := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "wrong_password",
	}

	userEntity := &domainuser.User{
		ID:        1,
		Name:      "Test User",
		Email:     req.Email,
		Password:  "hashed_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("GetByEmail", ctx, req.Email).Return(userEntity, nil)
	passwordHasher.On("Verify", userEntity.Password, req.Password).Return(false)

	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, domainuser.ErrInvalidCredentials, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
	passwordHasher.AssertExpectations(t)
	tokenGen.AssertNotCalled(t, "Generate")
}

func TestLoginUseCase_Execute_TokenGenerationError(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	tokenGen := &mockTokenGenerator{}

	uc := NewLoginUseCase(repo, passwordHasher, tokenGen)

	req := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	userEntity := &domainuser.User{
		ID:        1,
		Name:      "Test User",
		Email:     req.Email,
		Password:  "hashed_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tokenError := errors.New("token generation error")

	repo.On("GetByEmail", ctx, req.Email).Return(userEntity, nil)
	passwordHasher.On("Verify", userEntity.Password, req.Password).Return(true)
	tokenGen.On("Generate", userEntity.ID, userEntity.Email).Return("", tokenError)

	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, tokenError, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
	passwordHasher.AssertExpectations(t)
	tokenGen.AssertExpectations(t)
}

