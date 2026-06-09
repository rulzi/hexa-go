package article

import domainerrs "github.com/rulzi/hexa-go/internal/domain/errs"

const (
	CodeArticleNotFound   = "ARTICLE_NOT_FOUND"
	CodeTitleRequired     = "ARTICLE_TITLE_REQUIRED"
	CodeContentRequired   = "ARTICLE_CONTENT_REQUIRED"
	CodeAuthorIDRequired  = "ARTICLE_AUTHOR_ID_REQUIRED"
)

func NewArticleNotFound() error {
	return domainerrs.NewNotFound(CodeArticleNotFound, "article not found")
}

func NewTitleRequired() error {
	return domainerrs.NewValidation(CodeTitleRequired, "title is required")
}

func NewContentRequired() error {
	return domainerrs.NewValidation(CodeContentRequired, "content is required")
}

func NewAuthorIDRequired() error {
	return domainerrs.NewValidation(CodeAuthorIDRequired, "author id is required")
}

func IsArticleNotFound(err error) bool  { return domainerrs.IsCode(err, CodeArticleNotFound) }
func IsTitleRequired(err error) bool    { return domainerrs.IsCode(err, CodeTitleRequired) }
func IsContentRequired(err error) bool  { return domainerrs.IsCode(err, CodeContentRequired) }
func IsAuthorIDRequired(err error) bool { return domainerrs.IsCode(err, CodeAuthorIDRequired) }
