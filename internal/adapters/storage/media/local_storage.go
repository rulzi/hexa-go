package media

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	mediaentity "github.com/rulzi/hexa-go/internal/domain/media/entity"
)

// LocalStorageAdapter stores files on the local filesystem.
type LocalStorageAdapter struct {
	basePath string
}

// NewLocalStorageAdapter creates a local filesystem storage adapter.
func NewLocalStorageAdapter(basePath string) (*LocalStorageAdapter, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &LocalStorageAdapter{basePath: basePath}, nil
}

// Save saves a file and returns the storage path.
func (s *LocalStorageAdapter) Save(ctx context.Context, filename string, file io.Reader) (string, error) {
	relPath, err := generateStoragePath(filename)
	if err != nil {
		return "", fmt.Errorf("invalid filename: %w", err)
	}

	fullPath := filepath.Join(s.basePath, relPath)
	if err := validateResolvedPath(s.basePath, fullPath); err != nil {
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if err := dst.Close(); err != nil {
			log.Printf("Failed to close file: %v", err)
		}
	}()

	if _, err := io.Copy(dst, file); err != nil {
		if removeErr := os.Remove(fullPath); removeErr != nil {
			log.Printf("Failed to remove file: %v", removeErr)
		}
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return relPath, nil
}

// Delete deletes a file by path.
func (s *LocalStorageAdapter) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)
	if err := validateResolvedPath(s.basePath, fullPath); err != nil {
		return err
	}

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil
	}

	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// Get retrieves a file by path.
func (s *LocalStorageAdapter) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.basePath, path)
	if err := validateResolvedPath(s.basePath, fullPath); err != nil {
		return nil, err
	}

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, mediaentity.NewMediaNotFound()
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}
