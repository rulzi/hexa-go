package article

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestArticle_Validate(t *testing.T) {
	tests := []struct {
		name    string
		article Article
		wantErr error
	}{
		{
			name: "valid article",
			article: Article{
				ID:        1,
				Title:     "Test Article",
				Content:   "This is a test article content",
				AuthorID:  1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: nil,
		},
		{
			name: "missing title",
			article: Article{
				ID:        1,
				Title:     "",
				Content:   "This is a test article content",
				AuthorID:  1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: ErrTitleRequired,
		},
		{
			name: "missing content",
			article: Article{
				ID:        1,
				Title:     "Test Article",
				Content:   "",
				AuthorID:  1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: ErrContentRequired,
		},
		{
			name: "zero author ID",
			article: Article{
				ID:        1,
				Title:     "Test Article",
				Content:   "This is a test article content",
				AuthorID:  0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: ErrAuthorIDRequired,
		},
		{
			name: "negative author ID",
			article: Article{
				ID:        1,
				Title:     "Test Article",
				Content:   "This is a test article content",
				AuthorID:  -1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: ErrAuthorIDRequired,
		},
		{
			name: "missing title and content",
			article: Article{
				ID:        1,
				Title:     "",
				Content:   "",
				AuthorID:  1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: ErrTitleRequired, // Should return first error encountered
		},
		{
			name: "all fields invalid",
			article: Article{
				ID:        1,
				Title:     "",
				Content:   "",
				AuthorID:  0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: ErrTitleRequired, // Should return first error encountered
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.article.Validate()
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			}
		})
	}
}
