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

func TestNewUpdateUserUseCase(t *testing.T) {
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}

	uc := NewUpdateUserUseCase(repo, passwordHasher)

	assert.NotNil(t, uc)
	assert.Equal(t, repo, uc.userRepo)
	assert.Equal(t, passwordHasher, uc.passwordHasher)
}

func TestUpdateUserUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}

	uc := NewUpdateUserUseCase(repo, passwordHasher)

	userID := int64(1)
	existingUser := &domainuser.User{
		ID:        userID,
		Name:      "Old Name",
		Email:     "old@example.com",
		Password:  "old_hashed_password",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
	}

	req := dto.UpdateUserRequest{
		Name:     "New Name",
		Email:    "new@example.com",
		Password: "", // No password update
	}

	updatedUser := &domainuser.User{
		ID:        userID,
		Name:      req.Name,
		Email:     req.Email,
		Password:  existingUser.Password, // Password unchanged
		CreatedAt: existingUser.CreatedAt,
		UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, userID).Return(existingUser, nil)
	repo.On("GetByEmail", ctx, req.Email).Return(nil, errors.New("not found"))
	repo.On("Update", ctx, mock.AnythingOfType("*user.User")).Return(updatedUser, nil)

	result, err := uc.Execute(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updatedUser.ID, result.ID)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Email, result.Email)

	repo.AssertExpectations(t)
	passwordHasher.AssertNotCalled(t, "Hash")
}

func TestUpdateUserUseCase_Execute_SuccessWithPasswordUpdate(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}

	uc := NewUpdateUserUseCase(repo, passwordHasher)

	userID := int64(1)
	existingUser := &domainuser.User{
		ID:        userID,
		Name:      "Old Name",
		Email:     "old@example.com",
		Password:  "old_hashed_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := dto.UpdateUserRequest{
		Name:     "New Name",
		Email:    "old@example.com", // Same email
		Password: "new_password123",
	}

	newHashedPassword := "new_hashed_password"
	updatedUser := &domainuser.User{
		ID:        userID,
		Name:      req.Name,
		Email:     req.Email,
		Password:  newHashedPassword,
		CreatedAt: existingUser.CreatedAt,
		UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, userID).Return(existingUser, nil)
	passwordHasher.On("Hash", req.Password).Return(newHashedPassword, nil)
	repo.On("Update", ctx, mock.AnythingOfType("*user.User")).Return(updatedUser, nil)

	result, err := uc.Execute(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updatedUser.ID, result.ID)
	assert.Equal(t, req.Name, result.Name)

	repo.AssertExpectations(t)
	passwordHasher.AssertExpectations(t)
}

func TestUpdateUserUseCase_Execute_UserNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}

	uc := NewUpdateUserUseCase(repo, passwordHasher)

	userID := int64(1)
	req := dto.UpdateUserRequest{
		Name:     "New Name",
		Email:    "new@example.com",
		Password: "",
	}

	repo.On("GetByID", ctx, userID).Return(nil, nil)

	result, err := uc.Execute(ctx, userID, req)

	assert.Error(t, err)
	assert.Equal(t, domainuser.ErrUserNotFound, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Update")
}

func TestUpdateUserUseCase_Execute_GetByIDError(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}

	uc := NewUpdateUserUseCase(repo, passwordHasher)

	userID := int64(1)
	req := dto.UpdateUserRequest{
		Name:     "New Name",
		Email:    "new@example.com",
		Password: "",
	}

	repoError := errors.New("database error")
	repo.On("GetByID", ctx, userID).Return(nil, repoError)

	result, err := uc.Execute(ctx, userID, req)

	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Update")
}

func TestUpdateUserUseCase_Execute_EmailExists(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}

	uc := NewUpdateUserUseCase(repo, passwordHasher)

	userID := int64(1)
	existingUser := &domainuser.User{
		ID:        userID,
		Name:      "Old Name",
		Email:     "old@example.com",
		Password:  "hashed_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := dto.UpdateUserRequest{
		Name:     "New Name",
		Email:    "existing@example.com", // Different email
		Password: "",
	}

	existingEmailUser := &domainuser.User{
		ID:    2, // Different user
		Email: req.Email,
	}

	repo.On("GetByID", ctx, userID).Return(existingUser, nil)
	repo.On("GetByEmail", ctx, req.Email).Return(existingEmailUser, nil)

	result, err := uc.Execute(ctx, userID, req)

	assert.Error(t, err)
	assert.Equal(t, domainuser.ErrEmailExists, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Update")
}

func TestUpdateUserUseCase_Execute_SameEmailNoConflict(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}

	uc := NewUpdateUserUseCase(repo, passwordHasher)

	userID := int64(1)
	existingUser := &domainuser.User{
		ID:        userID,
		Name:      "Old Name",
		Email:     "test@example.com",
		Password:  "hashed_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := dto.UpdateUserRequest{
		Name:     "New Name",
		Email:    "test@example.com", // Same email
		Password: "",
	}

	updatedUser := &domainuser.User{
		ID:        userID,
		Name:      req.Name,
		Email:     req.Email,
		Password:  existingUser.Password,
		CreatedAt: existingUser.CreatedAt,
		UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, userID).Return(existingUser, nil)
	// GetByEmail should not be called when email is the same
	repo.On("Update", ctx, mock.AnythingOfType("*user.User")).Return(updatedUser, nil)

	result, err := uc.Execute(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.Name)

	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "GetByEmail")
}

func TestUpdateUserUseCase_Execute_PasswordHashError(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}

	uc := NewUpdateUserUseCase(repo, passwordHasher)

	userID := int64(1)
	existingUser := &domainuser.User{
		ID:        userID,
		Name:      "Old Name",
		Email:     "old@example.com",
		Password:  "old_hashed_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := dto.UpdateUserRequest{
		Name:     "New Name",
		Email:    "old@example.com",
		Password: "new_password123",
	}

	hashError := errors.New("hash error")

	repo.On("GetByID", ctx, userID).Return(existingUser, nil)
	passwordHasher.On("Hash", req.Password).Return("", hashError)

	result, err := uc.Execute(ctx, userID, req)

	assert.Error(t, err)
	assert.Equal(t, hashError, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
	passwordHasher.AssertExpectations(t)
	repo.AssertNotCalled(t, "Update")
}

func TestUpdateUserUseCase_Execute_ValidationError(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}

	uc := NewUpdateUserUseCase(repo, passwordHasher)

	userID := int64(1)
	existingUser := &domainuser.User{
		ID:        userID,
		Name:      "Old Name",
		Email:     "old@example.com",
		Password:  "hashed_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name string
		req  dto.UpdateUserRequest
	}{
		{
			name: "empty name",
			req: dto.UpdateUserRequest{
				Name:     "",
				Email:    "new@example.com",
				Password: "",
			},
		},
		{
			name: "empty email",
			req: dto.UpdateUserRequest{
				Name:     "New Name",
				Email:    "",
				Password: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.On("GetByID", ctx, userID).Return(existingUser, nil)
			if tt.req.Email != existingUser.Email {
				repo.On("GetByEmail", ctx, tt.req.Email).Return(nil, errors.New("not found"))
			}

			result, err := uc.Execute(ctx, userID, tt.req)

			assert.Error(t, err)
			assert.Nil(t, result)
			repo.AssertNotCalled(t, "Update")
		})
	}
}

func TestUpdateUserUseCase_Execute_RepositoryUpdateError(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}
	passwordHasher := &mockPasswordHasher{}

	uc := NewUpdateUserUseCase(repo, passwordHasher)

	userID := int64(1)
	existingUser := &domainuser.User{
		ID:        userID,
		Name:      "Old Name",
		Email:     "old@example.com",
		Password:  "hashed_password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := dto.UpdateUserRequest{
		Name:     "New Name",
		Email:    "new@example.com",
		Password: "",
	}

	updateError := errors.New("update error")

	repo.On("GetByID", ctx, userID).Return(existingUser, nil)
	repo.On("GetByEmail", ctx, req.Email).Return(nil, errors.New("not found"))
	repo.On("Update", ctx, mock.AnythingOfType("*user.User")).Return(nil, updateError)

	result, err := uc.Execute(ctx, userID, req)

	assert.Error(t, err)
	assert.Equal(t, updateError, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
}

