package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
	"github.com/stretchr/testify/assert"
)

func TestNewGetArticleUseCase(t *testing.T) {
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}

	uc := NewGetArticleUseCase(repo, cache)

	assert.NotNil(t, uc)
	assert.Equal(t, repo, uc.articleRepo)
	assert.Equal(t, cache, uc.cache)
}

func TestGetArticleUseCase_Execute_SuccessFromCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}

	uc := NewGetArticleUseCase(repo, cache)

	articleID := int64(1)
	cachedArticle := &domainarticle.Article{
		ID:        articleID,
		Title:     "Cached Article",
		Content:   "Cached Content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	cache.On("Get", ctx, articleID).Return(cachedArticle, nil)

	result, err := uc.Execute(ctx, articleID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cachedArticle.ID, result.ID)
	assert.Equal(t, cachedArticle.Title, result.Title)
	assert.Equal(t, cachedArticle.Content, result.Content)
	assert.Equal(t, cachedArticle.AuthorID, result.AuthorID)

	cache.AssertExpectations(t)
	repo.AssertNotCalled(t, "GetByID")
}

func TestGetArticleUseCase_Execute_SuccessFromRepository(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}

	uc := NewGetArticleUseCase(repo, cache)

	articleID := int64(1)
	articleEntity := &domainarticle.Article{
		ID:        articleID,
		Title:     "Test Article",
		Content:   "Test Content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	cache.On("Get", ctx, articleID).Return(nil, errors.New("cache miss"))
	repo.On("GetByID", ctx, articleID).Return(articleEntity, nil)
	cache.On("Set", ctx, articleID, articleEntity).Return(nil)

	result, err := uc.Execute(ctx, articleID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, articleEntity.ID, result.ID)
	assert.Equal(t, articleEntity.Title, result.Title)
	assert.Equal(t, articleEntity.Content, result.Content)
	assert.Equal(t, articleEntity.AuthorID, result.AuthorID)

	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestGetArticleUseCase_Execute_ArticleNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}

	uc := NewGetArticleUseCase(repo, cache)

	articleID := int64(1)

	cache.On("Get", ctx, articleID).Return(nil, errors.New("cache miss"))
	repo.On("GetByID", ctx, articleID).Return(nil, nil)

	result, err := uc.Execute(ctx, articleID)

	assert.Error(t, err)
	assert.Equal(t, domainarticle.ErrArticleNotFound, err)
	assert.Nil(t, result)

	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
	cache.AssertNotCalled(t, "Set")
}

func TestGetArticleUseCase_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}

	uc := NewGetArticleUseCase(repo, cache)

	articleID := int64(1)
	repoError := errors.New("database error")

	cache.On("Get", ctx, articleID).Return(nil, errors.New("cache miss"))
	repo.On("GetByID", ctx, articleID).Return(nil, repoError)

	result, err := uc.Execute(ctx, articleID)

	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	assert.Nil(t, result)

	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
	cache.AssertNotCalled(t, "Set")
}

func TestGetArticleUseCase_Execute_WithNilCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}

	uc := NewGetArticleUseCase(repo, nil)

	articleID := int64(1)
	articleEntity := &domainarticle.Article{
		ID:        articleID,
		Title:     "Test Article",
		Content:   "Test Content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, articleID).Return(articleEntity, nil)

	result, err := uc.Execute(ctx, articleID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, articleEntity.ID, result.ID)

	repo.AssertExpectations(t)
}

func TestGetArticleUseCase_Execute_CacheErrorButContinue(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}

	uc := NewGetArticleUseCase(repo, cache)

	articleID := int64(1)
	articleEntity := &domainarticle.Article{
		ID:        articleID,
		Title:     "Test Article",
		Content:   "Test Content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Cache returns error but we continue to repository
	cache.On("Get", ctx, articleID).Return(nil, errors.New("cache error"))
	repo.On("GetByID", ctx, articleID).Return(articleEntity, nil)
	cache.On("Set", ctx, articleID, articleEntity).Return(nil)

	result, err := uc.Execute(ctx, articleID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, articleEntity.ID, result.ID)

	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}
