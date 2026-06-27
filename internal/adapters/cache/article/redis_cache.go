package article

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	articleusecase "github.com/rulzi/hexa-go/internal/application/article/usecase"
	articleentity "github.com/rulzi/hexa-go/internal/domain/article/entity"
)

var (
	_ articleusecase.ArticleCache     = (*RedisCache)(nil)
	_ articleusecase.ArticleListCache = (*RedisCache)(nil)
)

// RedisCache handles caching for articles using Redis.
type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisCache creates a new RedisCache.
func NewRedisCache(client *redis.Client, ttl time.Duration) *RedisCache {
	if ttl == 0 {
		ttl = 5 * time.Minute
	}
	return &RedisCache{
		client: client,
		ttl:    ttl,
	}
}

// Get retrieves an article entity from cache by ID.
func (c *RedisCache) Get(ctx context.Context, id int64) (*articleentity.Article, error) {
	key := fmt.Sprintf("article:%d", id)

	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}

	var article articleentity.Article
	if err := json.Unmarshal([]byte(val), &article); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached article: %w", err)
	}

	return &article, nil
}

// Set stores an article entity in cache.
func (c *RedisCache) Set(ctx context.Context, id int64, article *articleentity.Article) error {
	key := fmt.Sprintf("article:%d", id)

	data, err := json.Marshal(article)
	if err != nil {
		return fmt.Errorf("failed to marshal article: %w", err)
	}

	if err := c.client.Set(ctx, key, data, c.ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// Delete removes an article from cache.
func (c *RedisCache) Delete(ctx context.Context, id int64) error {
	key := fmt.Sprintf("article:%d", id)
	return c.client.Del(ctx, key).Err()
}

// InvalidateList invalidates all article list caches.
func (c *RedisCache) InvalidateList(ctx context.Context) error {
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

// GetList retrieves a paginated list of article entities from cache.
func (c *RedisCache) GetList(ctx context.Context, limit, offset int) (*articleusecase.ArticleListPage, error) {
	key := fmt.Sprintf("article:list:%d:%d", limit, offset)

	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}

	var page articleusecase.ArticleListPage
	if err := json.Unmarshal([]byte(val), &page); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached list: %w", err)
	}

	return &page, nil
}

// SetList stores a paginated list of article entities in cache.
func (c *RedisCache) SetList(ctx context.Context, limit, offset int, page *articleusecase.ArticleListPage) error {
	key := fmt.Sprintf("article:list:%d:%d", limit, offset)

	data, err := json.Marshal(page)
	if err != nil {
		return fmt.Errorf("failed to marshal article list: %w", err)
	}

	if err := c.client.Set(ctx, key, data, c.ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}
