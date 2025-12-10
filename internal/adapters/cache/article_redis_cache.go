package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rulzi/hexa-go/internal/application/article"
)

// ArticleRedisCache handles caching for articles using Redis
type ArticleRedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewArticleRedisCache creates a new ArticleRedisCache
func NewArticleRedisCache(client *redis.Client, ttl time.Duration) *ArticleRedisCache {
	if ttl == 0 {
		ttl = 5 * time.Minute // Default TTL: 5 minutes
	}
	return &ArticleRedisCache{
		client: client,
		ttl:    ttl,
	}
}

// GetArticle retrieves an article from cache by ID
func (c *ArticleRedisCache) GetArticle(ctx context.Context, id int64) (*article.ArticleResponse, error) {
	key := fmt.Sprintf("article:%d", id)

	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}

	var articleResp article.ArticleResponse
	if err := json.Unmarshal([]byte(val), &articleResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached article: %w", err)
	}

	return &articleResp, nil
}

// SetArticle stores an article in cache
func (c *ArticleRedisCache) SetArticle(ctx context.Context, id int64, articleResp *article.ArticleResponse) error {
	key := fmt.Sprintf("article:%d", id)

	data, err := json.Marshal(articleResp)
	if err != nil {
		return fmt.Errorf("failed to marshal article: %w", err)
	}

	if err := c.client.Set(ctx, key, data, c.ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// DeleteArticle removes an article from cache
func (c *ArticleRedisCache) DeleteArticle(ctx context.Context, id int64) error {
	key := fmt.Sprintf("article:%d", id)
	return c.client.Del(ctx, key).Err()
}

// GetArticleList retrieves a list of articles from cache
func (c *ArticleRedisCache) GetArticleList(ctx context.Context, limit, offset int) (*article.ListArticlesResponse, error) {
	key := fmt.Sprintf("article:list:%d:%d", limit, offset)

	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}

	var listResp article.ListArticlesResponse
	if err := json.Unmarshal([]byte(val), &listResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached list: %w", err)
	}

	return &listResp, nil
}

// SetArticleList stores a list of articles in cache
func (c *ArticleRedisCache) SetArticleList(ctx context.Context, limit, offset int, listResp *article.ListArticlesResponse) error {
	key := fmt.Sprintf("article:list:%d:%d", limit, offset)

	data, err := json.Marshal(listResp)
	if err != nil {
		return fmt.Errorf("failed to marshal article list: %w", err)
	}

	if err := c.client.Set(ctx, key, data, c.ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// InvalidateArticleList invalidates all article list caches
func (c *ArticleRedisCache) InvalidateArticleList(ctx context.Context) error {
	pattern := "article:list:*"
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys: %w", err)
	}

	if len(keys) > 0 {
		if err := c.client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("failed to delete keys: %w", err)
		}
	}

	return nil
}
