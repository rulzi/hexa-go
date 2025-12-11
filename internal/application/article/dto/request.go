package dto

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
