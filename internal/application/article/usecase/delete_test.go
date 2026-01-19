package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
	"github.com/stretchr/testify/assert"
)

func TestNewDeleteArticleUseCase(t *testing.T) {
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewDeleteArticleUseCase(repo, cache, listCache)

	assert.NotNil(t, uc)
	assert.Equal(t, repo, uc.articleRepo)
	assert.Equal(t, cache, uc.cache)
	assert.Equal(t, listCache, uc.listCache)
}

func TestDeleteArticleUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewDeleteArticleUseCase(repo, cache, listCache)

	articleID := int64(1)
	existingArticle := &domainarticle.Article{
		ID:        articleID,
		Title:     "Test Article",
		Content:   "Test Content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

func TestDeleteArticleUseCase_Execute_ArticleNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewDeleteArticleUseCase(repo, cache, listCache)

	articleID := int64(1)

	repo.On("GetByID", ctx, articleID).Return(nil, nil)

	err := uc.Execute(ctx, articleID)

	assert.Error(t, err)
	assert.Equal(t, domainarticle.ErrArticleNotFound, err)
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Delete")
	cache.AssertNotCalled(t, "Delete")
}

func TestDeleteArticleUseCase_Execute_GetByIDError(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewDeleteArticleUseCase(repo, cache, listCache)

	articleID := int64(1)
	repoError := errors.New("database error")

	repo.On("GetByID", ctx, articleID).Return(nil, repoError)

	err := uc.Execute(ctx, articleID)

	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Delete")
}

func TestDeleteArticleUseCase_Execute_DeleteError(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewDeleteArticleUseCase(repo, cache, listCache)

	articleID := int64(1)
	existingArticle := &domainarticle.Article{
		ID:        articleID,
		Title:     "Test Article",
		Content:   "Test Content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	deleteError := errors.New("delete error")

	repo.On("GetByID", ctx, articleID).Return(existingArticle, nil)
	repo.On("Delete", ctx, articleID).Return(deleteError)

	err := uc.Execute(ctx, articleID)

	assert.Error(t, err)
	assert.Equal(t, deleteError, err)
	repo.AssertExpectations(t)
	cache.AssertNotCalled(t, "Delete")
}

func TestDeleteArticleUseCase_Execute_WithNilCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	listCache := &mockArticleListCache{}

	uc := NewDeleteArticleUseCase(repo, nil, listCache)

	articleID := int64(1)
	existingArticle := &domainarticle.Article{
		ID:        articleID,
		Title:     "Test Article",
		Content:   "Test Content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, articleID).Return(existingArticle, nil)
	repo.On("Delete", ctx, articleID).Return(nil)
	listCache.On("InvalidateArticleList", ctx).Return(nil)

	err := uc.Execute(ctx, articleID)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	listCache.AssertExpectations(t)
}

func TestDeleteArticleUseCase_Execute_WithNilListCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}

	uc := NewDeleteArticleUseCase(repo, cache, nil)

	articleID := int64(1)
	existingArticle := &domainarticle.Article{
		ID:        articleID,
		Title:     "Test Article",
		Content:   "Test Content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, articleID).Return(existingArticle, nil)
	repo.On("Delete", ctx, articleID).Return(nil)
	cache.On("Delete", ctx, articleID).Return(nil)
	cache.On("InvalidateList", ctx).Return(nil)

	err := uc.Execute(ctx, articleID)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}
