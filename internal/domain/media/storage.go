package media

import (
	"context"
	"io"
)

// Storage is the interface for file storage operations
type Storage interface {
	// Save saves a file and returns the storage path
	Save(ctx context.Context, filename string, file io.Reader) (string, error)

	// Delete deletes a file by path
	Delete(ctx context.Context, path string) error

	// Get retrieves a file by path
	Get(ctx context.Context, path string) (io.ReadCloser, error)
}
