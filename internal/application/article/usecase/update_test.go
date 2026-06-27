package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/rulzi/hexa-go/internal/application/article/dto"
	articleentity "github.com/rulzi/hexa-go/internal/domain/article/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateArticle_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, cache, listCache)

	articleID := int64(1)
	existingArticle := &articleentity.Article{
		ID: articleID, Title: "Old Title", Content: "Old Content", AuthorID: 1,
		CreatedAt: time.Now().Add(-24 * time.Hour), UpdatedAt: time.Now().Add(-24 * time.Hour),
	}
	req := dto.UpdateArticleRequest{Title: "New Title", Content: "New Content"}
	updatedArticle := &articleentity.Article{
		ID: articleID, Title: req.Title, Content: req.Content, AuthorID: existingArticle.AuthorID,
		CreatedAt: existingArticle.CreatedAt, UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, articleID).Return(existingArticle, nil)
	repo.On("Update", ctx, mock.Anything).Return(updatedArticle, nil)
	cache.On("Delete", ctx, articleID).Return(nil)
	cache.On("InvalidateList", ctx).Return(nil)
	listCache.On("InvalidateArticleList", ctx).Return(nil)

	result, err := uc.Update.Execute(ctx, articleID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Title, result.Title)
	assert.Equal(t, req.Content, result.Content)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
	listCache.AssertExpectations(t)
}

func TestUpdateArticle_Execute_ArticleNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, cache, listCache)

	articleID := int64(1)
	req := dto.UpdateArticleRequest{Title: "New Title", Content: "New Content"}
	repo.On("GetByID", ctx, articleID).Return(nil, nil)

	result, err := uc.Update.Execute(ctx, articleID, req)

	assert.Error(t, err)
	assert.True(t, articleentity.IsArticleNotFound(err))
	assert.Nil(t, result)
	repo.AssertNotCalled(t, "Update")
}

func TestUpdateArticle_Execute_ValidationError(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, cache, listCache)

	articleID := int64(1)
	existingArticle := &articleentity.Article{
		ID: articleID, Title: "Old Title", Content: "Old Content", AuthorID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	tests := []struct {
		name string
		req  dto.UpdateArticleRequest
	}{
		{"empty title", dto.UpdateArticleRequest{Title: "", Content: "New Content"}},
		{"empty content", dto.UpdateArticleRequest{Title: "New Title", Content: ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.On("GetByID", ctx, articleID).Return(existingArticle, nil)
			result, err := uc.Update.Execute(ctx, articleID, tt.req)
			assert.Error(t, err)
			assert.Nil(t, result)
			repo.AssertNotCalled(t, "Update")
		})
	}
}
