package usecase

import (
	"context"

	articleentity "github.com/rulzi/hexa-go/internal/domain/article/entity"
)

// DeleteArticle handles article deletion
type DeleteArticle struct {
	deps articleDeps
}

func newDeleteArticle(deps articleDeps) *DeleteArticle {
	return &DeleteArticle{deps: deps}
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
