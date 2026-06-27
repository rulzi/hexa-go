package entity

import "time"

// Article represents the article entity in the domain
type Article struct {
	ID        int64
	Title     string
	Content   string
	AuthorID  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Validate validates the article entity
func (a *Article) Validate() error {
	if a.Title == "" {
		return NewTitleRequired()
	}
	if a.Content == "" {
		return NewContentRequired()
	}
	if a.AuthorID <= 0 {
		return NewAuthorIDRequired()
	}
	return nil
}
