package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
	"github.com/stretchr/testify/assert"
)

func TestNewDeleteUserUseCase(t *testing.T) {
	repo := &mockUserRepository{}

	uc := NewDeleteUserUseCase(repo)

	assert.NotNil(t, uc)
	assert.Equal(t, repo, uc.userRepo)
}

func TestDeleteUserUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}

	uc := NewDeleteUserUseCase(repo)

	userID := int64(1)
	existingUser := &domainuser.User{
		ID:        userID,
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, userID).Return(existingUser, nil)
	repo.On("Delete", ctx, userID).Return(nil)

	err := uc.Execute(ctx, userID)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDeleteUserUseCase_Execute_UserNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}

	uc := NewDeleteUserUseCase(repo)

	userID := int64(1)

	repo.On("GetByID", ctx, userID).Return(nil, nil)

	err := uc.Execute(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, domainuser.ErrUserNotFound, err)

	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Delete")
}

func TestDeleteUserUseCase_Execute_GetByIDError(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}

	uc := NewDeleteUserUseCase(repo)

	userID := int64(1)
	repoError := errors.New("database error")

	repo.On("GetByID", ctx, userID).Return(nil, repoError)

	err := uc.Execute(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, repoError, err)

	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Delete")
}

func TestDeleteUserUseCase_Execute_DeleteError(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}

	uc := NewDeleteUserUseCase(repo)

	userID := int64(1)
	existingUser := &domainuser.User{
		ID:        userID,
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	deleteError := errors.New("delete error")

	repo.On("GetByID", ctx, userID).Return(existingUser, nil)
	repo.On("Delete", ctx, userID).Return(deleteError)

	err := uc.Execute(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, deleteError, err)

	repo.AssertExpectations(t)
}

