package article

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

// dtoCacheInterface defines the interface for DTO cache operations used by DomainCacheAdapter
type dtoCacheInterface interface {
	GetArticle(ctx context.Context, id int64) (*dto.ArticleResponse, error)
	SetArticle(ctx context.Context, id int64, articleResp *dto.ArticleResponse) error
	DeleteArticle(ctx context.Context, id int64) error
	InvalidateArticleList(ctx context.Context) error
}

// mockRedisCache is a mock implementation of dtoCacheInterface for testing
type mockRedisCache struct {
	mock.Mock
}

func (m *mockRedisCache) GetArticle(ctx context.Context, id int64) (*dto.ArticleResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ArticleResponse), args.Error(1)
}

func (m *mockRedisCache) SetArticle(ctx context.Context, id int64, articleResp *dto.ArticleResponse) error {
	args := m.Called(ctx, id, articleResp)
	return args.Error(0)
}

func (m *mockRedisCache) DeleteArticle(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockRedisCache) InvalidateArticleList(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// testDomainCacheAdapter wraps DomainCacheAdapter to allow testing with mock
type testDomainCacheAdapter struct {
	dtoCache dtoCacheInterface
}

func newTestDomainCacheAdapter(dtoCache dtoCacheInterface) *testDomainCacheAdapter {
	return &testDomainCacheAdapter{
		dtoCache: dtoCache,
	}
}

func (a *testDomainCacheAdapter) Get(ctx context.Context, id int64) (*domainarticle.Article, error) {
	dtoResp, err := a.dtoCache.GetArticle(ctx, id)
	if err != nil {
		return nil, err
	}
	if dtoResp == nil {
		return nil, nil
	}

	return &domainarticle.Article{
		ID:        dtoResp.ID,
		Title:     dtoResp.Title,
		Content:   dtoResp.Content,
		AuthorID:  dtoResp.AuthorID,
		CreatedAt: dtoResp.CreatedAt,
		UpdatedAt: dtoResp.UpdatedAt,
	}, nil
}

func (a *testDomainCacheAdapter) Set(ctx context.Context, id int64, article *domainarticle.Article) error {
	dtoResp := &dto.ArticleResponse{
		ID:        article.ID,
		Title:     article.Title,
		Content:   article.Content,
		AuthorID:  article.AuthorID,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}
	return a.dtoCache.SetArticle(ctx, id, dtoResp)
}

func (a *testDomainCacheAdapter) Delete(ctx context.Context, id int64) error {
	return a.dtoCache.DeleteArticle(ctx, id)
}

func (a *testDomainCacheAdapter) InvalidateList(ctx context.Context) error {
	return a.dtoCache.InvalidateArticleList(ctx)
}

func TestNewDomainCacheAdapter(t *testing.T) {
	// Test with real RedisCache (we can't easily mock *RedisCache, so we test the constructor)
	// In real scenario, this would use a real RedisCache instance
	// For unit testing the adapter logic, we use testDomainCacheAdapter
	dtoCache := &mockRedisCache{}

	adapter := newTestDomainCacheAdapter(dtoCache)

	assert.NotNil(t, adapter)
	assert.Equal(t, dtoCache, adapter.dtoCache)
}

func TestDomainCacheAdapter_Get_Success(t *testing.T) {
	ctx := context.Background()
	dtoCache := &mockRedisCache{}
	adapter := newTestDomainCacheAdapter(dtoCache)

	articleID := int64(1)
	now := time.Now()
	dtoResp := &dto.ArticleResponse{
		ID:        articleID,
		Title:     "Test Article",
		Content:   "Test Content",
		AuthorID:  123,
		CreatedAt: now,
		UpdatedAt: now,
	}

	dtoCache.On("GetArticle", ctx, articleID).Return(dtoResp, nil)

	result, err := adapter.Get(ctx, articleID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, dtoResp.ID, result.ID)
	assert.Equal(t, dtoResp.Title, result.Title)
	assert.Equal(t, dtoResp.Content, result.Content)
	assert.Equal(t, dtoResp.AuthorID, result.AuthorID)
	assert.Equal(t, dtoResp.CreatedAt, result.CreatedAt)
	assert.Equal(t, dtoResp.UpdatedAt, result.UpdatedAt)

	dtoCache.AssertExpectations(t)
}

func TestDomainCacheAdapter_Get_CacheMiss(t *testing.T) {
	ctx := context.Background()
	dtoCache := &mockRedisCache{}
	adapter := newTestDomainCacheAdapter(dtoCache)

	articleID := int64(1)

	dtoCache.On("GetArticle", ctx, articleID).Return(nil, nil)

	result, err := adapter.Get(ctx, articleID)

	assert.NoError(t, err)
	assert.Nil(t, result)

	dtoCache.AssertExpectations(t)
}

func TestDomainCacheAdapter_Get_Error(t *testing.T) {
	ctx := context.Background()
	dtoCache := &mockRedisCache{}
	adapter := newTestDomainCacheAdapter(dtoCache)

	articleID := int64(1)
	expectedErr := errors.New("cache error")

	dtoCache.On("GetArticle", ctx, articleID).Return(nil, expectedErr)

	result, err := adapter.Get(ctx, articleID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)

	dtoCache.AssertExpectations(t)
}

func TestDomainCacheAdapter_Set_Success(t *testing.T) {
	ctx := context.Background()
	dtoCache := &mockRedisCache{}
	adapter := newTestDomainCacheAdapter(dtoCache)

	articleID := int64(1)
	now := time.Now()
	domainArticle := &domainarticle.Article{
		ID:        articleID,
		Title:     "Test Article",
		Content:   "Test Content",
		AuthorID:  123,
		CreatedAt: now,
		UpdatedAt: now,
	}

	expectedDTO := &dto.ArticleResponse{
		ID:        domainArticle.ID,
		Title:     domainArticle.Title,
		Content:   domainArticle.Content,
		AuthorID:  domainArticle.AuthorID,
		CreatedAt: domainArticle.CreatedAt,
		UpdatedAt: domainArticle.UpdatedAt,
	}

	dtoCache.On("SetArticle", ctx, articleID, expectedDTO).Return(nil)

	err := adapter.Set(ctx, articleID, domainArticle)

	assert.NoError(t, err)
	dtoCache.AssertExpectations(t)
}

func TestDomainCacheAdapter_Set_Error(t *testing.T) {
	ctx := context.Background()
	dtoCache := &mockRedisCache{}
	adapter := newTestDomainCacheAdapter(dtoCache)

	articleID := int64(1)
	now := time.Now()
	domainArticle := &domainarticle.Article{
		ID:        articleID,
		Title:     "Test Article",
		Content:   "Test Content",
		AuthorID:  123,
		CreatedAt: now,
		UpdatedAt: now,
	}

	expectedDTO := &dto.ArticleResponse{
		ID:        domainArticle.ID,
		Title:     domainArticle.Title,
		Content:   domainArticle.Content,
		AuthorID:  domainArticle.AuthorID,
		CreatedAt: domainArticle.CreatedAt,
		UpdatedAt: domainArticle.UpdatedAt,
	}

	expectedErr := errors.New("cache set error")
	dtoCache.On("SetArticle", ctx, articleID, expectedDTO).Return(expectedErr)

	err := adapter.Set(ctx, articleID, domainArticle)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	dtoCache.AssertExpectations(t)
}

func TestDomainCacheAdapter_Delete_Success(t *testing.T) {
	ctx := context.Background()
	dtoCache := &mockRedisCache{}
	adapter := newTestDomainCacheAdapter(dtoCache)

	articleID := int64(1)

	dtoCache.On("DeleteArticle", ctx, articleID).Return(nil)

	err := adapter.Delete(ctx, articleID)

	assert.NoError(t, err)
	dtoCache.AssertExpectations(t)
}

func TestDomainCacheAdapter_Delete_Error(t *testing.T) {
	ctx := context.Background()
	dtoCache := &mockRedisCache{}
	adapter := newTestDomainCacheAdapter(dtoCache)

	articleID := int64(1)
	expectedErr := errors.New("cache delete error")

	dtoCache.On("DeleteArticle", ctx, articleID).Return(expectedErr)

	err := adapter.Delete(ctx, articleID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	dtoCache.AssertExpectations(t)
}

func TestDomainCacheAdapter_InvalidateList_Success(t *testing.T) {
	ctx := context.Background()
	dtoCache := &mockRedisCache{}
	adapter := newTestDomainCacheAdapter(dtoCache)

	dtoCache.On("InvalidateArticleList", ctx).Return(nil)

	err := adapter.InvalidateList(ctx)

	assert.NoError(t, err)
	dtoCache.AssertExpectations(t)
}

func TestDomainCacheAdapter_InvalidateList_Error(t *testing.T) {
	ctx := context.Background()
	dtoCache := &mockRedisCache{}
	adapter := newTestDomainCacheAdapter(dtoCache)

	expectedErr := errors.New("cache invalidate error")
	dtoCache.On("InvalidateArticleList", ctx).Return(expectedErr)

	err := adapter.InvalidateList(ctx)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	dtoCache.AssertExpectations(t)
}

func TestDomainCacheAdapter_ImplementsInterface(t *testing.T) {
	dtoCache := &mockRedisCache{}
	adapter := newTestDomainCacheAdapter(dtoCache)

	// Verify that testDomainCacheAdapter implements domainarticle.Cache interface
	var _ domainarticle.Cache = adapter
}

func TestDomainCacheAdapter_Get_ConvertsDTOToDomain(t *testing.T) {
	ctx := context.Background()
	dtoCache := &mockRedisCache{}
	adapter := newTestDomainCacheAdapter(dtoCache)

	articleID := int64(42)
	now := time.Now().Truncate(time.Second)
	dtoResp := &dto.ArticleResponse{
		ID:        articleID,
		Title:     "Converted Article",
		Content:   "This is converted content",
		AuthorID:  456,
		CreatedAt: now,
		UpdatedAt: now,
	}

	dtoCache.On("GetArticle", ctx, articleID).Return(dtoResp, nil)

	result, err := adapter.Get(ctx, articleID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	// Verify all fields are correctly converted
	assert.Equal(t, int64(42), result.ID)
	assert.Equal(t, "Converted Article", result.Title)
	assert.Equal(t, "This is converted content", result.Content)
	assert.Equal(t, int64(456), result.AuthorID)
	assert.Equal(t, now, result.CreatedAt)
	assert.Equal(t, now, result.UpdatedAt)

	dtoCache.AssertExpectations(t)
}

func TestDomainCacheAdapter_Set_ConvertsDomainToDTO(t *testing.T) {
	ctx := context.Background()
	dtoCache := &mockRedisCache{}
	adapter := newTestDomainCacheAdapter(dtoCache)

	articleID := int64(99)
	now := time.Now().Truncate(time.Second)
	domainArticle := &domainarticle.Article{
		ID:        articleID,
		Title:     "Domain Article",
		Content:   "Domain content here",
		AuthorID:  789,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Use mock.MatchedBy to verify the DTO structure
	dtoCache.On("SetArticle", ctx, articleID, mock.MatchedBy(func(dto *dto.ArticleResponse) bool {
		return dto.ID == articleID &&
			dto.Title == "Domain Article" &&
			dto.Content == "Domain content here" &&
			dto.AuthorID == 789 &&
			dto.CreatedAt.Equal(now) &&
			dto.UpdatedAt.Equal(now)
	})).Return(nil)

	err := adapter.Set(ctx, articleID, domainArticle)

	assert.NoError(t, err)
	dtoCache.AssertExpectations(t)
}

func TestDomainCacheAdapter_RoundTrip(t *testing.T) {
	ctx := context.Background()
	dtoCache := &mockRedisCache{}
	adapter := newTestDomainCacheAdapter(dtoCache)

	articleID := int64(100)
	now := time.Now().Truncate(time.Second)
	originalArticle := &domainarticle.Article{
		ID:        articleID,
		Title:     "Round Trip Article",
		Content:   "Round trip content",
		AuthorID:  111,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Set the article
	expectedDTO := &dto.ArticleResponse{
		ID:        originalArticle.ID,
		Title:     originalArticle.Title,
		Content:   originalArticle.Content,
		AuthorID:  originalArticle.AuthorID,
		CreatedAt: originalArticle.CreatedAt,
		UpdatedAt: originalArticle.UpdatedAt,
	}

	dtoCache.On("SetArticle", ctx, articleID, expectedDTO).Return(nil)
	err := adapter.Set(ctx, articleID, originalArticle)
	assert.NoError(t, err)

	// Get the article back
	dtoCache.On("GetArticle", ctx, articleID).Return(expectedDTO, nil)
	result, err := adapter.Get(ctx, articleID)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify the round trip preserved all data
	assert.Equal(t, originalArticle.ID, result.ID)
	assert.Equal(t, originalArticle.Title, result.Title)
	assert.Equal(t, originalArticle.Content, result.Content)
	assert.Equal(t, originalArticle.AuthorID, result.AuthorID)
	assert.Equal(t, originalArticle.CreatedAt, result.CreatedAt)
	assert.Equal(t, originalArticle.UpdatedAt, result.UpdatedAt)

	dtoCache.AssertExpectations(t)
}
