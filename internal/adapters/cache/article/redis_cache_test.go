package article

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rulzi/hexa-go/internal/application/article/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupRedisCache creates a RedisCache instance with a miniredis server
func setupRedisCache(t *testing.T, ttl time.Duration) (*RedisCache, *miniredis.Miniredis, func()) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	cache := NewRedisCache(client, ttl)

	cleanup := func() {
		_ = client.Close()
		mr.Close()
	}

	return cache, mr, cleanup
}

func TestNewRedisCache_DefaultTTL(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer func() {
		_ = client.Close()
	}()

	cache := NewRedisCache(client, 0)
	assert.Equal(t, 5*time.Minute, cache.ttl)
}

func TestNewRedisCache_CustomTTL(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer func() {
		_ = client.Close()
	}()

	customTTL := 10 * time.Minute
	cache := NewRedisCache(client, customTTL)
	assert.Equal(t, customTTL, cache.ttl)
}

// Test GetArticle - Success
func TestRedisCache_GetArticle_Success(t *testing.T) {
	cache, mr, cleanup := setupRedisCache(t, 5*time.Minute)
	defer cleanup()

	ctx := context.Background()
	articleID := int64(1)
	expectedArticle := &dto.ArticleResponse{
		ID:        articleID,
		Title:     "Test Article",
		Content:   "Test Content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set article in cache first
	data, err := json.Marshal(expectedArticle)
	require.NoError(t, err)
	key := fmt.Sprintf("article:%d", articleID)
	err = mr.Set(key, string(data))
	require.NoError(t, err)

	// Get article from cache
	result, err := cache.GetArticle(ctx, articleID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedArticle.ID, result.ID)
	assert.Equal(t, expectedArticle.Title, result.Title)
	assert.Equal(t, expectedArticle.Content, result.Content)
	assert.Equal(t, expectedArticle.AuthorID, result.AuthorID)
}

// Test GetArticle - Cache Miss
func TestRedisCache_GetArticle_CacheMiss(t *testing.T) {
	cache, _, cleanup := setupRedisCache(t, 5*time.Minute)
	defer cleanup()

	ctx := context.Background()
	articleID := int64(999)

	// Try to get non-existent article
	result, err := cache.GetArticle(ctx, articleID)
	require.NoError(t, err)
	assert.Nil(t, result) // Cache miss should return nil, not error
}

// Test GetArticle - Error (invalid JSON)
func TestRedisCache_GetArticle_InvalidJSON(t *testing.T) {
	cache, mr, cleanup := setupRedisCache(t, 5*time.Minute)
	defer cleanup()

	ctx := context.Background()
	articleID := int64(1)

	// Set invalid JSON in cache
	key := fmt.Sprintf("article:%d", articleID)
	err := mr.Set(key, "invalid json string")
	require.NoError(t, err)

	// Try to get article - should fail on unmarshal
	result, err := cache.GetArticle(ctx, articleID)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal cached article")
}

// Test SetArticle - Success
func TestRedisCache_SetArticle_Success(t *testing.T) {
	cache, mr, cleanup := setupRedisCache(t, 5*time.Minute)
	defer cleanup()

	ctx := context.Background()
	articleID := int64(1)
	article := &dto.ArticleResponse{
		ID:        articleID,
		Title:     "Test Article",
		Content:   "Test Content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set article in cache
	err := cache.SetArticle(ctx, articleID, article)
	require.NoError(t, err)

	// Verify it was stored
	key := fmt.Sprintf("article:%d", articleID)
	val, err := mr.Get(key)
	require.NoError(t, err)
	assert.NotEmpty(t, val)

	// Verify the content
	var storedArticle dto.ArticleResponse
	err = json.Unmarshal([]byte(val), &storedArticle)
	require.NoError(t, err)
	assert.Equal(t, article.ID, storedArticle.ID)
	assert.Equal(t, article.Title, storedArticle.Title)
}

// Test SetArticle - Error (Redis error)
func TestRedisCache_SetArticle_RedisError(t *testing.T) {
	cache, mr, cleanup := setupRedisCache(t, 5*time.Minute)
	defer cleanup()

	ctx := context.Background()
	articleID := int64(1)
	article := &dto.ArticleResponse{
		ID:        articleID,
		Title:     "Test Article",
		Content:   "Test Content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Close Redis to simulate error
	mr.Close()

	// Try to set article - should fail
	err := cache.SetArticle(ctx, articleID, article)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to set cache")
}

// Test DeleteArticle - Success
func TestRedisCache_DeleteArticle_Success(t *testing.T) {
	cache, mr, cleanup := setupRedisCache(t, 5*time.Minute)
	defer cleanup()

	ctx := context.Background()
	articleID := int64(1)

	// Set article in cache first
	key := fmt.Sprintf("article:%d", articleID)
	err := mr.Set(key, "test data")
	require.NoError(t, err)
	exists := mr.Exists(key)
	assert.True(t, exists)

	// Delete article from cache
	err = cache.DeleteArticle(ctx, articleID)
	require.NoError(t, err)

	// Verify it was deleted
	exists = mr.Exists(key)
	assert.False(t, exists)
}

// Test DeleteArticle - Error (Redis error)
func TestRedisCache_DeleteArticle_RedisError(t *testing.T) {
	cache, mr, cleanup := setupRedisCache(t, 5*time.Minute)
	defer cleanup()

	ctx := context.Background()
	articleID := int64(1)

	// Close Redis to simulate error
	mr.Close()

	// Try to delete article - should fail
	err := cache.DeleteArticle(ctx, articleID)
	assert.Error(t, err)
}

// Test GetArticleList - Success
func TestRedisCache_GetArticleList_Success(t *testing.T) {
	cache, mr, cleanup := setupRedisCache(t, 5*time.Minute)
	defer cleanup()

	ctx := context.Background()
	limit := 10
	offset := 0
	expectedList := &dto.ListArticlesResponse{
		Articles: []dto.ArticleResponse{
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
		},
		Total:  2,
		Limit:  limit,
		Offset: offset,
	}

	// Set list in cache first
	data, err := json.Marshal(expectedList)
	require.NoError(t, err)
	key := fmt.Sprintf("article:list:%d:%d", limit, offset)
	err = mr.Set(key, string(data))
	require.NoError(t, err)

	// Get list from cache
	result, err := cache.GetArticleList(ctx, limit, offset)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedList.Total, result.Total)
	assert.Equal(t, expectedList.Limit, result.Limit)
	assert.Equal(t, expectedList.Offset, result.Offset)
	assert.Len(t, result.Articles, 2)
	assert.Equal(t, expectedList.Articles[0].ID, result.Articles[0].ID)
	assert.Equal(t, expectedList.Articles[1].ID, result.Articles[1].ID)
}

// Test GetArticleList - Cache Miss
func TestRedisCache_GetArticleList_CacheMiss(t *testing.T) {
	cache, _, cleanup := setupRedisCache(t, 5*time.Minute)
	defer cleanup()

	ctx := context.Background()
	limit := 10
	offset := 100

	// Try to get non-existent list
	result, err := cache.GetArticleList(ctx, limit, offset)
	require.NoError(t, err)
	assert.Nil(t, result) // Cache miss should return nil, not error
}

// Test GetArticleList - Error (invalid JSON)
func TestRedisCache_GetArticleList_InvalidJSON(t *testing.T) {
	cache, mr, cleanup := setupRedisCache(t, 5*time.Minute)
	defer cleanup()

	ctx := context.Background()
	limit := 10
	offset := 0

	// Set invalid JSON in cache
	key := fmt.Sprintf("article:list:%d:%d", limit, offset)
	err := mr.Set(key, "invalid json string")
	require.NoError(t, err)

	// Try to get list - should fail on unmarshal
	result, err := cache.GetArticleList(ctx, limit, offset)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal cached list")
}

// Test SetArticleList - Success
func TestRedisCache_SetArticleList_Success(t *testing.T) {
	cache, mr, cleanup := setupRedisCache(t, 5*time.Minute)
	defer cleanup()

	ctx := context.Background()
	limit := 10
	offset := 0
	list := &dto.ListArticlesResponse{
		Articles: []dto.ArticleResponse{
			{
				ID:        1,
				Title:     "Article 1",
				Content:   "Content 1",
				AuthorID:  1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		Total:  1,
		Limit:  limit,
		Offset: offset,
	}

	// Set list in cache
	err := cache.SetArticleList(ctx, limit, offset, list)
	require.NoError(t, err)

	// Verify it was stored
	key := fmt.Sprintf("article:list:%d:%d", limit, offset)
	val, err := mr.Get(key)
	require.NoError(t, err)
	assert.NotEmpty(t, val)

	// Verify the content
	var storedList dto.ListArticlesResponse
	err = json.Unmarshal([]byte(val), &storedList)
	require.NoError(t, err)
	assert.Equal(t, list.Total, storedList.Total)
	assert.Equal(t, list.Limit, storedList.Limit)
	assert.Equal(t, list.Offset, storedList.Offset)
	assert.Len(t, storedList.Articles, 1)
}

// Test SetArticleList - Error (Redis error)
func TestRedisCache_SetArticleList_RedisError(t *testing.T) {
	cache, mr, cleanup := setupRedisCache(t, 5*time.Minute)
	defer cleanup()

	ctx := context.Background()
	limit := 10
	offset := 0
	list := &dto.ListArticlesResponse{
		Articles: []dto.ArticleResponse{},
		Total:    0,
		Limit:    limit,
		Offset:   offset,
	}

	// Close Redis to simulate error
	mr.Close()

	// Try to set list - should fail
	err := cache.SetArticleList(ctx, limit, offset, list)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to set cache")
}

// Test InvalidateArticleList - Success
func TestRedisCache_InvalidateArticleList_Success(t *testing.T) {
	cache, mr, cleanup := setupRedisCache(t, 5*time.Minute)
	defer cleanup()

	ctx := context.Background()

	// Set multiple list caches
	keys := []string{
		"article:list:10:0",
		"article:list:10:10",
		"article:list:20:0",
	}
	for _, key := range keys {
		err := mr.Set(key, "test data")
		require.NoError(t, err)
		exists := mr.Exists(key)
		assert.True(t, exists)
	}

	// Invalidate all article lists
	err := cache.InvalidateArticleList(ctx)
	require.NoError(t, err)

	// Verify all keys were deleted
	for _, key := range keys {
		exists := mr.Exists(key)
		assert.False(t, exists, "Key %s should be deleted", key)
	}
}

// Test InvalidateArticleList - Success (no keys to delete)
func TestRedisCache_InvalidateArticleList_NoKeys(t *testing.T) {
	cache, _, cleanup := setupRedisCache(t, 5*time.Minute)
	defer cleanup()

	ctx := context.Background()

	// Invalidate when no keys exist - should succeed
	err := cache.InvalidateArticleList(ctx)
	require.NoError(t, err)
}

// Test InvalidateArticleList - Error (Redis error on Keys)
func TestRedisCache_InvalidateArticleList_KeysError(t *testing.T) {
	cache, mr, cleanup := setupRedisCache(t, 5*time.Minute)
	defer cleanup()

	ctx := context.Background()

	// Close Redis to simulate error
	mr.Close()

	// Try to invalidate - should fail on Keys call
	err := cache.InvalidateArticleList(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get keys")
}

// Test InvalidateArticleList - Error (Redis error on Del)
func TestRedisCache_InvalidateArticleList_DelError(t *testing.T) {
	cache, mr, cleanup := setupRedisCache(t, 5*time.Minute)
	defer cleanup()

	ctx := context.Background()

	// Set a key first
	key := "article:list:10:0"
	err := mr.Set(key, "test data")
	require.NoError(t, err)

	// Close Redis after setting key to simulate error on Del
	mr.Close()

	// Try to invalidate - should fail on Del call
	err = cache.InvalidateArticleList(ctx)
	assert.Error(t, err)
}

// Additional tests for edge cases

func TestRedisCache_KeyGeneration_Article(t *testing.T) {
	testCases := []struct {
		name     string
		id       int64
		expected string
	}{
		{"single digit", 1, "article:1"},
		{"double digit", 10, "article:10"},
		{"large number", 12345, "article:12345"},
		{"zero", 0, "article:0"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			key := fmt.Sprintf("article:%d", tc.id)
			assert.Equal(t, tc.expected, key)
		})
	}
}

func TestRedisCache_KeyGeneration_List(t *testing.T) {
	testCases := []struct {
		name     string
		limit    int
		offset   int
		expected string
	}{
		{"default pagination", 10, 0, "article:list:10:0"},
		{"custom pagination", 20, 10, "article:list:20:10"},
		{"large numbers", 100, 50, "article:list:100:50"},
		{"zero offset", 5, 0, "article:list:5:0"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			key := fmt.Sprintf("article:list:%d:%d", tc.limit, tc.offset)
			assert.Equal(t, tc.expected, key)
		})
	}
}

func TestRedisCache_ArticleJSONMarshaling(t *testing.T) {
	article := &dto.ArticleResponse{
		ID:        1,
		Title:     "Test Article",
		Content:   "Test Content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	data, err := json.Marshal(article)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var unmarshaled dto.ArticleResponse
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, article.ID, unmarshaled.ID)
	assert.Equal(t, article.Title, unmarshaled.Title)
	assert.Equal(t, article.Content, unmarshaled.Content)
}

func TestRedisCache_ListJSONMarshaling(t *testing.T) {
	list := &dto.ListArticlesResponse{
		Articles: []dto.ArticleResponse{
			{
				ID:        1,
				Title:     "Article 1",
				Content:   "Content 1",
				AuthorID:  1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		Total:  1,
		Limit:  10,
		Offset: 0,
	}

	data, err := json.Marshal(list)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var unmarshaled dto.ListArticlesResponse
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, list.Total, unmarshaled.Total)
	assert.Equal(t, list.Limit, unmarshaled.Limit)
	assert.Equal(t, list.Offset, unmarshaled.Offset)
	assert.Len(t, unmarshaled.Articles, 1)
}

func TestRedisCache_InvalidJSONHandling(t *testing.T) {
	invalidJSON := "invalid json string"

	var article dto.ArticleResponse
	err := json.Unmarshal([]byte(invalidJSON), &article)
	assert.Error(t, err)
}

func TestRedisCache_EmptyArticleHandling(t *testing.T) {
	article := &dto.ArticleResponse{}

	data, err := json.Marshal(article)
	assert.NoError(t, err)

	var unmarshaled dto.ArticleResponse
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), unmarshaled.ID)
	assert.Empty(t, unmarshaled.Title)
}
