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

func TestDeleteUser_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := usermocks.NewMockRepository(ctrl)

	uc := NewDeleteUser(repo)

	userID := int64(1)
	existingUser := &userentity.User{ID: userID, Name: "Test", Email: "test@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	repo.EXPECT().GetByID(ctx, userID).Return(existingUser, nil)
	repo.EXPECT().Delete(ctx, userID).Return(nil)

	err := uc.Execute(ctx, userID)

	assert.NoError(t, err)
}

func TestDeleteUser_Execute_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := usermocks.NewMockRepository(ctrl)

	uc := NewDeleteUser(repo)

	userID := int64(999)
	repo.EXPECT().GetByID(ctx, userID).Return(nil, nil)

	err := uc.Execute(ctx, userID)

	assert.Error(t, err)
	assert.True(t, userentity.IsUserNotFound(err))
}
