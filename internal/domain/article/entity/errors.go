package entity

import domainerrs "github.com/rulzi/hexa-go/internal/domain/errs"

// Domain-specific error codes. Used for assertions and logging — not HTTP status.
const (
	CodeArticleNotFound  = "ARTICLE_NOT_FOUND"
	CodeTitleRequired    = "ARTICLE_TITLE_REQUIRED"
	CodeContentRequired  = "ARTICLE_CONTENT_REQUIRED"
	CodeAuthorIDRequired = "ARTICLE_AUTHOR_ID_REQUIRED"
)

// NewArticleNotFound returns an error when an article cannot be found.
func NewArticleNotFound() error {
	return domainerrs.NewNotFound(CodeArticleNotFound, "article not found")
}

// NewTitleRequired returns an error when article title is missing.
func NewTitleRequired() error {
	return domainerrs.NewValidation(CodeTitleRequired, "title is required")
}

// NewContentRequired returns an error when article content is missing.
func NewContentRequired() error {
	return domainerrs.NewValidation(CodeContentRequired, "content is required")
}

// NewAuthorIDRequired returns an error when article author ID is missing.
func NewAuthorIDRequired() error {
	return domainerrs.NewValidation(CodeAuthorIDRequired, "author id is required")
}

// IsArticleNotFound reports whether err is an article-not-found error.
func IsArticleNotFound(err error) bool { return domainerrs.IsCode(err, CodeArticleNotFound) }

// IsTitleRequired reports whether err is a missing-title validation error.
func IsTitleRequired(err error) bool { return domainerrs.IsCode(err, CodeTitleRequired) }

// IsContentRequired reports whether err is a missing-content validation error.
func IsContentRequired(err error) bool { return domainerrs.IsCode(err, CodeContentRequired) }

// IsAuthorIDRequired reports whether err is a missing-author-id validation error.
func IsAuthorIDRequired(err error) bool { return domainerrs.IsCode(err, CodeAuthorIDRequired) }
