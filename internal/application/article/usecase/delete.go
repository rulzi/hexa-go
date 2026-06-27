package usecase

import (
	"context"

	articleentity "github.com/rulzi/hexa-go/internal/domain/article/entity"
	articleport "github.com/rulzi/hexa-go/internal/domain/article/port"
)

type deleteArticleDeps struct {
	articleRepo articleport.Repository
	cache       articleport.Cache
	listCache   ArticleListCache
}

// DeleteArticle handles article deletion
type DeleteArticle struct {
	deps deleteArticleDeps
}

// NewDeleteArticle creates a new DeleteArticle use case.
func NewDeleteArticle(articleRepo articleport.Repository, cache articleport.Cache, listCache ArticleListCache) *DeleteArticle {
	return &DeleteArticle{deps: deleteArticleDeps{
		articleRepo: articleRepo,
		cache:       cache,
		listCache:   listCache,
	}}
}

// Execute deletes an article
func (uc *DeleteArticle) Execute(ctx context.Context, id int64) error {
	existingArticle, err := uc.deps.articleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingArticle == nil {
		return articleentity.NewArticleNotFound()
	}

	if err := uc.deps.articleRepo.Delete(ctx, id); err != nil {
		return err
	}

	if uc.deps.cache != nil {
		_ = uc.deps.cache.Delete(ctx, id)
		_ = uc.deps.cache.InvalidateList(ctx)
	}
	if uc.deps.listCache != nil {
		_ = uc.deps.listCache.InvalidateArticleList(ctx)
	}

	return nil
}
