package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	userentity "github.com/rulzi/hexa-go/internal/domain/user/entity"
	usermocks "github.com/rulzi/hexa-go/internal/domain/user/port/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthenticate(t *testing.T) {
	ctx := context.Background()
	tokenErr := errors.New("token generation failed")

	tests := []struct {
		name         string
		email        string
		password     string
		setup        func(repo *usermocks.MockRepository, hasher *usermocks.MockPasswordHasher, tokenGen *usermocks.MockTokenGenerator)
		wantErrCheck func(error) bool
		assertResult func(t *testing.T, result *Result)
	}{
		{
			name:     "success",
			email:    "test@example.com",
			password: "password123",
			setup: func(repo *usermocks.MockRepository, hasher *usermocks.MockPasswordHasher, tokenGen *usermocks.MockTokenGenerator) {
				user := &userentity.User{
					ID: 1, Name: "Test", Email: "test@example.com", Password: "hashed",
					CreatedAt: time.Now(), UpdatedAt: time.Now(),
				}
				repo.EXPECT().GetByEmail(ctx, "test@example.com").Return(user, nil)
				hasher.EXPECT().Verify("hashed", "password123").Return(true)
				tokenGen.EXPECT().Generate(int64(1), "test@example.com").Return("jwt-token", nil)
			},
			assertResult: func(t *testing.T, result *Result) {
				require.NotNil(t, result)
				assert.Equal(t, "jwt-token", result.Token)
				assert.Equal(t, int64(1), result.User.ID)
			},
		},
		{
			name:     "user not found",
			email:    "missing@example.com",
			password: "password123",
			setup: func(repo *usermocks.MockRepository, hasher *usermocks.MockPasswordHasher, tokenGen *usermocks.MockTokenGenerator) {
				repo.EXPECT().GetByEmail(ctx, "missing@example.com").Return(nil, nil)
			},
			wantErrCheck: userentity.IsInvalidCredentials,
		},
		{
			name:     "repository error",
			email:    "test@example.com",
			password: "password123",
			setup: func(repo *usermocks.MockRepository, hasher *usermocks.MockPasswordHasher, tokenGen *usermocks.MockTokenGenerator) {
				repo.EXPECT().GetByEmail(ctx, "test@example.com").Return(nil, errors.New("db down"))
			},
			wantErrCheck: userentity.IsInvalidCredentials,
		},
		{
			name:     "wrong password",
			email:    "test@example.com",
			password: "wrong",
			setup: func(repo *usermocks.MockRepository, hasher *usermocks.MockPasswordHasher, tokenGen *usermocks.MockTokenGenerator) {
				user := &userentity.User{ID: 1, Email: "test@example.com", Password: "hashed"}
				repo.EXPECT().GetByEmail(ctx, "test@example.com").Return(user, nil)
				hasher.EXPECT().Verify("hashed", "wrong").Return(false)
			},
			wantErrCheck: userentity.IsInvalidCredentials,
		},
		{
			name:     "token generation failed",
			email:    "test@example.com",
			password: "password123",
			setup: func(repo *usermocks.MockRepository, hasher *usermocks.MockPasswordHasher, tokenGen *usermocks.MockTokenGenerator) {
				user := &userentity.User{ID: 1, Email: "test@example.com", Password: "hashed"}
				repo.EXPECT().GetByEmail(ctx, "test@example.com").Return(user, nil)
				hasher.EXPECT().Verify("hashed", "password123").Return(true)
				tokenGen.EXPECT().Generate(int64(1), "test@example.com").Return("", tokenErr)
			},
			wantErrCheck: func(err error) bool { return errors.Is(err, tokenErr) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := usermocks.NewMockRepository(ctrl)
			hasher := usermocks.NewMockPasswordHasher(ctrl)
			tokenGen := usermocks.NewMockTokenGenerator(ctrl)

			if tt.setup != nil {
				tt.setup(repo, hasher, tokenGen)
			}

			result, err := Authenticate(ctx, tt.email, tt.password, repo, hasher, tokenGen)

			if tt.wantErrCheck != nil {
				require.Error(t, err)
				assert.True(t, tt.wantErrCheck(err))
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err)
			if tt.assertResult != nil {
				tt.assertResult(t, result)
			}
		})
	}
}
