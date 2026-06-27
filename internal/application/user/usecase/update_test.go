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
