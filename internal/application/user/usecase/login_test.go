package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/rulzi/hexa-go/internal/application/user/dto"
	userentity "github.com/rulzi/hexa-go/internal/domain/user/entity"
	usermocks "github.com/rulzi/hexa-go/internal/domain/user/port/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLoginUser_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := usermocks.NewMockRepository(ctrl)
	passwordHasher := usermocks.NewMockPasswordHasher(ctrl)
	tokenGen := usermocks.NewMockTokenGenerator(ctrl)

	uc := NewLoginUser(repo, passwordHasher, tokenGen)

	req := dto.LoginRequest{Email: "test@example.com", Password: "password123"}
	userEntity := &userentity.User{
		ID: 1, Name: "Test", Email: req.Email, Password: "hashed",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	token := "jwt-token"

	repo.EXPECT().GetByEmail(ctx, req.Email).Return(userEntity, nil)
	passwordHasher.EXPECT().Verify(userEntity.Password, req.Password).Return(true)
	tokenGen.EXPECT().Generate(userEntity.ID, userEntity.Email).Return(token, nil)

	result, err := uc.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, token, result.Token)
	assert.Equal(t, userEntity.Email, result.User.Email)
}

func TestLoginUser_Execute_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := usermocks.NewMockRepository(ctrl)
	passwordHasher := usermocks.NewMockPasswordHasher(ctrl)
	tokenGen := usermocks.NewMockTokenGenerator(ctrl)

	uc := NewLoginUser(repo, passwordHasher, tokenGen)

	req := dto.LoginRequest{Email: "test@example.com", Password: "wrong"}
	repo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, nil)

	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.True(t, userentity.IsInvalidCredentials(err))
	assert.Nil(t, result)
}
