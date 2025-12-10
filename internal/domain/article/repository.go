package article

import "context"

// Repository is the driven port (interface) for article persistence
// This defines what the domain needs, not how it's implemented
type Repository interface {
	// Create creates a new article
	Create(ctx context.Context, article *Article) (*Article, error)

	// GetByID retrieves an article by ID
	GetByID(ctx context.Context, id int64) (*Article, error)

	// Update updates an existing article
	Update(ctx context.Context, article *Article) (*Article, error)

	// Delete deletes an article by ID
	Delete(ctx context.Context, id int64) error

	// List retrieves all articles with pagination
	List(ctx context.Context, limit, offset int) ([]*Article, error)

	// ListByAuthor retrieves articles by author ID with pagination
	ListByAuthor(ctx context.Context, authorID int64, limit, offset int) ([]*Article, error)

	// Count returns the total number of articles
	Count(ctx context.Context) (int64, error)

	// CountByAuthor returns the total number of articles by author
	CountByAuthor(ctx context.Context, authorID int64) (int64, error)
}

