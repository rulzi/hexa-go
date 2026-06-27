package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rulzi/hexa-go/internal/application/article/dto"
	articleentity "github.com/rulzi/hexa-go/internal/domain/article/entity"
	articleservice "github.com/rulzi/hexa-go/internal/domain/article/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListArticle_Execute_SuccessFromCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := articleservice.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, service, cache, listCache)

	limit, offset := 10, 0
	cachedResponse := &dto.ListArticlesResponse{
		Articles: []dto.ArticleResponse{
			{ID: 1, Title: "Cached Article", Content: "Cached Content", AuthorID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
		Total: 1, Limit: limit, Offset: offset,
	}
	listCache.On("GetArticleList", ctx, limit, offset).Return(cachedResponse, nil)

	result, err := uc.List.Execute(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cachedResponse.Total, result.Total)
	assert.Equal(t, len(cachedResponse.Articles), len(result.Articles))
	listCache.AssertExpectations(t)
	repo.AssertNotCalled(t, "List")
}

func TestListArticle_Execute_SuccessFromRepository(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := articleservice.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, service, cache, listCache)

	limit, offset := 10, 0
	articles := []*articleentity.Article{
		{ID: 1, Title: "Article 1", Content: "Content 1", AuthorID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Title: "Article 2", Content: "Content 2", AuthorID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	total := int64(2)
	listCache.On("GetArticleList", ctx, limit, offset).Return(nil, errors.New("cache miss"))
	repo.On("List", ctx, limit, offset).Return(articles, nil)
	repo.On("Count", ctx).Return(total, nil)
	listCache.On("SetArticleList", ctx, limit, offset, mock.AnythingOfType("*dto.ListArticlesResponse")).Return(nil)

	result, err := uc.List.Execute(ctx, limit, offset)

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
	service := articleservice.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, service, cache, listCache)

	listCache.On("GetArticleList", ctx, 10, 0).Return(nil, errors.New("cache miss"))
	repo.On("List", ctx, 10, 0).Return([]*articleentity.Article{}, nil)
	repo.On("Count", ctx).Return(int64(0), nil)
	listCache.On("SetArticleList", ctx, 10, 0, mock.AnythingOfType("*dto.ListArticlesResponse")).Return(nil)

	result, err := uc.List.Execute(ctx, -1, -1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 0, result.Offset)
}

func TestListArticle_Execute_WithNilListCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := articleservice.NewService(repo)
	cache := &mockArticleCache{}

	uc := NewArticleUseCase(repo, service, cache, nil)

	limit, offset := 10, 0
	articles := []*articleentity.Article{
		{ID: 1, Title: "Article 1", Content: "Content 1", AuthorID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	total := int64(1)
	repo.On("List", ctx, limit, offset).Return(articles, nil)
	repo.On("Count", ctx).Return(total, nil)

	result, err := uc.List.Execute(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, total, result.Total)
	repo.AssertExpectations(t)
}
