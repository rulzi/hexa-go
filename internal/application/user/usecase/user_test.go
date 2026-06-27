package usecase

import (
	"testing"

	usermocks "github.com/rulzi/hexa-go/internal/domain/user/port/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewUserUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := usermocks.NewMockRepository(ctrl)
	passwordHasher := usermocks.NewMockPasswordHasher(ctrl)
	notificationService := usermocks.NewMockNotificationService(ctrl)
	tokenGen := usermocks.NewMockTokenGenerator(ctrl)

	uc := NewUserUseCase(repo, passwordHasher, notificationService, tokenGen)

	assert.NotNil(t, uc)
	assert.NotNil(t, uc.Create)
	assert.NotNil(t, uc.Get)
	assert.NotNil(t, uc.List)
	assert.NotNil(t, uc.Update)
	assert.NotNil(t, uc.Delete)
	assert.NotNil(t, uc.Login)
}
