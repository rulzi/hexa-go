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

func TestNewArticleUseCase(t *testing.T) {
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, service, cache, listCache)

	assert.NotNil(t, uc)
	assert.Equal(t, repo, uc.articleRepo)
	assert.Equal(t, service, uc.articleService)
	assert.Equal(t, cache, uc.cache)
	assert.Equal(t, listCache, uc.listCache)
}

// --- Create ---

func TestArticleUseCase_Create_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, service, cache, listCache)

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

	result, err := uc.Create(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedArticle.ID, result.ID)
	assert.Equal(t, expectedArticle.Title, result.Title)
	assert.Equal(t, expectedArticle.Content, result.Content)
	assert.Equal(t, expectedArticle.AuthorID, result.AuthorID)

	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestArticleUseCase_Create_ValidationError(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, service, cache, listCache)

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
			result, err := uc.Create(ctx, tt.req)
			assert.Error(t, err)
			assert.Nil(t, result)
			repo.AssertNotCalled(t, "Create")
		})
	}
}

func TestArticleUseCase_Create_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, service, cache, listCache)

	req := dto.CreateArticleRequest{Title: "Test Article", Content: "Test Content", AuthorID: 1}
	repoError := errors.New("repository error")
	repo.On("Create", ctx, mock.AnythingOfType("*article.Article")).Return(nil, repoError)

	result, err := uc.Create(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
	cache.AssertNotCalled(t, "InvalidateList")
}

func TestArticleUseCase_Create_WithNilCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)

	uc := NewArticleUseCase(repo, service, nil, nil)

	req := dto.CreateArticleRequest{Title: "Test Article", Content: "Test Content", AuthorID: 1}
	expectedArticle := &domainarticle.Article{
		ID:        1,
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  req.AuthorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.On("Create", ctx, mock.AnythingOfType("*article.Article")).Return(expectedArticle, nil)

	result, err := uc.Create(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	repo.AssertExpectations(t)
}

// --- Get ---

func TestArticleUseCase_Get_SuccessFromCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, service, cache, listCache)

	articleID := int64(1)
	cachedArticle := &domainarticle.Article{
		ID: articleID, Title: "Cached Article", Content: "Cached Content", AuthorID: 1,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	cache.On("Get", ctx, articleID).Return(cachedArticle, nil)

	result, err := uc.Get(ctx, articleID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cachedArticle.ID, result.ID)
	assert.Equal(t, cachedArticle.Title, result.Title)
	cache.AssertExpectations(t)
	repo.AssertNotCalled(t, "GetByID")
}

func TestArticleUseCase_Get_SuccessFromRepository(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, service, cache, listCache)

	articleID := int64(1)
	articleEntity := &domainarticle.Article{
		ID: articleID, Title: "Test Article", Content: "Test Content", AuthorID: 1,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	cache.On("Get", ctx, articleID).Return(nil, errors.New("cache miss"))
	repo.On("GetByID", ctx, articleID).Return(articleEntity, nil)
	cache.On("Set", ctx, articleID, articleEntity).Return(nil)

	result, err := uc.Get(ctx, articleID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, articleEntity.ID, result.ID)
	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestArticleUseCase_Get_ArticleNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, service, cache, listCache)

	articleID := int64(1)
	cache.On("Get", ctx, articleID).Return(nil, errors.New("cache miss"))
	repo.On("GetByID", ctx, articleID).Return(nil, nil)

	result, err := uc.Get(ctx, articleID)

	assert.Error(t, err)
	assert.Equal(t, domainarticle.ErrArticleNotFound, err)
	assert.Nil(t, result)
}

func TestArticleUseCase_Get_WithNilCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)

	uc := NewArticleUseCase(repo, service, nil, nil)

	articleID := int64(1)
	articleEntity := &domainarticle.Article{
		ID: articleID, Title: "Test Article", Content: "Test Content", AuthorID: 1,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	repo.On("GetByID", ctx, articleID).Return(articleEntity, nil)

	result, err := uc.Get(ctx, articleID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, articleEntity.ID, result.ID)
	repo.AssertExpectations(t)
}

// --- List ---

func TestArticleUseCase_List_SuccessFromCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
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

	result, err := uc.List(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cachedResponse.Total, result.Total)
	assert.Equal(t, len(cachedResponse.Articles), len(result.Articles))
	listCache.AssertExpectations(t)
	repo.AssertNotCalled(t, "List")
}

func TestArticleUseCase_List_SuccessFromRepository(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, service, cache, listCache)

	limit, offset := 10, 0
	articles := []*domainarticle.Article{
		{ID: 1, Title: "Article 1", Content: "Content 1", AuthorID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Title: "Article 2", Content: "Content 2", AuthorID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	total := int64(2)
	listCache.On("GetArticleList", ctx, limit, offset).Return(nil, errors.New("cache miss"))
	repo.On("List", ctx, limit, offset).Return(articles, nil)
	repo.On("Count", ctx).Return(total, nil)
	listCache.On("SetArticleList", ctx, limit, offset, mock.AnythingOfType("*dto.ListArticlesResponse")).Return(nil)

	result, err := uc.List(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, total, result.Total)
	assert.Equal(t, limit, result.Limit)
	assert.Equal(t, offset, result.Offset)
	assert.Equal(t, len(articles), len(result.Articles))
}

func TestArticleUseCase_List_DefaultPagination(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, service, cache, listCache)

	listCache.On("GetArticleList", ctx, 10, 0).Return(nil, errors.New("cache miss"))
	repo.On("List", ctx, 10, 0).Return([]*domainarticle.Article{}, nil)
	repo.On("Count", ctx).Return(int64(0), nil)
	listCache.On("SetArticleList", ctx, 10, 0, mock.AnythingOfType("*dto.ListArticlesResponse")).Return(nil)

	result, err := uc.List(ctx, -1, -1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 0, result.Offset)
}

func TestArticleUseCase_List_WithNilListCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}

	uc := NewArticleUseCase(repo, service, cache, nil)

	limit, offset := 10, 0
	articles := []*domainarticle.Article{
		{ID: 1, Title: "Article 1", Content: "Content 1", AuthorID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	total := int64(1)
	repo.On("List", ctx, limit, offset).Return(articles, nil)
	repo.On("Count", ctx).Return(total, nil)

	result, err := uc.List(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, total, result.Total)
	repo.AssertExpectations(t)
}

// --- Update ---

func TestArticleUseCase_Update_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, service, cache, listCache)

	articleID := int64(1)
	existingArticle := &domainarticle.Article{
		ID: articleID, Title: "Old Title", Content: "Old Content", AuthorID: 1,
		CreatedAt: time.Now().Add(-24 * time.Hour), UpdatedAt: time.Now().Add(-24 * time.Hour),
	}
	req := dto.UpdateArticleRequest{Title: "New Title", Content: "New Content"}
	updatedArticle := &domainarticle.Article{
		ID: articleID, Title: req.Title, Content: req.Content, AuthorID: existingArticle.AuthorID,
		CreatedAt: existingArticle.CreatedAt, UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, articleID).Return(existingArticle, nil)
	repo.On("Update", ctx, mock.AnythingOfType("*article.Article")).Return(updatedArticle, nil)
	cache.On("Delete", ctx, articleID).Return(nil)
	cache.On("InvalidateList", ctx).Return(nil)
	listCache.On("InvalidateArticleList", ctx).Return(nil)

	result, err := uc.Update(ctx, articleID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Title, result.Title)
	assert.Equal(t, req.Content, result.Content)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
	listCache.AssertExpectations(t)
}

func TestArticleUseCase_Update_ArticleNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, service, cache, listCache)

	articleID := int64(1)
	req := dto.UpdateArticleRequest{Title: "New Title", Content: "New Content"}
	repo.On("GetByID", ctx, articleID).Return(nil, nil)

	result, err := uc.Update(ctx, articleID, req)

	assert.Error(t, err)
	assert.Equal(t, domainarticle.ErrArticleNotFound, err)
	assert.Nil(t, result)
	repo.AssertNotCalled(t, "Update")
}

func TestArticleUseCase_Update_ValidationError(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, service, cache, listCache)

	articleID := int64(1)
	existingArticle := &domainarticle.Article{
		ID: articleID, Title: "Old Title", Content: "Old Content", AuthorID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	tests := []struct {
		name string
		req  dto.UpdateArticleRequest
	}{
		{"empty title", dto.UpdateArticleRequest{Title: "", Content: "New Content"}},
		{"empty content", dto.UpdateArticleRequest{Title: "New Title", Content: ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.On("GetByID", ctx, articleID).Return(existingArticle, nil)
			result, err := uc.Update(ctx, articleID, tt.req)
			assert.Error(t, err)
			assert.Nil(t, result)
			repo.AssertNotCalled(t, "Update")
		})
	}
}

// --- Delete ---

func TestArticleUseCase_Delete_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, service, cache, listCache)

	articleID := int64(1)
	existingArticle := &domainarticle.Article{
		ID: articleID, Title: "Test Article", Content: "Test Content", AuthorID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	repo.On("GetByID", ctx, articleID).Return(existingArticle, nil)
	repo.On("Delete", ctx, articleID).Return(nil)
	cache.On("Delete", ctx, articleID).Return(nil)
	cache.On("InvalidateList", ctx).Return(nil)
	listCache.On("InvalidateArticleList", ctx).Return(nil)

	err := uc.Delete(ctx, articleID)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
	listCache.AssertExpectations(t)
}

func TestArticleUseCase_Delete_ArticleNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, service, cache, listCache)

	articleID := int64(1)
	repo.On("GetByID", ctx, articleID).Return(nil, nil)

	err := uc.Delete(ctx, articleID)

	assert.Error(t, err)
	assert.Equal(t, domainarticle.ErrArticleNotFound, err)
	repo.AssertNotCalled(t, "Delete")
}

func TestArticleUseCase_Delete_WithNilCache(t *testing.T) {
	ctx := context.Background()
	repo := &mockArticleRepository{}
	service := domainarticle.NewService(repo)
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, service, nil, listCache)

	articleID := int64(1)
	existingArticle := &domainarticle.Article{
		ID: articleID, Title: "Test Article", Content: "Test Content", AuthorID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	repo.On("GetByID", ctx, articleID).Return(existingArticle, nil)
	repo.On("Delete", ctx, articleID).Return(nil)
	listCache.On("InvalidateArticleList", ctx).Return(nil)

	err := uc.Delete(ctx, articleID)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	listCache.AssertExpectations(t)
}
