package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rulzi/hexa-go/internal/application/user/dto"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
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

func TestUserUseCase_Create_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	notificationService := &mockNotificationService{}
	tokenGen := &mockTokenGenerator{}

	uc := NewUserUseCase(repo, passwordHasher, notificationService, tokenGen)

	req := dto.CreateUserRequest{Name: "Test User", Email: "test@example.com", Password: "password123"}
	hashedPassword := "hashed_123"
	expectedUser := &domainuser.User{
		ID: 1, Name: req.Name, Email: req.Email, Password: hashedPassword,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	repo.On("GetByEmail", ctx, req.Email).Return(nil, errors.New("not found"))
	passwordHasher.On("Hash", req.Password).Return(hashedPassword, nil)
	repo.On("Create", ctx, mock.AnythingOfType("*user.User")).Return(expectedUser, nil)
	notificationService.On("SendWelcomeEmail", ctx, req.Email, req.Name).Return(nil)

	result, err := uc.Create(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Email, result.Email)
	repo.AssertExpectations(t)
	passwordHasher.AssertExpectations(t)
	notificationService.AssertExpectations(t)
}

func TestUserUseCase_Create_EmailExists(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	notificationService := &mockNotificationService{}
	tokenGen := &mockTokenGenerator{}

	uc := NewUserUseCase(repo, passwordHasher, notificationService, tokenGen)

	req := dto.CreateUserRequest{Name: "Test", Email: "exist@example.com", Password: "pass123"}
	existingUser := &domainuser.User{ID: 1, Email: req.Email}
	repo.On("GetByEmail", ctx, req.Email).Return(existingUser, nil)

	result, err := uc.Create(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, domainuser.ErrEmailExists, err)
	assert.Nil(t, result)
	passwordHasher.AssertNotCalled(t, "Hash")
	repo.AssertNotCalled(t, "Create")
}

func TestUserUseCase_Get_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	notificationService := &mockNotificationService{}
	tokenGen := &mockTokenGenerator{}

	uc := NewUserUseCase(repo, passwordHasher, notificationService, tokenGen)

	userID := int64(1)
	userEntity := &domainuser.User{
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
	assert.Equal(t, domainuser.ErrUserNotFound, err)
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
	users := []*domainuser.User{
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
	existingUser := &domainuser.User{
		ID: userID, Name: "Old", Email: "old@example.com", Password: "hash",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	req := dto.UpdateUserRequest{Name: "New", Email: "new@example.com"}
	updatedUser := &domainuser.User{
		ID: userID, Name: req.Name, Email: req.Email, Password: existingUser.Password,
		CreatedAt: existingUser.CreatedAt, UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, userID).Return(existingUser, nil)
	repo.On("GetByEmail", ctx, req.Email).Return(nil, errors.New("not found"))
	repo.On("Update", ctx, mock.AnythingOfType("*user.User")).Return(updatedUser, nil)

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
	assert.Equal(t, domainuser.ErrUserNotFound, err)
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
	existingUser := &domainuser.User{ID: userID, Name: "Test", Email: "test@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}
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
	assert.Equal(t, domainuser.ErrUserNotFound, err)
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
	userEntity := &domainuser.User{
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
	assert.Equal(t, domainuser.ErrInvalidCredentials, err)
	assert.Nil(t, result)
	passwordHasher.AssertNotCalled(t, "Verify")
	tokenGen.AssertNotCalled(t, "Generate")
}
