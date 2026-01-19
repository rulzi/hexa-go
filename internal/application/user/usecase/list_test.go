package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rulzi/hexa-go/internal/application/user/dto"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
	"github.com/stretchr/testify/assert"
)

// Ensure dto is used
var _ *dto.ListUsersResponse

func TestNewListUsersUseCase(t *testing.T) {
	repo := &mockUserRepository{}

	uc := NewListUsersUseCase(repo)

	assert.NotNil(t, uc)
	assert.Equal(t, repo, uc.userRepo)
}

func TestListUsersUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}

	uc := NewListUsersUseCase(repo)

	limit := 10
	offset := 0
	users := []*domainuser.User{
		{
			ID:        1,
			Name:      "User 1",
			Email:     "user1@example.com",
			Password:  "hashed1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Name:      "User 2",
			Email:     "user2@example.com",
			Password:  "hashed2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	total := int64(2)

	repo.On("List", ctx, limit, offset).Return(users, nil)
	repo.On("Count", ctx).Return(total, nil)

	result, err := uc.Execute(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, total, result.Total)
	assert.Equal(t, limit, result.Limit)
	assert.Equal(t, offset, result.Offset)
	assert.Equal(t, len(users), len(result.Users))
	assert.Equal(t, users[0].ID, result.Users[0].ID)
	assert.Equal(t, users[1].ID, result.Users[1].ID)
	// Passwords should not be in response (UserResponse doesn't have Password field)

	repo.AssertExpectations(t)
}

func TestListUsersUseCase_Execute_DefaultPagination(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}

	uc := NewListUsersUseCase(repo)

	users := []*domainuser.User{}
	total := int64(0)

	repo.On("List", ctx, 10, 0).Return(users, nil)
	repo.On("Count", ctx).Return(total, nil)

	result, err := uc.Execute(ctx, -1, -1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 0, result.Offset)

	repo.AssertExpectations(t)
}

func TestListUsersUseCase_Execute_ListError(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}

	uc := NewListUsersUseCase(repo)

	limit := 10
	offset := 0
	listError := errors.New("list error")

	repo.On("List", ctx, limit, offset).Return(nil, listError)

	result, err := uc.Execute(ctx, limit, offset)

	assert.Error(t, err)
	assert.Equal(t, listError, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Count")
}

func TestListUsersUseCase_Execute_CountError(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}

	uc := NewListUsersUseCase(repo)

	limit := 10
	offset := 0
	users := []*domainuser.User{
		{
			ID:        1,
			Name:      "User 1",
			Email:     "user1@example.com",
			Password:  "hashed1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	countError := errors.New("count error")

	repo.On("List", ctx, limit, offset).Return(users, nil)
	repo.On("Count", ctx).Return(int64(0), countError)

	result, err := uc.Execute(ctx, limit, offset)

	assert.Error(t, err)
	assert.Equal(t, countError, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
}

func TestListUsersUseCase_Execute_EmptyList(t *testing.T) {
	ctx := context.Background()
	repo := &mockUserRepository{}

	uc := NewListUsersUseCase(repo)

	limit := 10
	offset := 0
	users := []*domainuser.User{}
	total := int64(0)

	repo.On("List", ctx, limit, offset).Return(users, nil)
	repo.On("Count", ctx).Return(total, nil)

	result, err := uc.Execute(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(0), result.Total)
	assert.Equal(t, 0, len(result.Users))

	repo.AssertExpectations(t)
}

