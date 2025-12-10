package article

import "time"

// Article represents the article entity in the domain
type Article struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  int64     `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate validates the article entity
func (a *Article) Validate() error {
	if a.Title == "" {
		return ErrTitleRequired
	}
	if a.Content == "" {
		return ErrContentRequired
	}
	if a.AuthorID <= 0 {
		return ErrAuthorIDRequired
	}
	return nil
}

