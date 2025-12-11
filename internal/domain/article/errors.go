package article

import "errors"

var (
	// ErrArticleNotFound is returned when an article is not found
	ErrArticleNotFound = errors.New("article not found")
	// ErrTitleRequired is returned when article title is missing
	ErrTitleRequired = errors.New("title is required")
	// ErrContentRequired is returned when article content is missing
	ErrContentRequired = errors.New("content is required")
	// ErrAuthorIDRequired is returned when author ID is missing
	ErrAuthorIDRequired = errors.New("author id is required")
)
