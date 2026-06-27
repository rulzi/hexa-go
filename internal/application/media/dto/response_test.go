package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildURL(t *testing.T) {
	testCases := []struct {
		name     string
		baseURL  string
		path     string
		expected string
	}{
		{
			name:     "trailing slash on base",
			baseURL:  "http://localhost:8080/",
			path:     "2025/01/01/file.jpg",
			expected: "http://localhost:8080/api/v1/media/files/2025/01/01/file.jpg",
		},
		{
			name:     "path with leading slash",
			baseURL:  "http://localhost:8080",
			path:     "/2025/01/01/file.jpg",
			expected: "http://localhost:8080/api/v1/media/files/2025/01/01/file.jpg",
		},
		{
			name:     "windows path separators",
			baseURL:  "http://localhost:8080",
			path:     `2025\01\01\file.jpg`,
			expected: "http://localhost:8080/api/v1/media/files/2025/01/01/file.jpg",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, BuildURL(tc.baseURL, tc.path))
		})
	}
}
