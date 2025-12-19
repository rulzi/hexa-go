package media

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	domainmedia "github.com/rulzi/hexa-go/internal/domain/media"
)

// LocalStorage is the local file system implementation of media.Storage
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new LocalStorage instance
func NewLocalStorage(basePath string) (*LocalStorage, error) {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &LocalStorage{
		basePath: basePath,
	}, nil
}

// Save saves a file and returns the storage path
func (s *LocalStorage) Save(ctx context.Context, filename string, file io.Reader) (string, error) {
	// Generate unique filename with timestamp to avoid conflicts
	timestamp := time.Now().Unix()
	ext := filepath.Ext(filename)
	nameWithoutExt := filename[:len(filename)-len(ext)]
	uniqueFilename := fmt.Sprintf("%s_%d%s", nameWithoutExt, timestamp, ext)

	// Create directory structure: YYYY/MM/DD
	now := time.Now()
	dirPath := filepath.Join(s.basePath, fmt.Sprintf("%d", now.Year()), fmt.Sprintf("%02d", now.Month()), fmt.Sprintf("%02d", now.Day()))

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Full file path
	fullPath := filepath.Join(dirPath, uniqueFilename)

	// Create file
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(fullPath) // Clean up on error
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return relative path from basePath
	relPath, err := filepath.Rel(s.basePath, fullPath)
	if err != nil {
		// If relative path fails, return the path after basePath
		relPath = fullPath[len(s.basePath)+1:]
	}

	return relPath, nil
}

// Delete deletes a file by path
func (s *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// File doesn't exist, but we don't return error (idempotent)
		return nil
	}

	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// Get retrieves a file by path
func (s *LocalStorage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.basePath, path)

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, domainmedia.ErrMediaNotFound
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}
