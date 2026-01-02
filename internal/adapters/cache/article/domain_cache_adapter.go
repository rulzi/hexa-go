package article

import (
	"context"

	"github.com/rulzi/hexa-go/internal/application/article/dto"
	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
)

// DomainCacheAdapter adapts the DTO-based cache to domain cache port
type DomainCacheAdapter struct {
	dtoCache *RedisCache
}

// NewDomainCacheAdapter creates a new domain cache adapter
func NewDomainCacheAdapter(dtoCache *RedisCache) *DomainCacheAdapter {
	return &DomainCacheAdapter{
		dtoCache: dtoCache,
	}
}

// Get implements domainarticle.Cache interface
func (a *DomainCacheAdapter) Get(ctx context.Context, id int64) (*domainarticle.Article, error) {
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

// Set implements domainarticle.Cache interface
func (a *DomainCacheAdapter) Set(ctx context.Context, id int64, article *domainarticle.Article) error {
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

// Delete implements domainarticle.Cache interface
func (a *DomainCacheAdapter) Delete(ctx context.Context, id int64) error {
	return a.dtoCache.DeleteArticle(ctx, id)
}

// InvalidateList implements domainarticle.Cache interface
func (a *DomainCacheAdapter) InvalidateList(ctx context.Context) error {
	return a.dtoCache.InvalidateArticleList(ctx)
}

