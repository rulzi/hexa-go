package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rulzi/hexa-go/internal/application/user/dto"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
	usermocks "github.com/rulzi/hexa-go/internal/domain/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserUseCase_Create(t *testing.T) {
	ctx := context.Background()
	dbDownErr := errors.New("database connection refused")

	tests := []struct {
		name          string
		req           dto.CreateUserRequest
		setupRepo     func(repo *usermocks.MockRepository)
		setupHasher   func(hasher *mockPasswordHasher)
		setupNotifier func(notifier *mockNotificationService)
		wantErrCheck  func(error) bool
		assertResult  func(t *testing.T, resp *dto.UserResponse)
	}{
		{
			name: "happy path - user created successfully",
			req: dto.CreateUserRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupRepo: func(repo *usermocks.MockRepository) {
				created := &domainuser.User{
					ID:        1,
					Name:      "Test User",
					Email:     "test@example.com",
					Password:  "hashed_123",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				repo.EXPECT().GetByEmail(ctx, "test@example.com").Return(nil, domainuser.NewUserNotFound())
				repo.EXPECT().Create(ctx, gomock.AssignableToTypeOf(&domainuser.User{})).Return(created, nil)
			},
			setupHasher: func(hasher *mockPasswordHasher) {
				hasher.On("Hash", "password123").Return("hashed_123", nil)
			},
			setupNotifier: func(notifier *mockNotificationService) {
				notifier.On("SendWelcomeEmail", ctx, "test@example.com", "Test User").Return(nil)
			},
			wantErrCheck: nil,
			assertResult: func(t *testing.T, resp *dto.UserResponse) {
				require.NotNil(t, resp)
				assert.Equal(t, int64(1), resp.ID)
				assert.Equal(t, "Test User", resp.Name)
				assert.Equal(t, "test@example.com", resp.Email)
			},
		},
		{
			name: "email already registered",
			req: dto.CreateUserRequest{
				Name:     "Duplicate User",
				Email:    "exist@example.com",
				Password: "password123",
			},
			setupRepo: func(repo *usermocks.MockRepository) {
				existing := &domainuser.User{ID: 42, Email: "exist@example.com"}
				repo.EXPECT().GetByEmail(ctx, "exist@example.com").Return(existing, nil)
			},
			setupHasher: func(hasher *mockPasswordHasher) {},
			wantErrCheck: domainuser.IsEmailExists,
			assertResult: func(t *testing.T, resp *dto.UserResponse) {
				assert.Nil(t, resp)
			},
		},
		{
			name: "validation failed - empty email",
			req: dto.CreateUserRequest{
				Name:     "Test User",
				Email:    "",
				Password: "password123",
			},
			setupRepo:   func(repo *usermocks.MockRepository) {},
			setupHasher: func(hasher *mockPasswordHasher) {},
			wantErrCheck: domainuser.IsEmailRequired,
			assertResult: func(t *testing.T, resp *dto.UserResponse) {
				assert.Nil(t, resp)
			},
		},
		{
			name: "validation failed - password too short",
			req: dto.CreateUserRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "12345",
			},
			setupRepo:   func(repo *usermocks.MockRepository) {},
			setupHasher: func(hasher *mockPasswordHasher) {},
			wantErrCheck: domainuser.IsPasswordTooShort,
			assertResult: func(t *testing.T, resp *dto.UserResponse) {
				assert.Nil(t, resp)
			},
		},
		{
			name: "repository error - database down on create",
			req: dto.CreateUserRequest{
				Name:     "Test User",
				Email:    "new@example.com",
				Password: "password123",
			},
			setupRepo: func(repo *usermocks.MockRepository) {
				repo.EXPECT().GetByEmail(ctx, "new@example.com").Return(nil, domainuser.NewUserNotFound())
				repo.EXPECT().Create(ctx, gomock.AssignableToTypeOf(&domainuser.User{})).Return(nil, dbDownErr)
			},
			setupHasher: func(hasher *mockPasswordHasher) {
				hasher.On("Hash", "password123").Return("hashed_123", nil)
			},
			wantErrCheck: func(err error) bool { return errors.Is(err, dbDownErr) },
			assertResult: func(t *testing.T, resp *dto.UserResponse) {
				assert.Nil(t, resp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := usermocks.NewMockRepository(ctrl)
			hasher := &mockPasswordHasher{}
			notifier := &mockNotificationService{}
			tokenGen := &mockTokenGenerator{}

			if tt.setupRepo != nil {
				tt.setupRepo(repo)
			}
			if tt.setupHasher != nil {
				tt.setupHasher(hasher)
			}
			if tt.setupNotifier != nil {
				tt.setupNotifier(notifier)
			}

			uc := NewUserUseCase(repo, hasher, notifier, tokenGen)
			resp, err := uc.Create(ctx, tt.req)

			if tt.wantErrCheck != nil {
				require.Error(t, err)
				assert.True(t, tt.wantErrCheck(err))
			} else {
				require.NoError(t, err)
			}

			if tt.assertResult != nil {
				tt.assertResult(t, resp)
			}

			hasher.AssertExpectations(t)
			notifier.AssertExpectations(t)
		})
	}
}

func TestUserUseCase_Create_ValidationDoesNotTouchRepository(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := usermocks.NewMockRepository(ctrl)
	hasher := &mockPasswordHasher{}
	notifier := &mockNotificationService{}

	uc := NewUserUseCase(repo, hasher, notifier, &mockTokenGenerator{})

	invalidRequests := []dto.CreateUserRequest{
		{Name: "Test", Email: "", Password: "password123"},
		{Name: "Test", Email: "test@example.com", Password: "short"},
	}

	for _, req := range invalidRequests {
		resp, err := uc.Create(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	}

	hasher.AssertNotCalled(t, "Hash", mock.Anything)
	notifier.AssertNotCalled(t, "SendWelcomeEmail", mock.Anything, mock.Anything, mock.Anything)
}

func TestUserUseCase_Create_EmailExistsDoesNotPersist(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := usermocks.NewMockRepository(ctrl)
	hasher := &mockPasswordHasher{}
	notifier := &mockNotificationService{}

	existing := &domainuser.User{ID: 1, Email: "exist@example.com"}
	repo.EXPECT().GetByEmail(ctx, "exist@example.com").Return(existing, nil)

	uc := NewUserUseCase(repo, hasher, notifier, &mockTokenGenerator{})
	resp, err := uc.Create(ctx, dto.CreateUserRequest{
		Name:     "Test",
		Email:    "exist@example.com",
		Password: "password123",
	})

	assert.True(t, domainuser.IsEmailExists(err))
	assert.Nil(t, resp)
	hasher.AssertNotCalled(t, "Hash", mock.Anything)
	notifier.AssertNotCalled(t, "SendWelcomeEmail", mock.Anything, mock.Anything, mock.Anything)
}
