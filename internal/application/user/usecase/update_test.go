package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rulzi/hexa-go/internal/application/user/dto"
	userentity "github.com/rulzi/hexa-go/internal/domain/user/entity"
	usermocks "github.com/rulzi/hexa-go/internal/domain/user/port/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpdateUser_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := usermocks.NewMockRepository(ctrl)
	passwordHasher := usermocks.NewMockPasswordHasher(ctrl)

	uc := NewUpdateUser(repo, passwordHasher)

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

	repo.EXPECT().GetByID(ctx, userID).Return(existingUser, nil)
	repo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, errors.New("not found"))
	repo.EXPECT().Update(ctx, gomock.AssignableToTypeOf(&userentity.User{})).Return(updatedUser, nil)

	result, err := uc.Execute(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Email, result.Email)
}

func TestUpdateUser_Execute_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := usermocks.NewMockRepository(ctrl)
	passwordHasher := usermocks.NewMockPasswordHasher(ctrl)

	uc := NewUpdateUser(repo, passwordHasher)

	userID := int64(999)
	req := dto.UpdateUserRequest{Name: "New", Email: "new@example.com"}
	repo.EXPECT().GetByID(ctx, userID).Return(nil, nil)

	result, err := uc.Execute(ctx, userID, req)

	assert.Error(t, err)
	assert.True(t, userentity.IsUserNotFound(err))
	assert.Nil(t, result)
}

func TestUpdateUser_Execute_EmailAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := usermocks.NewMockRepository(ctrl)
	passwordHasher := usermocks.NewMockPasswordHasher(ctrl)
	uc := NewUpdateUser(repo, passwordHasher)

	userID := int64(1)
	existingUser := &userentity.User{
		ID: userID, Name: "Old", Email: "old@example.com", Password: "hash",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	req := dto.UpdateUserRequest{Name: "New", Email: "taken@example.com"}

	repo.EXPECT().GetByID(ctx, userID).Return(existingUser, nil)
	repo.EXPECT().GetByEmail(ctx, req.Email).Return(&userentity.User{ID: 2}, nil)

	result, err := uc.Execute(ctx, userID, req)

	assert.Error(t, err)
	assert.True(t, userentity.IsEmailExists(err))
	assert.Nil(t, result)
}

func TestUpdateUser_Execute_WithPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := usermocks.NewMockRepository(ctrl)
	passwordHasher := usermocks.NewMockPasswordHasher(ctrl)
	uc := NewUpdateUser(repo, passwordHasher)

	userID := int64(1)
	existingUser := &userentity.User{
		ID: userID, Name: "User", Email: "user@example.com", Password: "old-hash",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	req := dto.UpdateUserRequest{Name: "User", Email: "user@example.com", Password: "newpassword123"}
	updatedUser := &userentity.User{
		ID: userID, Name: req.Name, Email: req.Email, Password: "new-hash",
		CreatedAt: existingUser.CreatedAt, UpdatedAt: time.Now(),
	}

	repo.EXPECT().GetByID(ctx, userID).Return(existingUser, nil)
	passwordHasher.EXPECT().Hash(req.Password).Return("new-hash", nil)
	repo.EXPECT().Update(ctx, gomock.AssignableToTypeOf(&userentity.User{})).Return(updatedUser, nil)

	result, err := uc.Execute(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestUpdateUser_Execute_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := usermocks.NewMockRepository(ctrl)
	passwordHasher := usermocks.NewMockPasswordHasher(ctrl)
	uc := NewUpdateUser(repo, passwordHasher)

	userID := int64(1)
	existingUser := &userentity.User{
		ID: userID, Name: "User", Email: "user@example.com", Password: "hash",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	req := dto.UpdateUserRequest{Name: "User", Email: "user@example.com", Password: "short"}

	repo.EXPECT().GetByID(ctx, userID).Return(existingUser, nil)

	result, err := uc.Execute(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdateUser_Execute_HashError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := usermocks.NewMockRepository(ctrl)
	passwordHasher := usermocks.NewMockPasswordHasher(ctrl)
	uc := NewUpdateUser(repo, passwordHasher)

	userID := int64(1)
	existingUser := &userentity.User{
		ID: userID, Name: "User", Email: "user@example.com", Password: "hash",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	req := dto.UpdateUserRequest{Name: "User", Email: "user@example.com", Password: "validpassword123"}

	repo.EXPECT().GetByID(ctx, userID).Return(existingUser, nil)
	passwordHasher.EXPECT().Hash(req.Password).Return("", errors.New("hash failed"))

	result, err := uc.Execute(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
}
