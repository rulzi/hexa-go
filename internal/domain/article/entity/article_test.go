package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestArticle_Validate(t *testing.T) {
	tests := []struct {
		name         string
		article      Article
		wantErrCheck func(error) bool
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
			wantErrCheck: nil,
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
			wantErrCheck: IsTitleRequired,
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
			wantErrCheck: IsContentRequired,
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
			wantErrCheck: IsAuthorIDRequired,
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
			wantErrCheck: IsAuthorIDRequired,
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
			wantErrCheck: IsTitleRequired,
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
			wantErrCheck: IsTitleRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.article.Validate()
			if tt.wantErrCheck == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.True(t, tt.wantErrCheck(err))
			}
		})
	}
}
