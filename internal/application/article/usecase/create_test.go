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

func TestNewCreateArticleUseCase(t *testing.T) {
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}

	uc := NewCreateArticleUseCase(repo, service, cache)

	assert.NotNil(t, uc)
	assert.Equal(t, repo, uc.articleRepo)
	assert.Equal(t, service, uc.articleService)
	assert.Equal(t, cache, uc.cache)
}

func TestCreateArticleUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}

	uc := NewCreateArticleUseCase(repo, service, cache)

	req := dto.CreateArticleRequest{
		Title:    "Test Article",
		Content:  "Test Content",
		AuthorID: 1,
	}

	expectedArticle := &domainarticle.Article{
		ID:        1,
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  req.AuthorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("Create", ctx, mock.AnythingOfType("*article.Article")).Return(expectedArticle, nil)
	cache.On("InvalidateList", ctx).Return(nil)

	result, err := uc.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedArticle.ID, result.ID)
	assert.Equal(t, expectedArticle.Title, result.Title)
	assert.Equal(t, expectedArticle.Content, result.Content)
	assert.Equal(t, expectedArticle.AuthorID, result.AuthorID)

	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestCreateArticleUseCase_Execute_ValidationError(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}

	uc := NewCreateArticleUseCase(repo, service, cache)

	tests := []struct {
		name string
		req  dto.CreateArticleRequest
	}{
		{
			name: "empty title",
			req: dto.CreateArticleRequest{
				Title:    "",
				Content:  "Test Content",
				AuthorID: 1,
			},
		},
		{
			name: "empty content",
			req: dto.CreateArticleRequest{
				Title:    "Test Title",
				Content:  "",
				AuthorID: 1,
			},
		},
		{
			name: "invalid author id",
			req: dto.CreateArticleRequest{
				Title:    "Test Title",
				Content:  "Test Content",
				AuthorID: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := uc.Execute(ctx, tt.req)

			assert.Error(t, err)
			assert.Nil(t, result)
			repo.AssertNotCalled(t, "Create")
		})
	}
}

func TestCreateArticleUseCase_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}

	uc := NewCreateArticleUseCase(repo, service, cache)

	req := dto.CreateArticleRequest{
		Title:    "Test Article",
		Content:  "Test Content",
		AuthorID: 1,
	}

	repoError := errors.New("repository error")
	repo.On("Create", ctx, mock.AnythingOfType("*article.Article")).Return(nil, repoError)

	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
	cache.AssertNotCalled(t, "InvalidateList")
}

func TestCreateArticleUseCase_Execute_WithNilCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)

	uc := NewCreateArticleUseCase(repo, service, nil)

	req := dto.CreateArticleRequest{
		Title:    "Test Article",
		Content:  "Test Content",
		AuthorID: 1,
	}

	expectedArticle := &domainarticle.Article{
		ID:        1,
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  req.AuthorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("Create", ctx, mock.AnythingOfType("*article.Article")).Return(expectedArticle, nil)

	result, err := uc.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	repo.AssertExpectations(t)
}
