package article

import "time"

// CreateArticleRequest represents the request DTO for creating an article
type CreateArticleRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	AuthorID int64  `json:"author_id" binding:"required"`
}

// UpdateArticleRequest represents the request DTO for updating an article
type UpdateArticleRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

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
