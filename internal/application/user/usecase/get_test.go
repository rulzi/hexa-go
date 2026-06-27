package usecase

import (
	"context"
	"testing"
	"time"

	userentity "github.com/rulzi/hexa-go/internal/domain/user/entity"
	usermocks "github.com/rulzi/hexa-go/internal/domain/user/port/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUser_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := usermocks.NewMockRepository(ctrl)
	passwordHasher := usermocks.NewMockPasswordHasher(ctrl)
	notificationService := usermocks.NewMockNotificationService(ctrl)
	tokenGen := usermocks.NewMockTokenGenerator(ctrl)

	uc := NewUserUseCase(repo, passwordHasher, notificationService, tokenGen)

	userID := int64(1)
	userEntity := &userentity.User{
		ID: userID, Name: "Test", Email: "test@example.com",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	repo.EXPECT().GetByID(ctx, userID).Return(userEntity, nil)

	result, err := uc.Get.Execute(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userEntity.ID, result.ID)
	assert.Equal(t, userEntity.Email, result.Email)
}

func TestGetUser_Execute_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := usermocks.NewMockRepository(ctrl)
	passwordHasher := usermocks.NewMockPasswordHasher(ctrl)
	notificationService := usermocks.NewMockNotificationService(ctrl)
	tokenGen := usermocks.NewMockTokenGenerator(ctrl)

	uc := NewUserUseCase(repo, passwordHasher, notificationService, tokenGen)

	userID := int64(999)
	repo.EXPECT().GetByID(ctx, userID).Return(nil, nil)

	result, err := uc.Get.Execute(ctx, userID)

	assert.Error(t, err)
	assert.True(t, userentity.IsUserNotFound(err))
	assert.Nil(t, result)
}
