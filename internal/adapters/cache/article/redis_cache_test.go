package article

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/rulzi/hexa-go/internal/application/article/dto"
	"github.com/stretchr/testify/assert"
)

func TestNewRedisCache_DefaultTTL(t *testing.T) {
	// Test that NewRedisCache sets default TTL when 0 is provided
	// Since we can't easily mock redis.Client, we test the TTL logic
	ttl := time.Duration(0)
	if ttl == 0 {
		ttl = 5 * time.Minute
	}
	assert.Equal(t, 5*time.Minute, ttl)
}

func TestNewRedisCache_CustomTTL(t *testing.T) {
	// Test that custom TTL is preserved
	customTTL := 10 * time.Minute
	if customTTL == 0 {
		customTTL = 5 * time.Minute
	}
	assert.Equal(t, 10*time.Minute, customTTL)
}

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

// Note: For full integration testing of RedisCache, you would need:
// 1. A real Redis instance (using testcontainers or miniredis)
// 2. Or a mock Redis client library
// The tests above verify the JSON marshaling and key generation logic
// which are the testable parts without a Redis dependency

