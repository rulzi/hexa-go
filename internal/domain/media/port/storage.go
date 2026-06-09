package port

import (
	"context"
	"io"
)

// StoragePort is the driven port for file storage operations.
type StoragePort interface {
	Save(ctx context.Context, filename string, file io.Reader) (string, error)
	Delete(ctx context.Context, path string) error
	Get(ctx context.Context, path string) (io.ReadCloser, error)
}

// Storage is an alias for StoragePort.
type Storage = StoragePort
