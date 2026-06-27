package usecase

import (
	"context"
	"testing"
	"time"

	articleentity "github.com/rulzi/hexa-go/internal/domain/article/entity"
	"github.com/stretchr/testify/assert"
)

func TestDeleteArticle_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewDeleteArticle(repo, cache, listCache)

	articleID := int64(1)
	existingArticle := &articleentity.Article{
		ID: articleID, Title: "Test Article", Content: "Test Content", AuthorID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	repo.On("GetByID", ctx, articleID).Return(existingArticle, nil)
	repo.On("Delete", ctx, articleID).Return(nil)
	cache.On("Delete", ctx, articleID).Return(nil)
	cache.On("InvalidateList", ctx).Return(nil)
	listCache.On("InvalidateArticleList", ctx).Return(nil)

	err := uc.Execute(ctx, articleID)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
	listCache.AssertExpectations(t)
}

func TestDeleteArticle_Execute_ArticleNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewDeleteArticle(repo, cache, listCache)

	articleID := int64(1)
	repo.On("GetByID", ctx, articleID).Return(nil, nil)

	err := uc.Execute(ctx, articleID)

	assert.Error(t, err)
	assert.True(t, articleentity.IsArticleNotFound(err))
	repo.AssertNotCalled(t, "Delete")
}

func TestDeleteArticle_Execute_WithNilCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	listCache := &mockArticleListCache{}

	uc := NewDeleteArticle(repo, nil, listCache)

	articleID := int64(1)
	existingArticle := &articleentity.Article{
		ID: articleID, Title: "Test Article", Content: "Test Content", AuthorID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	repo.On("GetByID", ctx, articleID).Return(existingArticle, nil)
	repo.On("Delete", ctx, articleID).Return(nil)
	listCache.On("InvalidateArticleList", ctx).Return(nil)

	err := uc.Execute(ctx, articleID)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	listCache.AssertExpectations(t)
}
