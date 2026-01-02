package usecase

import (
	"context"
	"time"

	"github.com/rulzi/hexa-go/internal/application/article/dto"
	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
)

// UpdateArticleUseCase handles updating an article
type UpdateArticleUseCase struct {
	articleRepo    domainarticle.Repository
	articleService *domainarticle.Service
	cache          domainarticle.Cache
	listCache      ArticleListCache
}

// NewUpdateArticleUseCase creates a new UpdateArticleUseCase
func NewUpdateArticleUseCase(
	articleRepo domainarticle.Repository,
	articleService *domainarticle.Service,
	cache domainarticle.Cache,
	listCache ArticleListCache,
) *UpdateArticleUseCase {
	return &UpdateArticleUseCase{
		articleRepo:    articleRepo,
		articleService: articleService,
		cache:          cache,
		listCache:      listCache,
	}
}

// Execute executes the update article use case
func (uc *UpdateArticleUseCase) Execute(ctx context.Context, id int64, req dto.UpdateArticleRequest) (*dto.ArticleResponse, error) {
	// Get existing article
	existingArticle, err := uc.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingArticle == nil {
		return nil, domainarticle.ErrArticleNotFound
	}

	// Update fields
	existingArticle.Title = req.Title
	existingArticle.Content = req.Content
	existingArticle.UpdatedAt = time.Now()

	// Validate entity
	if err := existingArticle.Validate(); err != nil {
		return nil, err
	}

	// Update in repository
	updatedArticle, err := uc.articleRepo.Update(ctx, existingArticle)
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

	// Invalidate cache
	if uc.cache != nil {
		_ = uc.cache.Delete(ctx, id)
		_ = uc.cache.InvalidateList(ctx)
	}
	if uc.listCache != nil {
		_ = uc.listCache.InvalidateArticleList(ctx)
	}

	return response, nil
}
