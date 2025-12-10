package article

import (
	"context"

	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
)

// DeleteArticleUseCase handles deleting an article
type DeleteArticleUseCase struct {
	articleRepo domainarticle.Repository
	singleCache ArticleSingleCache
	listCache   ArticleCache
}

// NewDeleteArticleUseCase creates a new DeleteArticleUseCase
func NewDeleteArticleUseCase(articleRepo domainarticle.Repository, singleCache ArticleSingleCache, listCache ArticleCache) *DeleteArticleUseCase {
	return &DeleteArticleUseCase{
		articleRepo: articleRepo,
		singleCache: singleCache,
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
	if uc.singleCache != nil {
		_ = uc.singleCache.DeleteArticle(ctx, id)
	}
	if uc.listCache != nil {
		_ = uc.listCache.InvalidateArticleList(ctx)
	}

	return nil
}

