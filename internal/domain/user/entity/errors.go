package entity

import domainerrs "github.com/rulzi/hexa-go/internal/domain/errs"

// Domain-specific error codes. Used for assertions and logging — not HTTP status.
const (
	CodeUserNotFound       = "USER_NOT_FOUND"
	CodeEmailExists        = "USER_EMAIL_EXISTS"
	CodeNameRequired       = "USER_NAME_REQUIRED"
	CodeEmailRequired      = "USER_EMAIL_REQUIRED"
	CodePasswordRequired   = "USER_PASSWORD_REQUIRED"
	CodePasswordTooShort   = "USER_PASSWORD_TOO_SHORT"
	CodeInvalidEmail       = "USER_INVALID_EMAIL"
	CodeInvalidCredentials = "USER_INVALID_CREDENTIALS"
)

// NewUserNotFound returns an error when a user cannot be found.
func NewUserNotFound() error {
	return domainerrs.NewNotFound(CodeUserNotFound, "user not found")
}

// NewEmailExists returns an error when registering with a duplicate email.
func NewEmailExists() error {
	return domainerrs.NewConflict(CodeEmailExists, "email already exists")
}

// NewNameRequired returns an error when user name is missing.
func NewNameRequired() error {
	return domainerrs.NewValidation(CodeNameRequired, "name is required")
}

// NewEmailRequired returns an error when user email is missing.
func NewEmailRequired() error {
	return domainerrs.NewValidation(CodeEmailRequired, "email is required")
}

// NewPasswordRequired returns an error when user password is missing.
func NewPasswordRequired() error {
	return domainerrs.NewValidation(CodePasswordRequired, "password is required")
}

// NewPasswordTooShort returns an error when password does not meet length requirements.
func NewPasswordTooShort() error {
	return domainerrs.NewValidation(CodePasswordTooShort, "password must be at least 6 characters")
}

// NewInvalidEmail returns an error when email format is invalid.
func NewInvalidEmail() error {
	return domainerrs.NewValidation(CodeInvalidEmail, "invalid email format")
}

// NewInvalidCredentials returns an error when login credentials are invalid.
func NewInvalidCredentials() error {
	return domainerrs.NewUnauthorized(CodeInvalidCredentials, "invalid email or password")
}

// IsUserNotFound reports whether err is a user-not-found error.
func IsUserNotFound(err error) bool { return domainerrs.IsCode(err, CodeUserNotFound) }

// IsEmailExists reports whether err is a duplicate-email conflict error.
func IsEmailExists(err error) bool { return domainerrs.IsCode(err, CodeEmailExists) }

// IsNameRequired reports whether err is a missing-name validation error.
func IsNameRequired(err error) bool { return domainerrs.IsCode(err, CodeNameRequired) }

// IsEmailRequired reports whether err is a missing-email validation error.
func IsEmailRequired(err error) bool { return domainerrs.IsCode(err, CodeEmailRequired) }

// IsPasswordRequired reports whether err is a missing-password validation error.
func IsPasswordRequired(err error) bool { return domainerrs.IsCode(err, CodePasswordRequired) }

// IsPasswordTooShort reports whether err is a password-too-short validation error.
func IsPasswordTooShort(err error) bool { return domainerrs.IsCode(err, CodePasswordTooShort) }

// IsInvalidCredentials reports whether err is an invalid-credentials error.
func IsInvalidCredentials(err error) bool { return domainerrs.IsCode(err, CodeInvalidCredentials) }
