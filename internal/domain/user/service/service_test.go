package service

import (
	"context"
	"testing"

	"github.com/rulzi/hexa-go/internal/domain/user/entity"
	"github.com/rulzi/hexa-go/internal/domain/user/port"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	createFunc     func(ctx context.Context, user *entity.User) (*entity.User, error)
	getByIDFunc    func(ctx context.Context, id int64) (*entity.User, error)
	getByEmailFunc func(ctx context.Context, email string) (*entity.User, error)
	updateFunc     func(ctx context.Context, user *entity.User) (*entity.User, error)
	deleteFunc     func(ctx context.Context, id int64) error
	listFunc       func(ctx context.Context, limit, offset int) ([]*entity.User, error)
	countFunc      func(ctx context.Context) (int64, error)
}

func (m *mockRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, user)
	}
	return nil, nil
}

func (m *mockRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	if m.getByEmailFunc != nil {
		return m.getByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *mockRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, user)
	}
	return nil, nil
}

func (m *mockRepository) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, limit, offset)
	}
	return nil, nil
}

func (m *mockRepository) Count(ctx context.Context) (int64, error) {
	if m.countFunc != nil {
		return m.countFunc(ctx)
	}
	return 0, nil
}

type mockTokenGenerator struct {
	generateFunc func(userID int64, email string) (string, error)
}

func (m *mockTokenGenerator) Generate(userID int64, email string) (string, error) {
	if m.generateFunc != nil {
		return m.generateFunc(userID, email)
	}
	return "mock-token", nil
}

type mockTokenValidator struct {
	validateFunc func(token string) (*port.TokenClaims, error)
}

func (m *mockTokenValidator) Validate(token string) (*port.TokenClaims, error) {
	if m.validateFunc != nil {
		return m.validateFunc(token)
	}
	return &port.TokenClaims{UserID: 1, Email: "test@example.com"}, nil
}

type mockPasswordHasher struct {
	hashFunc   func(password string) (string, error)
	verifyFunc func(hashedPassword, password string) bool
}

func (m *mockPasswordHasher) Hash(password string) (string, error) {
	if m.hashFunc != nil {
		return m.hashFunc(password)
	}
	return "hashed-password", nil
}

func (m *mockPasswordHasher) Verify(hashedPassword, password string) bool {
	if m.verifyFunc != nil {
		return m.verifyFunc(hashedPassword, password)
	}
	return true
}

func TestNewService(t *testing.T) {
	tests := []struct {
		name           string
		repo           port.Repository
		tokenGen       port.TokenGenerator
		tokenValidator port.TokenValidator
		passwordHasher port.PasswordHasher
	}{
		{
			name:           "create service with all dependencies",
			repo:           &mockRepository{},
			tokenGen:       &mockTokenGenerator{},
			tokenValidator: &mockTokenValidator{},
			passwordHasher: &mockPasswordHasher{},
		},
		{
			name:           "create service with nil dependencies",
			repo:           nil,
			tokenGen:       nil,
			tokenValidator: nil,
			passwordHasher: nil,
		},
		{
			name:           "create service with partial dependencies",
			repo:           &mockRepository{},
			tokenGen:       nil,
			tokenValidator: &mockTokenValidator{},
			passwordHasher: &mockPasswordHasher{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(tt.repo, tt.tokenGen, tt.tokenValidator, tt.passwordHasher)
			assert.NotNil(t, svc)
			assert.Equal(t, tt.repo, svc.repo)
			assert.Equal(t, tt.tokenGen, svc.tokenGen)
			assert.Equal(t, tt.tokenValidator, svc.tokenValidator)
			assert.Equal(t, tt.passwordHasher, svc.passwordHasher)
		})
	}
}
