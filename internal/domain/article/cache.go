package article

import "context"

// Cache is a port for caching article data
// This is an infrastructure concern but defined as a port in domain
// to maintain dependency inversion
type Cache interface {
	// Get retrieves a cached article by ID
	Get(ctx context.Context, id int64) (*Article, error)

	// Set stores an article in cache
	Set(ctx context.Context, id int64, article *Article) error

	// Delete removes an article from cache
	Delete(ctx context.Context, id int64) error

	// InvalidateList invalidates all article list caches
	InvalidateList(ctx context.Context) error
}

