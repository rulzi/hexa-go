package user

import (
	"fmt"
	"testing"

	domainerrs "github.com/rulzi/hexa-go/internal/domain/errs"
	"github.com/stretchr/testify/assert"
)

func TestUserErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		check    func(error) bool
		kind     func(error) bool
	}{
		{"not found", NewUserNotFound(), IsUserNotFound, domainerrs.IsNotFound},
		{"email exists", NewEmailExists(), IsEmailExists, domainerrs.IsConflict},
		{"name required", NewNameRequired(), IsNameRequired, domainerrs.IsValidation},
		{"invalid credentials", NewInvalidCredentials(), IsInvalidCredentials, domainerrs.IsUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, tt.check(tt.err))
			assert.True(t, tt.kind(tt.err))
		})
	}
}

func TestUserErrors_Wrapped(t *testing.T) {
	wrapped := fmt.Errorf("usecase: %w", NewUserNotFound())
	assert.True(t, IsUserNotFound(wrapped))
	assert.True(t, domainerrs.IsNotFound(wrapped))
}
