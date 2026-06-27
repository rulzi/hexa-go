package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rulzi/hexa-go/internal/application/user/dto"
	userentity "github.com/rulzi/hexa-go/internal/domain/user/entity"
	usermocks "github.com/rulzi/hexa-go/internal/domain/user/port/mocks"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
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
		setupHasher   func(hasher *usermocks.MockPasswordHasher)
		setupNotifier func(notifier *usermocks.MockNotificationService)
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
				created := &userentity.User{
					ID:        1,
					Name:      "Test User",
					Email:     "test@example.com",
					Password:  "hashed_123",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				repo.EXPECT().GetByEmail(ctx, "test@example.com").Return(nil, userentity.NewUserNotFound())
				repo.EXPECT().Create(ctx, gomock.AssignableToTypeOf(&userentity.User{})).Return(created, nil)
			},
			setupHasher: func(hasher *usermocks.MockPasswordHasher) {
				hasher.EXPECT().Hash("password123").Return("hashed_123", nil)
			},
			setupNotifier: func(notifier *usermocks.MockNotificationService) {
				notifier.EXPECT().SendWelcomeEmail(ctx, "test@example.com", "Test User").Return(nil)
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
				existing := &userentity.User{ID: 42, Email: "exist@example.com"}
				repo.EXPECT().GetByEmail(ctx, "exist@example.com").Return(existing, nil)
			},
			setupHasher:  func(hasher *usermocks.MockPasswordHasher) {},
			wantErrCheck: userentity.IsEmailExists,
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
			setupRepo:    func(repo *usermocks.MockRepository) {},
			setupHasher:  func(hasher *usermocks.MockPasswordHasher) {},
			wantErrCheck: userentity.IsEmailRequired,
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
			setupRepo:    func(repo *usermocks.MockRepository) {},
			setupHasher:  func(hasher *usermocks.MockPasswordHasher) {},
			wantErrCheck: userentity.IsPasswordTooShort,
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
				repo.EXPECT().GetByEmail(ctx, "new@example.com").Return(nil, userentity.NewUserNotFound())
				repo.EXPECT().Create(ctx, gomock.AssignableToTypeOf(&userentity.User{})).Return(nil, dbDownErr)
			},
			setupHasher: func(hasher *usermocks.MockPasswordHasher) {
				hasher.EXPECT().Hash("password123").Return("hashed_123", nil)
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
			hasher := usermocks.NewMockPasswordHasher(ctrl)
			notifier := usermocks.NewMockNotificationService(ctrl)

			if tt.setupRepo != nil {
				tt.setupRepo(repo)
			}
			if tt.setupHasher != nil {
				tt.setupHasher(hasher)
			}
			if tt.setupNotifier != nil {
				tt.setupNotifier(notifier)
			}

			uc := NewCreateUser(repo, hasher, notifier, logger.NewSimpleLogger())
			resp, err := uc.Execute(ctx, tt.req)

			if tt.wantErrCheck != nil {
				require.Error(t, err)
				assert.True(t, tt.wantErrCheck(err))
			} else {
				require.NoError(t, err)
			}

			if tt.assertResult != nil {
				tt.assertResult(t, resp)
			}
		})
	}
}

func TestUserUseCase_Create_ValidationDoesNotTouchRepository(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := usermocks.NewMockRepository(ctrl)
	hasher := usermocks.NewMockPasswordHasher(ctrl)
	notifier := usermocks.NewMockNotificationService(ctrl)

	uc := NewCreateUser(repo, hasher, notifier, logger.NewSimpleLogger())

	invalidRequests := []dto.CreateUserRequest{
		{Name: "Test", Email: "", Password: "password123"},
		{Name: "Test", Email: "test@example.com", Password: "short"},
	}

	for _, req := range invalidRequests {
		resp, err := uc.Execute(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	}
}

func TestUserUseCase_Create_EmailExistsDoesNotPersist(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := usermocks.NewMockRepository(ctrl)
	hasher := usermocks.NewMockPasswordHasher(ctrl)
	notifier := usermocks.NewMockNotificationService(ctrl)

	existing := &userentity.User{ID: 1, Email: "exist@example.com"}
	repo.EXPECT().GetByEmail(ctx, "exist@example.com").Return(existing, nil)

	uc := NewCreateUser(repo, hasher, notifier, logger.NewSimpleLogger())
	resp, err := uc.Execute(ctx, dto.CreateUserRequest{
		Name:     "Test",
		Email:    "exist@example.com",
		Password: "password123",
	})

	assert.True(t, userentity.IsEmailExists(err))
	assert.Nil(t, resp)
}
