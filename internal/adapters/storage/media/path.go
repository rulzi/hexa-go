package media

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// generateStoragePath builds a date-partitioned relative path with a timestamped filename.
func generateStoragePath(filename string) (string, error) {
	safeName, err := SanitizeFilename(filename)
	if err != nil {
		return "", err
	}

	timestamp := time.Now().Unix()
	ext := filepath.Ext(safeName)
	nameWithoutExt := strings.TrimSuffix(safeName, ext)

	now := time.Now()
	return fmt.Sprintf("%d/%02d/%02d/%s_%d%s",
		now.Year(), now.Month(), now.Day(),
		nameWithoutExt, timestamp, ext,
	), nil
}
