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

func TestWrap_Nil(t *testing.T) {
	assert.Nil(t, Wrap(nil, "context"))
}

func TestWrap_AllDomainTypes(t *testing.T) {
	testCases := []struct {
		name   string
		err    error
		check  func(error) bool
	}{
		{"not found", NewNotFound("NF", "not found"), IsNotFound},
		{"validation", NewValidation("VAL", "invalid"), IsValidation},
		{"conflict", NewConflict("CF", "conflict"), IsConflict},
		{"unauthorized", NewUnauthorized("UA", "unauthorized"), IsUnauthorized},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wrapped := Wrap(tc.err, "operation")
			assert.True(t, tc.check(wrapped))
			assert.ErrorIs(t, wrapped, tc.err)
		})
	}
}

func TestWrap_GenericError(t *testing.T) {
	original := errors.New("db down")
	wrapped := Wrap(original, "query")

	assert.False(t, IsNotFound(wrapped))
	assert.ErrorIs(t, wrapped, original)
}

func TestIsCode_False(t *testing.T) {
	err := NewNotFound("USER_NOT_FOUND", "missing")
	assert.False(t, IsCode(err, "OTHER_CODE"))
	assert.False(t, IsCode(errors.New("plain"), "USER_NOT_FOUND"))
}

func TestDomainError_Unwrap(t *testing.T) {
	cause := errors.New("root cause")

	notFound := &NotFoundError{code: "NF", message: "missing", cause: cause}
	assert.Equal(t, cause, notFound.Unwrap())

	validation := &ValidationError{code: "VAL", message: "bad", cause: cause}
	assert.Equal(t, cause, validation.Unwrap())

	conflict := &ConflictError{code: "CF", message: "dup", cause: cause}
	assert.Equal(t, cause, conflict.Unwrap())

	unauthorized := &UnauthorizedError{code: "UA", message: "denied", cause: cause}
	assert.Equal(t, cause, unauthorized.Unwrap())
}

func TestDomainError_CodeAndError(t *testing.T) {
	conflict := NewConflict("CF", "duplicate")
	assert.Equal(t, "duplicate", conflict.Error())
	assert.Equal(t, "CF", conflict.Code())

	unauthorized := NewUnauthorized("UA", "denied")
	assert.Equal(t, "denied", unauthorized.Error())
	assert.Equal(t, "UA", unauthorized.Code())
}
