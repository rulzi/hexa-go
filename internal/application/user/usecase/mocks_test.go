package usecase

import (
	"context"

	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
	"github.com/stretchr/testify/mock"
)

// mockUserRepository is a mock implementation of Repository
type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) Create(ctx context.Context, user *domainuser.User) (*domainuser.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainuser.User), args.Error(1)
}

func (m *mockUserRepository) GetByID(ctx context.Context, id int64) (*domainuser.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainuser.User), args.Error(1)
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domainuser.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainuser.User), args.Error(1)
}

func (m *mockUserRepository) Update(ctx context.Context, user *domainuser.User) (*domainuser.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainuser.User), args.Error(1)
}

func (m *mockUserRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockUserRepository) List(ctx context.Context, limit, offset int) ([]*domainuser.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domainuser.User), args.Error(1)
}

func (m *mockUserRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// mockPasswordHasher is a mock implementation of PasswordHasher
type mockPasswordHasher struct {
	mock.Mock
}

func (m *mockPasswordHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *mockPasswordHasher) Verify(hashedPassword, password string) bool {
	args := m.Called(hashedPassword, password)
	return args.Bool(0)
}

// mockNotificationService is a mock implementation of NotificationService
type mockNotificationService struct {
	mock.Mock
}

func (m *mockNotificationService) SendWelcomeEmail(ctx context.Context, email, name string) error {
	args := m.Called(ctx, email, name)
	return args.Error(0)
}

// mockTokenGenerator is a mock implementation of TokenGenerator
type mockTokenGenerator struct {
	mock.Mock
}

func (m *mockTokenGenerator) Generate(userID int64, email string) (string, error) {
	args := m.Called(userID, email)
	return args.String(0), args.Error(1)
}

