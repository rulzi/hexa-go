package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	articleentity "github.com/rulzi/hexa-go/internal/domain/article/entity"
	"github.com/stretchr/testify/assert"
)

func TestGetArticle_Execute_SuccessFromCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}

	uc := NewGetArticle(repo, cache)

	articleID := int64(1)
	cachedArticle := &articleentity.Article{
		ID: articleID, Title: "Cached Article", Content: "Cached Content", AuthorID: 1,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	cache.On("Get", ctx, articleID).Return(cachedArticle, nil)

	result, err := uc.Execute(ctx, articleID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cachedArticle.ID, result.ID)
	assert.Equal(t, cachedArticle.Title, result.Title)
	cache.AssertExpectations(t)
	repo.AssertNotCalled(t, "GetByID")
}

func TestGetArticle_Execute_SuccessFromRepository(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}

	uc := NewGetArticle(repo, cache)

	articleID := int64(1)
	articleEntity := &articleentity.Article{
		ID: articleID, Title: "Test Article", Content: "Test Content", AuthorID: 1,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	cache.On("Get", ctx, articleID).Return(nil, errors.New("cache miss"))
	repo.On("GetByID", ctx, articleID).Return(articleEntity, nil)
	cache.On("Set", ctx, articleID, articleEntity).Return(nil)

	result, err := uc.Execute(ctx, articleID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, articleEntity.ID, result.ID)
	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestGetArticle_Execute_ArticleNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}

	uc := NewGetArticle(repo, cache)

	articleID := int64(1)
	cache.On("Get", ctx, articleID).Return(nil, errors.New("cache miss"))
	repo.On("GetByID", ctx, articleID).Return(nil, nil)

	result, err := uc.Execute(ctx, articleID)

	assert.Error(t, err)
	assert.True(t, articleentity.IsArticleNotFound(err))
	assert.Nil(t, result)
}

func TestGetArticle_Execute_WithNilCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}

	uc := NewGetArticle(repo, nil)

	articleID := int64(1)
	articleEntity := &articleentity.Article{
		ID: articleID, Title: "Test Article", Content: "Test Content", AuthorID: 1,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	repo.On("GetByID", ctx, articleID).Return(articleEntity, nil)

	result, err := uc.Execute(ctx, articleID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, articleEntity.ID, result.ID)
	repo.AssertExpectations(t)
}
