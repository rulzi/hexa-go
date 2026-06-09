package errs

import (
	"errors"
	"fmt"
)

// DomainError is the common contract for typed domain errors.
// HTTP adapters map these types to status codes without the domain knowing HTTP.
type DomainError interface {
	error
	Code() string
}

// NotFoundError indicates a requested resource does not exist.
type NotFoundError struct {
	code    string
	message string
	cause   error
}

func NewNotFound(code, message string) *NotFoundError {
	return &NotFoundError{code: code, message: message}
}

func (e *NotFoundError) Error() string { return e.message }
func (e *NotFoundError) Code() string  { return e.code }
func (e *NotFoundError) Unwrap() error { return e.cause }

// ValidationError indicates invalid input or business rule violation.
type ValidationError struct {
	code    string
	message string
	field   string
	cause   error
}

func NewValidation(code, message string) *ValidationError {
	return &ValidationError{code: code, message: message}
}

func NewValidationField(code, message, field string) *ValidationError {
	return &ValidationError{code: code, message: message, field: field}
}

func (e *ValidationError) Error() string { return e.message }
func (e *ValidationError) Code() string  { return e.code }
func (e *ValidationError) Field() string { return e.field }
func (e *ValidationError) Unwrap() error { return e.cause }

// ConflictError indicates a state conflict (e.g. duplicate unique value).
type ConflictError struct {
	code    string
	message string
	cause   error
}

func NewConflict(code, message string) *ConflictError {
	return &ConflictError{code: code, message: message}
}

func (e *ConflictError) Error() string { return e.message }
func (e *ConflictError) Code() string  { return e.code }
func (e *ConflictError) Unwrap() error { return e.cause }

// UnauthorizedError indicates authentication or authorization failure.
type UnauthorizedError struct {
	code    string
	message string
	cause   error
}

func NewUnauthorized(code, message string) *UnauthorizedError {
	return &UnauthorizedError{code: code, message: message}
}

func (e *UnauthorizedError) Error() string { return e.message }
func (e *UnauthorizedError) Code() string  { return e.code }
func (e *UnauthorizedError) Unwrap() error { return e.cause }

// IsNotFound reports whether err is a NotFoundError (including wrapped).
func IsNotFound(err error) bool {
	var e *NotFoundError
	return errors.As(err, &e)
}

// IsValidation reports whether err is a ValidationError (including wrapped).
func IsValidation(err error) bool {
	var e *ValidationError
	return errors.As(err, &e)
}

// IsConflict reports whether err is a ConflictError (including wrapped).
func IsConflict(err error) bool {
	var e *ConflictError
	return errors.As(err, &e)
}

// IsUnauthorized reports whether err is an UnauthorizedError (including wrapped).
func IsUnauthorized(err error) bool {
	var e *UnauthorizedError
	return errors.As(err, &e)
}

// IsCode reports whether err is a DomainError with the given code.
func IsCode(err error, code string) bool {
	var de DomainError
	if errors.As(err, &de) {
		return de.Code() == code
	}
	return false
}

// Message returns the client-safe message from a domain error, or a generic fallback.
func Message(err error) string {
	var de DomainError
	if errors.As(err, &de) {
		return de.Error()
	}
	return "internal server error"
}

// Wrap adds context to a domain error while preserving its type for errors.As.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case *NotFoundError:
		return &NotFoundError{code: e.code, message: fmt.Sprintf("%s: %s", message, e.message), cause: err}
	case *ValidationError:
		return &ValidationError{code: e.code, message: fmt.Sprintf("%s: %s", message, e.message), field: e.field, cause: err}
	case *ConflictError:
		return &ConflictError{code: e.code, message: fmt.Sprintf("%s: %s", message, e.message), cause: err}
	case *UnauthorizedError:
		return &UnauthorizedError{code: e.code, message: fmt.Sprintf("%s: %s", message, e.message), cause: err}
	default:
		return fmt.Errorf("%s: %w", message, err)
	}
}
