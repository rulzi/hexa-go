package media

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMedia_Validate(t *testing.T) {
	tests := []struct {
		name         string
		media        Media
		wantErrCheck func(error) bool
	}{
		{
			name: "valid media",
			media: Media{
				ID:        1,
				Name:      "test-image.jpg",
				Path:      "/storage/2025/12/19/test-image.jpg",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErrCheck: nil,
		},
		{
			name: "missing name",
			media: Media{
				ID:        1,
				Name:      "",
				Path:      "/storage/2025/12/19/test-image.jpg",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErrCheck: IsNameRequired,
		},
		{
			name: "missing path",
			media: Media{
				ID:        1,
				Name:      "test-image.jpg",
				Path:      "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErrCheck: IsPathRequired,
		},
		{
			name: "missing name and path",
			media: Media{
				ID:        1,
				Name:      "",
				Path:      "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErrCheck: IsNameRequired, // Should return first error encountered
		},
		{
			name: "valid media with long name",
			media: Media{
				ID:        1,
				Name:      "very-long-filename-with-many-characters-and-extension.png",
				Path:      "/storage/2025/12/19/very-long-filename-with-many-characters-and-extension.png",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErrCheck: nil,
		},
		{
			name: "valid media with nested path",
			media: Media{
				ID:        1,
				Name:      "document.pdf",
				Path:      "/storage/2025/12/19/subfolder/document.pdf",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErrCheck: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.media.Validate()
			if tt.wantErrCheck == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.True(t, tt.wantErrCheck(err))
			}
		})
	}
}
