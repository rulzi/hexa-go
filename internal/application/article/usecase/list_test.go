package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rulzi/hexa-go/internal/application/article/dto"
	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewListArticlesUseCase(t *testing.T) {
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}
	dtoCache := &mockArticleListCache{}

	uc := NewListArticlesUseCase(repo, cache, dtoCache)

	assert.NotNil(t, uc)
	assert.Equal(t, repo, uc.articleRepo)
	assert.Equal(t, cache, uc.cache)
	assert.Equal(t, dtoCache, uc.dtoCache)
}

func TestListArticlesUseCase_Execute_SuccessFromCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}
	dtoCache := &mockArticleListCache{}

	uc := NewListArticlesUseCase(repo, cache, dtoCache)

	limit := 10
	offset := 0
	cachedResponse := &dto.ListArticlesResponse{
		Articles: []dto.ArticleResponse{
			{
				ID:        1,
				Title:     "Cached Article",
				Content:   "Cached Content",
				AuthorID:  1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		Total:  1,
		Limit: limit,
		Offset: offset,
	}

	dtoCache.On("GetArticleList", ctx, limit, offset).Return(cachedResponse, nil)

	result, err := uc.Execute(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cachedResponse.Total, result.Total)
	assert.Equal(t, len(cachedResponse.Articles), len(result.Articles))
	assert.Equal(t, cachedResponse.Articles[0].ID, result.Articles[0].ID)

	dtoCache.AssertExpectations(t)
	repo.AssertNotCalled(t, "List")
	repo.AssertNotCalled(t, "Count")
}

func TestListArticlesUseCase_Execute_SuccessFromRepository(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}
	dtoCache := &mockArticleListCache{}

	uc := NewListArticlesUseCase(repo, cache, dtoCache)

	limit := 10
	offset := 0
	articles := []*domainarticle.Article{
		{
			ID:        1,
			Title:     "Article 1",
			Content:   "Content 1",
			AuthorID:  1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Title:     "Article 2",
			Content:   "Content 2",
			AuthorID:  1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	total := int64(2)

	dtoCache.On("GetArticleList", ctx, limit, offset).Return(nil, errors.New("cache miss"))
	repo.On("List", ctx, limit, offset).Return(articles, nil)
	repo.On("Count", ctx).Return(total, nil)
	dtoCache.On("SetArticleList", ctx, limit, offset, mock.AnythingOfType("*dto.ListArticlesResponse")).Return(nil)

	result, err := uc.Execute(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, total, result.Total)
	assert.Equal(t, limit, result.Limit)
	assert.Equal(t, offset, result.Offset)
	assert.Equal(t, len(articles), len(result.Articles))
	assert.Equal(t, articles[0].ID, result.Articles[0].ID)
	assert.Equal(t, articles[1].ID, result.Articles[1].ID)

	dtoCache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestListArticlesUseCase_Execute_DefaultPagination(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}
	dtoCache := &mockArticleListCache{}

	uc := NewListArticlesUseCase(repo, cache, dtoCache)

	// Test with invalid limit and offset
	articles := []*domainarticle.Article{}
	total := int64(0)

	dtoCache.On("GetArticleList", ctx, 10, 0).Return(nil, errors.New("cache miss"))
	repo.On("List", ctx, 10, 0).Return(articles, nil)
	repo.On("Count", ctx).Return(total, nil)
	dtoCache.On("SetArticleList", ctx, 10, 0, mock.AnythingOfType("*dto.ListArticlesResponse")).Return(nil)

	result, err := uc.Execute(ctx, -1, -1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 0, result.Offset)

	dtoCache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestListArticlesUseCase_Execute_ListError(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}
	dtoCache := &mockArticleListCache{}

	uc := NewListArticlesUseCase(repo, cache, dtoCache)

	limit := 10
	offset := 0
	listError := errors.New("list error")

	dtoCache.On("GetArticleList", ctx, limit, offset).Return(nil, errors.New("cache miss"))
	repo.On("List", ctx, limit, offset).Return(nil, listError)

	result, err := uc.Execute(ctx, limit, offset)

	assert.Error(t, err)
	assert.Equal(t, listError, err)
	assert.Nil(t, result)

	dtoCache.AssertExpectations(t)
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Count")
}

func TestListArticlesUseCase_Execute_CountError(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}
	dtoCache := &mockArticleListCache{}

	uc := NewListArticlesUseCase(repo, cache, dtoCache)

	limit := 10
	offset := 0
	articles := []*domainarticle.Article{
		{
			ID:        1,
			Title:     "Article 1",
			Content:   "Content 1",
			AuthorID:  1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	countError := errors.New("count error")

	dtoCache.On("GetArticleList", ctx, limit, offset).Return(nil, errors.New("cache miss"))
	repo.On("List", ctx, limit, offset).Return(articles, nil)
	repo.On("Count", ctx).Return(int64(0), countError)

	result, err := uc.Execute(ctx, limit, offset)

	assert.Error(t, err)
	assert.Equal(t, countError, err)
	assert.Nil(t, result)

	dtoCache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestListArticlesUseCase_Execute_WithNilDtoCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}

	uc := NewListArticlesUseCase(repo, cache, nil)

	limit := 10
	offset := 0
	articles := []*domainarticle.Article{
		{
			ID:        1,
			Title:     "Article 1",
			Content:   "Content 1",
			AuthorID:  1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	total := int64(1)

	repo.On("List", ctx, limit, offset).Return(articles, nil)
	repo.On("Count", ctx).Return(total, nil)

	result, err := uc.Execute(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, total, result.Total)
	assert.Equal(t, len(articles), len(result.Articles))

	repo.AssertExpectations(t)
}

func TestListArticlesUseCase_Execute_EmptyList(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}
	dtoCache := &mockArticleListCache{}

	uc := NewListArticlesUseCase(repo, cache, dtoCache)

	limit := 10
	offset := 0
	articles := []*domainarticle.Article{}
	total := int64(0)

	dtoCache.On("GetArticleList", ctx, limit, offset).Return(nil, errors.New("cache miss"))
	repo.On("List", ctx, limit, offset).Return(articles, nil)
	repo.On("Count", ctx).Return(total, nil)
	dtoCache.On("SetArticleList", ctx, limit, offset, mock.AnythingOfType("*dto.ListArticlesResponse")).Return(nil)

	result, err := uc.Execute(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(0), result.Total)
	assert.Equal(t, 0, len(result.Articles))

	dtoCache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

