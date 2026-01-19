package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockRepository is a mock implementation of Repository for testing
type mockRepository struct {
	createFunc    func(ctx context.Context, user *User) (*User, error)
	getByIDFunc   func(ctx context.Context, id int64) (*User, error)
	getByEmailFunc func(ctx context.Context, email string) (*User, error)
	updateFunc    func(ctx context.Context, user *User) (*User, error)
	deleteFunc    func(ctx context.Context, id int64) error
	listFunc      func(ctx context.Context, limit, offset int) ([]*User, error)
	countFunc     func(ctx context.Context) (int64, error)
}

func (m *mockRepository) Create(ctx context.Context, user *User) (*User, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, user)
	}
	return nil, nil
}

func (m *mockRepository) GetByID(ctx context.Context, id int64) (*User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	if m.getByEmailFunc != nil {
		return m.getByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *mockRepository) Update(ctx context.Context, user *User) (*User, error) {
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

func (m *mockRepository) List(ctx context.Context, limit, offset int) ([]*User, error) {
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

// mockTokenGenerator is a mock implementation of TokenGenerator for testing
type mockTokenGenerator struct {
	generateFunc func(userID int64, email string) (string, error)
}

func (m *mockTokenGenerator) Generate(userID int64, email string) (string, error) {
	if m.generateFunc != nil {
		return m.generateFunc(userID, email)
	}
	return "mock-token", nil
}

// mockTokenValidator is a mock implementation of TokenValidator for testing
type mockTokenValidator struct {
	validateFunc func(token string) (*TokenClaims, error)
}

func (m *mockTokenValidator) Validate(token string) (*TokenClaims, error) {
	if m.validateFunc != nil {
		return m.validateFunc(token)
	}
	return &TokenClaims{UserID: 1, Email: "test@example.com"}, nil
}

// mockPasswordHasher is a mock implementation of PasswordHasher for testing
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
		repo           Repository
		tokenGen       TokenGenerator
		tokenValidator TokenValidator
		passwordHasher PasswordHasher
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
			service := NewService(tt.repo, tt.tokenGen, tt.tokenValidator, tt.passwordHasher)
			assert.NotNil(t, service)
			assert.Equal(t, tt.repo, service.repo)
			assert.Equal(t, tt.tokenGen, service.tokenGen)
			assert.Equal(t, tt.tokenValidator, service.tokenValidator)
			assert.Equal(t, tt.passwordHasher, service.passwordHasher)
		})
	}
}

