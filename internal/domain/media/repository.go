package media

import "context"

// Repository is the driven port (interface) for media persistence
// This defines what the domain needs, not how it's implemented
type Repository interface {
	// Create creates a new media
	Create(ctx context.Context, media *Media) (*Media, error)

	// GetByID retrieves a media by ID
	GetByID(ctx context.Context, id int64) (*Media, error)

	// Update updates an existing media
	Update(ctx context.Context, media *Media) (*Media, error)

	// Delete deletes a media by ID
	Delete(ctx context.Context, id int64) error

	// List retrieves all media with pagination
	List(ctx context.Context, limit, offset int) ([]*Media, error)

	// Count returns the total number of media
	Count(ctx context.Context) (int64, error)
}
