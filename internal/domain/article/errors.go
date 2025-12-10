package article

import "errors"

var (
	ErrArticleNotFound  = errors.New("article not found")
	ErrTitleRequired    = errors.New("title is required")
	ErrContentRequired  = errors.New("content is required")
	ErrAuthorIDRequired = errors.New("author id is required")
)

