package article

import (
	"context"

	"github.com/rulzi/hexa-go/internal/application/article/dto"
	articleentity "github.com/rulzi/hexa-go/internal/domain/article/entity"
	articleport "github.com/rulzi/hexa-go/internal/domain/article/port"
)

var _ articleport.Cache = (*DomainCacheAdapter)(nil)

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

// Get implements articleport.Cache interface
func (a *DomainCacheAdapter) Get(ctx context.Context, id int64) (*articleentity.Article, error) {
	dtoResp, err := a.dtoCache.GetArticle(ctx, id)
	if err != nil {
		return nil, err
	}
	if dtoResp == nil {
		return nil, nil
	}

	return &articleentity.Article{
		ID:        dtoResp.ID,
		Title:     dtoResp.Title,
		Content:   dtoResp.Content,
		AuthorID:  dtoResp.AuthorID,
		CreatedAt: dtoResp.CreatedAt,
		UpdatedAt: dtoResp.UpdatedAt,
	}, nil
}

// Set implements articleport.Cache interface
func (a *DomainCacheAdapter) Set(ctx context.Context, id int64, article *articleentity.Article) error {
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

// Delete implements articleport.Cache interface
func (a *DomainCacheAdapter) Delete(ctx context.Context, id int64) error {
	return a.dtoCache.DeleteArticle(ctx, id)
}

// InvalidateList implements articleport.Cache interface
func (a *DomainCacheAdapter) InvalidateList(ctx context.Context) error {
	return a.dtoCache.InvalidateArticleList(ctx)
}

