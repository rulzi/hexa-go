package entity

import domainerrs "github.com/rulzi/hexa-go/internal/domain/errs"

// Domain-specific error codes. Used for assertions and logging — not HTTP status.
const (
	CodeMediaNotFound = "MEDIA_NOT_FOUND"
	CodeNameRequired  = "MEDIA_NAME_REQUIRED"
	CodePathRequired  = "MEDIA_PATH_REQUIRED"
)

// NewMediaNotFound returns an error when media cannot be found.
func NewMediaNotFound() error {
	return domainerrs.NewNotFound(CodeMediaNotFound, "media not found")
}

// NewNameRequired returns an error when media name is missing.
func NewNameRequired() error {
	return domainerrs.NewValidation(CodeNameRequired, "name is required")
}

// NewPathRequired returns an error when media path is missing.
func NewPathRequired() error {
	return domainerrs.NewValidation(CodePathRequired, "path is required")
}

// IsMediaNotFound reports whether err is a media-not-found error.
func IsMediaNotFound(err error) bool { return domainerrs.IsCode(err, CodeMediaNotFound) }

// IsNameRequired reports whether err is a missing-name validation error.
func IsNameRequired(err error) bool { return domainerrs.IsCode(err, CodeNameRequired) }

// IsPathRequired reports whether err is a missing-path validation error.
func IsPathRequired(err error) bool { return domainerrs.IsCode(err, CodePathRequired) }
