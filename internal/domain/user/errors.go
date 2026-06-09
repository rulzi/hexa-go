package user

import domainerrs "github.com/rulzi/hexa-go/internal/domain/errs"

// Domain-specific error codes. Used for assertions and logging — not HTTP status.
const (
	CodeUserNotFound        = "USER_NOT_FOUND"
	CodeEmailExists         = "USER_EMAIL_EXISTS"
	CodeNameRequired        = "USER_NAME_REQUIRED"
	CodeEmailRequired       = "USER_EMAIL_REQUIRED"
	CodePasswordRequired    = "USER_PASSWORD_REQUIRED"
	CodePasswordTooShort    = "USER_PASSWORD_TOO_SHORT"
	CodeInvalidEmail        = "USER_INVALID_EMAIL"
	CodeInvalidCredentials  = "USER_INVALID_CREDENTIALS"
)

func NewUserNotFound() error {
	return domainerrs.NewNotFound(CodeUserNotFound, "user not found")
}

func NewEmailExists() error {
	return domainerrs.NewConflict(CodeEmailExists, "email already exists")
}

func NewNameRequired() error {
	return domainerrs.NewValidation(CodeNameRequired, "name is required")
}

func NewEmailRequired() error {
	return domainerrs.NewValidation(CodeEmailRequired, "email is required")
}

func NewPasswordRequired() error {
	return domainerrs.NewValidation(CodePasswordRequired, "password is required")
}

func NewPasswordTooShort() error {
	return domainerrs.NewValidation(CodePasswordTooShort, "password must be at least 6 characters")
}

func NewInvalidEmail() error {
	return domainerrs.NewValidation(CodeInvalidEmail, "invalid email format")
}

func NewInvalidCredentials() error {
	return domainerrs.NewUnauthorized(CodeInvalidCredentials, "invalid email or password")
}

func IsUserNotFound(err error) bool       { return domainerrs.IsCode(err, CodeUserNotFound) }
func IsEmailExists(err error) bool        { return domainerrs.IsCode(err, CodeEmailExists) }
func IsNameRequired(err error) bool       { return domainerrs.IsCode(err, CodeNameRequired) }
func IsEmailRequired(err error) bool      { return domainerrs.IsCode(err, CodeEmailRequired) }
func IsPasswordRequired(err error) bool   { return domainerrs.IsCode(err, CodePasswordRequired) }
func IsPasswordTooShort(err error) bool   { return domainerrs.IsCode(err, CodePasswordTooShort) }
func IsInvalidCredentials(err error) bool { return domainerrs.IsCode(err, CodeInvalidCredentials) }
