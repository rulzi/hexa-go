package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	userentity "github.com/rulzi/hexa-go/internal/domain/user/entity"
	usermocks "github.com/rulzi/hexa-go/internal/domain/user/port/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestListUser_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := usermocks.NewMockRepository(ctrl)

	uc := NewListUser(repo)

	limit, offset := 10, 0
	users := []*userentity.User{
		{ID: 1, Name: "User1", Email: "u1@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	total := int64(1)
	repo.EXPECT().List(ctx, limit, offset).Return(users, nil)
	repo.EXPECT().Count(ctx).Return(total, nil)

	result, err := uc.Execute(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, total, result.Total)
	assert.Equal(t, limit, result.Limit)
	assert.Len(t, result.Users, 1)
}

func TestListUser_Execute_DefaultPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := usermocks.NewMockRepository(ctrl)
	uc := NewListUser(repo)

	repo.EXPECT().List(ctx, 10, 0).Return([]*userentity.User{}, nil)
	repo.EXPECT().Count(ctx).Return(int64(0), nil)

	result, err := uc.Execute(ctx, 0, -1)

	assert.NoError(t, err)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 0, result.Offset)
}

func TestListUser_Execute_ListError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := usermocks.NewMockRepository(ctrl)
	uc := NewListUser(repo)

	repo.EXPECT().List(ctx, 10, 0).Return(nil, errors.New("db error"))

	result, err := uc.Execute(ctx, 10, 0)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestListUser_Execute_CountError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := usermocks.NewMockRepository(ctrl)
	uc := NewListUser(repo)

	repo.EXPECT().List(ctx, 10, 0).Return([]*userentity.User{}, nil)
	repo.EXPECT().Count(ctx).Return(int64(0), errors.New("count error"))

	result, err := uc.Execute(ctx, 10, 0)

	assert.Error(t, err)
	assert.Nil(t, result)
}
