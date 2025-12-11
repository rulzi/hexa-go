package usecase

import (
	"context"

	"github.com/rulzi/hexa-go/internal/application/article/dto"
	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
)

// ArticleSingleCache defines the interface for single article caching
type ArticleSingleCache interface {
	GetArticle(ctx context.Context, id int64) (*dto.ArticleResponse, error)
	SetArticle(ctx context.Context, id int64, articleResp *dto.ArticleResponse) error
	DeleteArticle(ctx context.Context, id int64) error
}

// GetArticleUseCase handles retrieving an article by ID
type GetArticleUseCase struct {
	articleRepo domainarticle.Repository
	cache       ArticleSingleCache
}

// NewGetArticleUseCase creates a new GetArticleUseCase
func NewGetArticleUseCase(articleRepo domainarticle.Repository, cache ArticleSingleCache) *GetArticleUseCase {
	return &GetArticleUseCase{
		articleRepo: articleRepo,
		cache:       cache,
	}
}

// Execute executes the get article use case
func (uc *GetArticleUseCase) Execute(ctx context.Context, id int64) (*dto.ArticleResponse, error) {
	// Try to get from cache first
	if uc.cache != nil {
		cached, err := uc.cache.GetArticle(ctx, id)
		if err == nil && cached != nil {
			return cached, nil
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
		_ = uc.cache.SetArticle(ctx, id, response)
	}

	return response, nil
}
