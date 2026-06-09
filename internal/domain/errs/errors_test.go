package errs

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotFoundError(t *testing.T) {
	err := NewNotFound("USER_NOT_FOUND", "user not found")

	assert.Equal(t, "user not found", err.Error())
	assert.Equal(t, "USER_NOT_FOUND", err.Code())
	assert.True(t, IsNotFound(err))
	assert.True(t, IsCode(err, "USER_NOT_FOUND"))
	assert.False(t, IsValidation(err))
}

func TestValidationError(t *testing.T) {
	err := NewValidationField("USER_NAME_REQUIRED", "name is required", "name")

	assert.Equal(t, "name is required", err.Error())
	assert.Equal(t, "USER_NAME_REQUIRED", err.Code())
	assert.Equal(t, "name", err.Field())
	assert.True(t, IsValidation(err))
}

func TestConflictError(t *testing.T) {
	err := NewConflict("USER_EMAIL_EXISTS", "email already exists")

	assert.True(t, IsConflict(err))
	assert.Equal(t, "email already exists", Message(err))
}

func TestUnauthorizedError(t *testing.T) {
	err := NewUnauthorized("USER_INVALID_CREDENTIALS", "invalid email or password")

	assert.True(t, IsUnauthorized(err))
}

func TestMessage_Fallback(t *testing.T) {
	assert.Equal(t, "internal server error", Message(errors.New("db down")))
}

func TestWrap_PreservesType(t *testing.T) {
	original := NewNotFound("USER_NOT_FOUND", "user not found")
	wrapped := Wrap(original, "get user")

	assert.True(t, IsNotFound(wrapped))
	assert.True(t, IsCode(wrapped, "USER_NOT_FOUND"))
	assert.ErrorIs(t, wrapped, original)
}

func TestErrorsAs_WithFmtWrap(t *testing.T) {
	original := NewValidation("USER_EMAIL_REQUIRED", "email is required")
	wrapped := fmt.Errorf("validate user: %w", original)

	assert.True(t, IsValidation(wrapped))
	assert.Equal(t, "email is required", Message(wrapped))
}
