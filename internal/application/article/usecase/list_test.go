package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	articleentity "github.com/rulzi/hexa-go/internal/domain/article/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListArticle_Execute_SuccessFromCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	listCache := &mockArticleListCache{}

	uc := NewListArticle(repo, listCache)

	limit, offset := 10, 0
	cachedPage := &ArticleListPage{
		Articles: []*articleentity.Article{
			{ID: 1, Title: "Cached Article", Content: "Cached Content", AuthorID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
		Total: 1,
	}
	listCache.On("GetList", ctx, limit, offset).Return(cachedPage, nil)

	result, err := uc.Execute(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cachedPage.Total, result.Total)
	assert.Equal(t, len(cachedPage.Articles), len(result.Articles))
	listCache.AssertExpectations(t)
	repo.AssertNotCalled(t, "List")
}

func TestListArticle_Execute_SuccessFromRepository(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	listCache := &mockArticleListCache{}

	uc := NewListArticle(repo, listCache)

	limit, offset := 10, 0
	articles := []*articleentity.Article{
		{ID: 1, Title: "Article 1", Content: "Content 1", AuthorID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Title: "Article 2", Content: "Content 2", AuthorID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	total := int64(2)
	listCache.On("GetList", ctx, limit, offset).Return(nil, errors.New("cache miss"))
	repo.On("List", ctx, limit, offset).Return(articles, nil)
	repo.On("Count", ctx).Return(total, nil)
	listCache.On("SetList", ctx, limit, offset, mock.AnythingOfType("*usecase.ArticleListPage")).Return(nil)

	result, err := uc.Execute(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, total, result.Total)
	assert.Equal(t, limit, result.Limit)
	assert.Equal(t, offset, result.Offset)
	assert.Equal(t, len(articles), len(result.Articles))
}

func TestListArticle_Execute_DefaultPagination(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	listCache := &mockArticleListCache{}

	uc := NewListArticle(repo, listCache)

	listCache.On("GetList", ctx, 10, 0).Return(nil, errors.New("cache miss"))
	repo.On("List", ctx, 10, 0).Return([]*articleentity.Article{}, nil)
	repo.On("Count", ctx).Return(int64(0), nil)
	listCache.On("SetList", ctx, 10, 0, mock.AnythingOfType("*usecase.ArticleListPage")).Return(nil)

	result, err := uc.Execute(ctx, -1, -1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 0, result.Offset)
}

func TestListArticle_Execute_WithNoopCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}

	uc := NewListArticle(repo, NoopCache{})

	limit, offset := 10, 0
	articles := []*articleentity.Article{
		{ID: 1, Title: "Article 1", Content: "Content 1", AuthorID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	total := int64(1)
	repo.On("List", ctx, limit, offset).Return(articles, nil)
	repo.On("Count", ctx).Return(total, nil)

	result, err := uc.Execute(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, total, result.Total)
	repo.AssertExpectations(t)
}
