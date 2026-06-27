package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rulzi/hexa-go/internal/application/article/dto"
	articleentity "github.com/rulzi/hexa-go/internal/domain/article/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateArticle_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}

	uc := NewCreateArticle(repo, cache)

	req := dto.CreateArticleRequest{
		Title:    "Test Article",
		Content:  "Test Content",
		AuthorID: 1,
	}

	expectedArticle := &articleentity.Article{
		ID:        1,
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  req.AuthorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("Create", ctx, mock.Anything).Return(expectedArticle, nil)
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

func TestCreateArticle_Execute_ValidationError(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}

	uc := NewCreateArticle(repo, cache)

	tests := []struct {
		name string
		req  dto.CreateArticleRequest
	}{
		{"empty title", dto.CreateArticleRequest{Title: "", Content: "Test Content", AuthorID: 1}},
		{"empty content", dto.CreateArticleRequest{Title: "Test Title", Content: "", AuthorID: 1}},
		{"invalid author id", dto.CreateArticleRequest{Title: "Test Title", Content: "Test Content", AuthorID: 0}},
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

func TestCreateArticle_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}

	uc := NewCreateArticle(repo, cache)

	req := dto.CreateArticleRequest{Title: "Test Article", Content: "Test Content", AuthorID: 1}
	repoError := errors.New("repository error")
	repo.On("Create", ctx, mock.Anything).Return(nil, repoError)

	result, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
	cache.AssertNotCalled(t, "InvalidateList")
}

func TestCreateArticle_Execute_WithNoopCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}

	uc := NewCreateArticle(repo, NoopCache{})

	req := dto.CreateArticleRequest{Title: "Test Article", Content: "Test Content", AuthorID: 1}
	expectedArticle := &articleentity.Article{
		ID:        1,
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  req.AuthorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.On("Create", ctx, mock.Anything).Return(expectedArticle, nil)

	result, err := uc.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	repo.AssertExpectations(t)
}
