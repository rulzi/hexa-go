package media

import (
	"context"
	"io"
)

// StoragePort is the driven port for file storage operations.
// Implementations may use local filesystem, S3, or any other backend.
type StoragePort interface {
	// Save stores a file and returns a relative storage path (object key).
	Save(ctx context.Context, filename string, file io.Reader) (string, error)

	// Delete removes a file by path. Idempotent when the file does not exist.
	Delete(ctx context.Context, path string) error

	// Get retrieves a file by path.
	Get(ctx context.Context, path string) (io.ReadCloser, error)
}

// Storage is an alias for StoragePort for backward compatibility.
type Storage = StoragePort
