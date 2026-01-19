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

// Ensure dto is used
var _ *dto.UserResponse

func TestNewCreateUserUseCase(t *testing.T) {
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	notificationService := &mockNotificationService{}

	uc := NewCreateUserUseCase(repo, passwordHasher, notificationService)

	assert.NotNil(t, uc)
	assert.Equal(t, repo, uc.userRepo)
	assert.Equal(t, passwordHasher, uc.passwordHasher)
	assert.Equal(t, notificationService, uc.notificationService)
}

func TestCreateUserUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	notificationService := &mockNotificationService{}

	uc := NewCreateUserUseCase(repo, passwordHasher, notificationService)

	req := dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	hashedPassword := "hashed_password_123"
	expectedUser := &domainuser.User{
		ID:        1,
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("GetByEmail", ctx, req.Email).Return(nil, errors.New("not found"))
	passwordHasher.On("Hash", req.Password).Return(hashedPassword, nil)
	repo.On("Create", ctx, mock.AnythingOfType("*user.User")).Return(expectedUser, nil)
	notificationService.On("SendWelcomeEmail", ctx, req.Email, req.Name).Return(nil)

	result, err := uc.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Name, result.Name)
	assert.Equal(t, expectedUser.Email, result.Email)
	// Password should not be in response (UserResponse doesn't have Password field)

	repo.AssertExpectations(t)
	passwordHasher.AssertExpectations(t)
	notificationService.AssertExpectations(t)
}

func TestCreateUserUseCase_Execute_EmailExists(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	notificationService := &mockNotificationService{}

	uc := NewCreateUserUseCase(repo, passwordHasher, notificationService)

	req := dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	existingUser := &domainuser.User{
		ID:    1,
		Email: req.Email,
	}

	repo.On("GetByEmail", ctx, req.Email).Return(existingUser, nil)

	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, domainuser.ErrEmailExists, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Create")
	notificationService.AssertNotCalled(t, "SendWelcomeEmail")
}

func TestCreateUserUseCase_Execute_GetByEmailErrorContinues(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	notificationService := &mockNotificationService{}

	uc := NewCreateUserUseCase(repo, passwordHasher, notificationService)

	req := dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	hashedPassword := "hashed_password_123"
	expectedUser := &domainuser.User{
		ID:        1,
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// GetByEmail error means email not found, so we continue
	repo.On("GetByEmail", ctx, req.Email).Return(nil, errors.New("database error"))
	passwordHasher.On("Hash", req.Password).Return(hashedPassword, nil)
	repo.On("Create", ctx, mock.AnythingOfType("*user.User")).Return(expectedUser, nil)
	notificationService.On("SendWelcomeEmail", ctx, req.Email, req.Name).Return(nil)

	result, err := uc.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)

	repo.AssertExpectations(t)
	passwordHasher.AssertExpectations(t)
	notificationService.AssertExpectations(t)
}

func TestCreateUserUseCase_Execute_PasswordHashError(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	notificationService := &mockNotificationService{}

	uc := NewCreateUserUseCase(repo, passwordHasher, notificationService)

	req := dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	hashError := errors.New("hash error")

	repo.On("GetByEmail", ctx, req.Email).Return(nil, errors.New("not found"))
	passwordHasher.On("Hash", req.Password).Return("", hashError)

	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, hashError, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
	passwordHasher.AssertExpectations(t)
	repo.AssertNotCalled(t, "Create")
}

func TestCreateUserUseCase_Execute_ValidationError(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	notificationService := &mockNotificationService{}

	uc := NewCreateUserUseCase(repo, passwordHasher, notificationService)

	tests := []struct {
		name string
		req  dto.CreateUserRequest
	}{
		{
			name: "empty name",
			req: dto.CreateUserRequest{
				Name:     "",
				Email:    "test@example.com",
				Password: "password123",
			},
		},
		{
			name: "empty email",
			req: dto.CreateUserRequest{
				Name:     "Test User",
				Email:    "",
				Password: "password123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword := "hashed_password_123"

			repo.On("GetByEmail", ctx, tt.req.Email).Return(nil, errors.New("not found"))
			passwordHasher.On("Hash", tt.req.Password).Return(hashedPassword, nil)

			result, err := uc.Execute(ctx, tt.req)

			assert.Error(t, err)
			assert.Nil(t, result)
			repo.AssertNotCalled(t, "Create")
		})
	}
}

func TestCreateUserUseCase_Execute_RepositoryCreateError(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	notificationService := &mockNotificationService{}

	uc := NewCreateUserUseCase(repo, passwordHasher, notificationService)

	req := dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	hashedPassword := "hashed_password_123"
	repoError := errors.New("database error")

	repo.On("GetByEmail", ctx, req.Email).Return(nil, errors.New("not found"))
	passwordHasher.On("Hash", req.Password).Return(hashedPassword, nil)
	repo.On("Create", ctx, mock.AnythingOfType("*user.User")).Return(nil, repoError)

	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
	passwordHasher.AssertExpectations(t)
	notificationService.AssertNotCalled(t, "SendWelcomeEmail")
}

func TestCreateUserUseCase_Execute_NotificationServiceErrorIgnored(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}
	notificationService := &mockNotificationService{}

	uc := NewCreateUserUseCase(repo, passwordHasher, notificationService)

	req := dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	hashedPassword := "hashed_password_123"
	expectedUser := &domainuser.User{
		ID:        1,
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	notificationError := errors.New("email service error")

	repo.On("GetByEmail", ctx, req.Email).Return(nil, errors.New("not found"))
	passwordHasher.On("Hash", req.Password).Return(hashedPassword, nil)
	repo.On("Create", ctx, mock.AnythingOfType("*user.User")).Return(expectedUser, nil)
	notificationService.On("SendWelcomeEmail", ctx, req.Email, req.Name).Return(notificationError)

	// Notification error should not fail the operation
	result, err := uc.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)

	repo.AssertExpectations(t)
	passwordHasher.AssertExpectations(t)
	notificationService.AssertExpectations(t)
}
