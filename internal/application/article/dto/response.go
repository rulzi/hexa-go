package dto

import "time"

// ArticleResponse represents the response DTO for article
type ArticleResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  int64     `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListArticlesResponse represents the response DTO for listing articles
type ListArticlesResponse struct {
	Articles []ArticleResponse `json:"articles"`
	Total    int64             `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}
