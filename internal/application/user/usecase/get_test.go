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
var _ *dto.UserResponse

func TestNewGetUserUseCase(t *testing.T) {
	repo := &mockUserRepository{}

	uc := NewGetUserUseCase(repo)

	assert.NotNil(t, uc)
	assert.Equal(t, repo, uc.userRepo)
}

func TestGetUserUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}

	uc := NewGetUserUseCase(repo)

	userID := int64(1)
	userEntity := &domainuser.User{
		ID:        userID,
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "hashed_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, userID).Return(userEntity, nil)

	result, err := uc.Execute(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userEntity.ID, result.ID)
	assert.Equal(t, userEntity.Name, result.Name)
	assert.Equal(t, userEntity.Email, result.Email)
	// Password should not be in response (UserResponse doesn't have Password field)

	repo.AssertExpectations(t)
}

func TestGetUserUseCase_Execute_UserNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}

	uc := NewGetUserUseCase(repo)

	userID := int64(1)

	repo.On("GetByID", ctx, userID).Return(nil, nil)

	result, err := uc.Execute(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, domainuser.ErrUserNotFound, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
}

func TestGetUserUseCase_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}

	uc := NewGetUserUseCase(repo)

	userID := int64(1)
	repoError := errors.New("database error")

	repo.On("GetByID", ctx, userID).Return(nil, repoError)

	result, err := uc.Execute(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
}

