package usecase

import (
	"context"
	"time"

	"github.com/rulzi/hexa-go/internal/application/article/dto"
	articleentity "github.com/rulzi/hexa-go/internal/domain/article/entity"
	articleport "github.com/rulzi/hexa-go/internal/domain/article/port"
)

type updateArticleDeps struct {
	articleRepo articleport.Repository
	cache       articleport.Cache
	listCache   ArticleListCache
}

// UpdateArticle handles article updates
type UpdateArticle struct {
	deps updateArticleDeps
}

// NewUpdateArticle creates a new UpdateArticle use case.
func NewUpdateArticle(articleRepo articleport.Repository, cache articleport.Cache, listCache ArticleListCache) *UpdateArticle {
	return &UpdateArticle{deps: updateArticleDeps{
		articleRepo: articleRepo,
		cache:       cache,
		listCache:   listCache,
	}}
}

// Execute updates an article
func (uc *UpdateArticle) Execute(ctx context.Context, id int64, req dto.UpdateArticleRequest) (*dto.ArticleResponse, error) {
	existingArticle, err := uc.deps.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingArticle == nil {
		return nil, articleentity.NewArticleNotFound()
	}

	existingArticle.Title = req.Title
	existingArticle.Content = req.Content
	existingArticle.UpdatedAt = time.Now()

	if err := existingArticle.Validate(); err != nil {
		return nil, err
	}

	updatedArticle, err := uc.deps.articleRepo.Update(ctx, existingArticle)
	if err != nil {
		return nil, err
	}

	response := &dto.ArticleResponse{
		ID:        updatedArticle.ID,
		Title:     updatedArticle.Title,
		Content:   updatedArticle.Content,
		AuthorID:  updatedArticle.AuthorID,
		CreatedAt: updatedArticle.CreatedAt,
		UpdatedAt: updatedArticle.UpdatedAt,
	}

	if uc.deps.cache != nil {
		_ = uc.deps.cache.Delete(ctx, id)
		_ = uc.deps.cache.InvalidateList(ctx)
	}
	if uc.deps.listCache != nil {
		_ = uc.deps.listCache.InvalidateArticleList(ctx)
	}

	return response, nil
}
