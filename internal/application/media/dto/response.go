package dto

import (
	"fmt"
	"strings"
	"time"
)

// MediaResponse represents the response DTO for media
type MediaResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"` // Storage path (relative)
	URL       string    `json:"url"`  // Full URL to access the file
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BuildURL builds the full URL from base URL and path
func BuildURL(baseURL, path string) string {
	// Remove trailing slash from baseURL
	baseURL = strings.TrimSuffix(baseURL, "/")

	// Ensure path starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Normalize path separators for URL
	path = strings.ReplaceAll(path, "\\", "/")

	return fmt.Sprintf("%s/api/v1/media/files%s", baseURL, path)
}

// ListMediaResponse represents the response DTO for listing media
type ListMediaResponse struct {
	Media  []MediaResponse `json:"media"`
	Total  int64           `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}
