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

func TestNewUpdateArticleUseCase(t *testing.T) {
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewUpdateArticleUseCase(repo, service, cache, listCache)

	assert.NotNil(t, uc)
	assert.Equal(t, repo, uc.articleRepo)
	assert.Equal(t, service, uc.articleService)
	assert.Equal(t, cache, uc.cache)
	assert.Equal(t, listCache, uc.listCache)
}

func TestUpdateArticleUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewUpdateArticleUseCase(repo, service, cache, listCache)

	articleID := int64(1)
	existingArticle := &domainarticle.Article{
		ID:        articleID,
		Title:     "Old Title",
		Content:   "Old Content",
		AuthorID:  1,
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
	}

	req := dto.UpdateArticleRequest{
		Title:   "New Title",
		Content: "New Content",
	}

	updatedArticle := &domainarticle.Article{
		ID:        articleID,
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  existingArticle.AuthorID,
		CreatedAt: existingArticle.CreatedAt,
		UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, articleID).Return(existingArticle, nil)
	repo.On("Update", ctx, mock.AnythingOfType("*article.Article")).Return(updatedArticle, nil)
	cache.On("Delete", ctx, articleID).Return(nil)
	cache.On("InvalidateList", ctx).Return(nil)
	listCache.On("InvalidateArticleList", ctx).Return(nil)

	result, err := uc.Execute(ctx, articleID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updatedArticle.ID, result.ID)
	assert.Equal(t, req.Title, result.Title)
	assert.Equal(t, req.Content, result.Content)
	assert.Equal(t, existingArticle.AuthorID, result.AuthorID)

	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
	listCache.AssertExpectations(t)
}

func TestUpdateArticleUseCase_Execute_ArticleNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewUpdateArticleUseCase(repo, service, cache, listCache)

	articleID := int64(1)
	req := dto.UpdateArticleRequest{
		Title:   "New Title",
		Content: "New Content",
	}

	repo.On("GetByID", ctx, articleID).Return(nil, nil)

	result, err := uc.Execute(ctx, articleID, req)

	assert.Error(t, err)
	assert.Equal(t, domainarticle.ErrArticleNotFound, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Update")
	cache.AssertNotCalled(t, "Delete")
}

func TestUpdateArticleUseCase_Execute_GetByIDError(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewUpdateArticleUseCase(repo, service, cache, listCache)

	articleID := int64(1)
	req := dto.UpdateArticleRequest{
		Title:   "New Title",
		Content: "New Content",
	}
	repoError := errors.New("database error")

	repo.On("GetByID", ctx, articleID).Return(nil, repoError)

	result, err := uc.Execute(ctx, articleID, req)

	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Update")
}

func TestUpdateArticleUseCase_Execute_ValidationError(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewUpdateArticleUseCase(repo, service, cache, listCache)

	articleID := int64(1)
	existingArticle := &domainarticle.Article{
		ID:        articleID,
		Title:     "Old Title",
		Content:   "Old Content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name string
		req  dto.UpdateArticleRequest
	}{
		{
			name: "empty title",
			req: dto.UpdateArticleRequest{
				Title:   "",
				Content: "New Content",
			},
		},
		{
			name: "empty content",
			req: dto.UpdateArticleRequest{
				Title:   "New Title",
				Content: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.On("GetByID", ctx, articleID).Return(existingArticle, nil)

			result, err := uc.Execute(ctx, articleID, tt.req)

			assert.Error(t, err)
			assert.Nil(t, result)
			repo.AssertNotCalled(t, "Update")
		})
	}
}

func TestUpdateArticleUseCase_Execute_UpdateError(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewUpdateArticleUseCase(repo, service, cache, listCache)

	articleID := int64(1)
	existingArticle := &domainarticle.Article{
		ID:        articleID,
		Title:     "Old Title",
		Content:   "Old Content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := dto.UpdateArticleRequest{
		Title:   "New Title",
		Content: "New Content",
	}

	updateError := errors.New("update error")

	repo.On("GetByID", ctx, articleID).Return(existingArticle, nil)
	repo.On("Update", ctx, mock.AnythingOfType("*article.Article")).Return(nil, updateError)

	result, err := uc.Execute(ctx, articleID, req)

	assert.Error(t, err)
	assert.Equal(t, updateError, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
	cache.AssertNotCalled(t, "Delete")
}

func TestUpdateArticleUseCase_Execute_WithNilCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	listCache := &mockArticleListCache{}

	uc := NewUpdateArticleUseCase(repo, service, nil, listCache)

	articleID := int64(1)
	existingArticle := &domainarticle.Article{
		ID:        articleID,
		Title:     "Old Title",
		Content:   "Old Content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := dto.UpdateArticleRequest{
		Title:   "New Title",
		Content: "New Content",
	}

	updatedArticle := &domainarticle.Article{
		ID:        articleID,
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  existingArticle.AuthorID,
		CreatedAt: existingArticle.CreatedAt,
		UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, articleID).Return(existingArticle, nil)
	repo.On("Update", ctx, mock.AnythingOfType("*article.Article")).Return(updatedArticle, nil)
	listCache.On("InvalidateArticleList", ctx).Return(nil)

	result, err := uc.Execute(ctx, articleID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Title, result.Title)

	repo.AssertExpectations(t)
	listCache.AssertExpectations(t)
}

func TestUpdateArticleUseCase_Execute_WithNilListCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}

	uc := NewUpdateArticleUseCase(repo, service, cache, nil)

	articleID := int64(1)
	existingArticle := &domainarticle.Article{
		ID:        articleID,
		Title:     "Old Title",
		Content:   "Old Content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := dto.UpdateArticleRequest{
		Title:   "New Title",
		Content: "New Content",
	}

	updatedArticle := &domainarticle.Article{
		ID:        articleID,
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  existingArticle.AuthorID,
		CreatedAt: existingArticle.CreatedAt,
		UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, articleID).Return(existingArticle, nil)
	repo.On("Update", ctx, mock.AnythingOfType("*article.Article")).Return(updatedArticle, nil)
	cache.On("Delete", ctx, articleID).Return(nil)
	cache.On("InvalidateList", ctx).Return(nil)

	result, err := uc.Execute(ctx, articleID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Title, result.Title)

	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

