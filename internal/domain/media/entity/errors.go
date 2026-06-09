package entity

import domainerrs "github.com/rulzi/hexa-go/internal/domain/errs"

const (
	CodeMediaNotFound = "MEDIA_NOT_FOUND"
	CodeNameRequired  = "MEDIA_NAME_REQUIRED"
	CodePathRequired  = "MEDIA_PATH_REQUIRED"
)

func NewMediaNotFound() error {
	return domainerrs.NewNotFound(CodeMediaNotFound, "media not found")
}

func NewNameRequired() error {
	return domainerrs.NewValidation(CodeNameRequired, "name is required")
}

func NewPathRequired() error {
	return domainerrs.NewValidation(CodePathRequired, "path is required")
}

func IsMediaNotFound(err error) bool { return domainerrs.IsCode(err, CodeMediaNotFound) }
func IsNameRequired(err error) bool  { return domainerrs.IsCode(err, CodeNameRequired) }
func IsPathRequired(err error) bool  { return domainerrs.IsCode(err, CodePathRequired) }
