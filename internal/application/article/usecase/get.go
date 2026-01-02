package usecase

import (
	"context"

	"github.com/rulzi/hexa-go/internal/application/article/dto"
	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
)

// GetArticleUseCase handles retrieving an article by ID
type GetArticleUseCase struct {
	articleRepo domainarticle.Repository
	cache       domainarticle.Cache
}

// NewGetArticleUseCase creates a new GetArticleUseCase
func NewGetArticleUseCase(articleRepo domainarticle.Repository, cache domainarticle.Cache) *GetArticleUseCase {
	return &GetArticleUseCase{
		articleRepo: articleRepo,
		cache:       cache,
	}
}

// Execute executes the get article use case
func (uc *GetArticleUseCase) Execute(ctx context.Context, id int64) (*dto.ArticleResponse, error) {
	// Try to get from cache first
	if uc.cache != nil {
		cached, err := uc.cache.Get(ctx, id)
		if err == nil && cached != nil {
			return &dto.ArticleResponse{
				ID:        cached.ID,
				Title:     cached.Title,
				Content:   cached.Content,
				AuthorID:  cached.AuthorID,
				CreatedAt: cached.CreatedAt,
				UpdatedAt: cached.UpdatedAt,
			}, nil
		}
	}

	// Get from repository
	articleEntity, err := uc.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if articleEntity == nil {
		return nil, domainarticle.ErrArticleNotFound
	}

	response := &dto.ArticleResponse{
		ID:        articleEntity.ID,
		Title:     articleEntity.Title,
		Content:   articleEntity.Content,
		AuthorID:  articleEntity.AuthorID,
		CreatedAt: articleEntity.CreatedAt,
		UpdatedAt: articleEntity.UpdatedAt,
	}

	// Store in cache
	if uc.cache != nil {
		_ = uc.cache.Set(ctx, id, articleEntity)
	}

	return response, nil
}
