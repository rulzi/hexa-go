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
