package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rulzi/hexa-go/internal/application/user/dto"
	userentity "github.com/rulzi/hexa-go/internal/domain/user/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewUserUseCase(t *testing.T) {
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	notificationService := &mockNotificationService{}
	tokenGen := &mockTokenGenerator{}

	uc := NewUserUseCase(repo, passwordHasher, notificationService, tokenGen)

	assert.NotNil(t, uc)
	assert.Equal(t, repo, uc.userRepo)
	assert.Equal(t, passwordHasher, uc.passwordHasher)
	assert.Equal(t, notificationService, uc.notificationService)
	assert.Equal(t, tokenGen, uc.tokenGen)
}

func TestUserUseCase_Get_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	notificationService := &mockNotificationService{}
	tokenGen := &mockTokenGenerator{}

	uc := NewUserUseCase(repo, passwordHasher, notificationService, tokenGen)

	userID := int64(1)
	userEntity := &userentity.User{
		ID: userID, Name: "Test", Email: "test@example.com",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	repo.On("GetByID", ctx, userID).Return(userEntity, nil)

	result, err := uc.Get(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userEntity.ID, result.ID)
	assert.Equal(t, userEntity.Email, result.Email)
	repo.AssertExpectations(t)
}

func TestUserUseCase_Get_UserNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	notificationService := &mockNotificationService{}
	tokenGen := &mockTokenGenerator{}

	uc := NewUserUseCase(repo, passwordHasher, notificationService, tokenGen)

	userID := int64(999)
	repo.On("GetByID", ctx, userID).Return(nil, nil)

	result, err := uc.Get(ctx, userID)

	assert.Error(t, err)
	assert.True(t, userentity.IsUserNotFound(err))
	assert.Nil(t, result)
}

func TestUserUseCase_List_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	notificationService := &mockNotificationService{}
	tokenGen := &mockTokenGenerator{}

	uc := NewUserUseCase(repo, passwordHasher, notificationService, tokenGen)

	limit, offset := 10, 0
	users := []*userentity.User{
		{ID: 1, Name: "User1", Email: "u1@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	total := int64(1)
	repo.On("List", ctx, limit, offset).Return(users, nil)
	repo.On("Count", ctx).Return(total, nil)

	result, err := uc.List(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, total, result.Total)
	assert.Equal(t, limit, result.Limit)
	assert.Len(t, result.Users, 1)
}

func TestUserUseCase_Update_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	notificationService := &mockNotificationService{}
	tokenGen := &mockTokenGenerator{}

	uc := NewUserUseCase(repo, passwordHasher, notificationService, tokenGen)

	userID := int64(1)
	existingUser := &userentity.User{
		ID: userID, Name: "Old", Email: "old@example.com", Password: "hash",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	req := dto.UpdateUserRequest{Name: "New", Email: "new@example.com"}
	updatedUser := &userentity.User{
		ID: userID, Name: req.Name, Email: req.Email, Password: existingUser.Password,
		CreatedAt: existingUser.CreatedAt, UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, userID).Return(existingUser, nil)
	repo.On("GetByEmail", ctx, req.Email).Return(nil, errors.New("not found"))
	repo.On("Update", ctx, mock.Anything).Return(updatedUser, nil)

	result, err := uc.Update(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Email, result.Email)
	repo.AssertExpectations(t)
}

func TestUserUseCase_Update_UserNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	notificationService := &mockNotificationService{}
	tokenGen := &mockTokenGenerator{}

	uc := NewUserUseCase(repo, passwordHasher, notificationService, tokenGen)

	userID := int64(999)
	req := dto.UpdateUserRequest{Name: "New", Email: "new@example.com"}
	repo.On("GetByID", ctx, userID).Return(nil, nil)

	result, err := uc.Update(ctx, userID, req)

	assert.Error(t, err)
	assert.True(t, userentity.IsUserNotFound(err))
	assert.Nil(t, result)
	repo.AssertNotCalled(t, "Update")
}

func TestUserUseCase_Delete_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	notificationService := &mockNotificationService{}
	tokenGen := &mockTokenGenerator{}

	uc := NewUserUseCase(repo, passwordHasher, notificationService, tokenGen)

	userID := int64(1)
	existingUser := &userentity.User{ID: userID, Name: "Test", Email: "test@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	repo.On("GetByID", ctx, userID).Return(existingUser, nil)
	repo.On("Delete", ctx, userID).Return(nil)

	err := uc.Delete(ctx, userID)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUserUseCase_Delete_UserNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	notificationService := &mockNotificationService{}
	tokenGen := &mockTokenGenerator{}

	uc := NewUserUseCase(repo, passwordHasher, notificationService, tokenGen)

	userID := int64(999)
	repo.On("GetByID", ctx, userID).Return(nil, nil)

	err := uc.Delete(ctx, userID)

	assert.Error(t, err)
	assert.True(t, userentity.IsUserNotFound(err))
	repo.AssertNotCalled(t, "Delete")
}

func TestUserUseCase_Login_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	notificationService := &mockNotificationService{}
	tokenGen := &mockTokenGenerator{}

	uc := NewUserUseCase(repo, passwordHasher, notificationService, tokenGen)

	req := dto.LoginRequest{Email: "test@example.com", Password: "password123"}
	userEntity := &userentity.User{
		ID: 1, Name: "Test", Email: req.Email, Password: "hashed",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	token := "jwt-token"

	repo.On("GetByEmail", ctx, req.Email).Return(userEntity, nil)
	passwordHasher.On("Verify", userEntity.Password, req.Password).Return(true)
	tokenGen.On("Generate", userEntity.ID, userEntity.Email).Return(token, nil)

	result, err := uc.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, token, result.Token)
	assert.Equal(t, userEntity.Email, result.User.Email)
	repo.AssertExpectations(t)
	passwordHasher.AssertExpectations(t)
	tokenGen.AssertExpectations(t)
}

func TestUserUseCase_Login_InvalidCredentials(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	notificationService := &mockNotificationService{}
	tokenGen := &mockTokenGenerator{}

	uc := NewUserUseCase(repo, passwordHasher, notificationService, tokenGen)

	req := dto.LoginRequest{Email: "test@example.com", Password: "wrong"}
	repo.On("GetByEmail", ctx, req.Email).Return(nil, nil)

	result, err := uc.Login(ctx, req)

	assert.Error(t, err)
	assert.True(t, userentity.IsInvalidCredentials(err))
	assert.Nil(t, result)
	passwordHasher.AssertNotCalled(t, "Verify")
	tokenGen.AssertNotCalled(t, "Generate")
}
