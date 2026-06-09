package media

import (
	"fmt"
	"path/filepath"
	"strings"
)

func validateResolvedPath(basePath, fullPath string) error {
	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return fmt.Errorf("invalid base path: %w", err)
	}

	absFull, err := filepath.Abs(fullPath)
	if err != nil {
		return fmt.Errorf("invalid storage path: %w", err)
	}

	rel, err := filepath.Rel(absBase, absFull)
	if err != nil {
		return fmt.Errorf("invalid storage path: %w", err)
	}

	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("path traversal detected")
	}

	return nil
}
