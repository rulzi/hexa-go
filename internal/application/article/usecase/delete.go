package usecase

import (
	"context"

	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
)

// DeleteArticleUseCase handles deleting an article
type DeleteArticleUseCase struct {
	articleRepo domainarticle.Repository
	cache       domainarticle.Cache
	listCache   ArticleListCache
}

// NewDeleteArticleUseCase creates a new DeleteArticleUseCase
func NewDeleteArticleUseCase(articleRepo domainarticle.Repository, cache domainarticle.Cache, listCache ArticleListCache) *DeleteArticleUseCase {
	return &DeleteArticleUseCase{
		articleRepo: articleRepo,
		cache:       cache,
		listCache:   listCache,
	}
}

// Execute executes the delete article use case
func (uc *DeleteArticleUseCase) Execute(ctx context.Context, id int64) error {
	// Check if article exists
	existingArticle, err := uc.articleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingArticle == nil {
		return domainarticle.ErrArticleNotFound
	}

	// Delete article
	if err := uc.articleRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Invalidate cache
	if uc.cache != nil {
		_ = uc.cache.Delete(ctx, id)
		_ = uc.cache.InvalidateList(ctx)
	}
	if uc.listCache != nil {
		_ = uc.listCache.InvalidateArticleList(ctx)
	}

	return nil
}
