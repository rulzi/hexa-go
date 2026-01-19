package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr error
	}{
		{
			name: "valid user",
			user: User{
				ID:        1,
				Name:      "John Doe",
				Email:     "john@example.com",
				Password:  "password123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: nil,
		},
		{
			name: "missing name",
			user: User{
				ID:        1,
				Name:      "",
				Email:     "john@example.com",
				Password:  "password123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: ErrNameRequired,
		},
		{
			name: "missing email",
			user: User{
				ID:        1,
				Name:      "John Doe",
				Email:     "",
				Password:  "password123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: ErrEmailRequired,
		},
		{
			name: "missing password",
			user: User{
				ID:        1,
				Name:      "John Doe",
				Email:     "john@example.com",
				Password:  "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: ErrPasswordRequired,
		},
		{
			name: "missing name and email",
			user: User{
				ID:        1,
				Name:      "",
				Email:     "",
				Password:  "password123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: ErrNameRequired, // Should return first error encountered
		},
		{
			name: "missing name and password",
			user: User{
				ID:        1,
				Name:      "",
				Email:     "john@example.com",
				Password:  "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: ErrNameRequired, // Should return first error encountered
		},
		{
			name: "missing email and password",
			user: User{
				ID:        1,
				Name:      "John Doe",
				Email:     "",
				Password:  "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: ErrEmailRequired, // Should return first error encountered
		},
		{
			name: "all fields missing",
			user: User{
				ID:        1,
				Name:      "",
				Email:     "",
				Password:  "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: ErrNameRequired, // Should return first error encountered
		},
		{
			name: "valid user with long name",
			user: User{
				ID:        1,
				Name:      "Very Long Name That Exceeds Normal Length",
				Email:     "longname@example.com",
				Password:  "password123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: nil,
		},
		{
			name: "valid user with complex email",
			user: User{
				ID:        1,
				Name:      "Jane Doe",
				Email:     "jane.doe+test@example.co.uk",
				Password:  "password123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			}
		})
	}
}

