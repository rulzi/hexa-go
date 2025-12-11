package user

import "errors"

var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
	// ErrEmailExists is returned when email already exists in the system
	ErrEmailExists = errors.New("email already exists")
	// ErrNameRequired is returned when user name is missing
	ErrNameRequired = errors.New("name is required")
	// ErrEmailRequired is returned when user email is missing
	ErrEmailRequired = errors.New("email is required")
	// ErrPasswordRequired is returned when user password is missing
	ErrPasswordRequired = errors.New("password is required")
	// ErrInvalidEmail is returned when email format is invalid
	ErrInvalidEmail = errors.New("invalid email format")
	// ErrInvalidCredentials is returned when login credentials are invalid
	ErrInvalidCredentials = errors.New("invalid email or password")
)
