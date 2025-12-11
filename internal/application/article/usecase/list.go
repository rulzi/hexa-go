package usecase

import (
	"context"

	"github.com/rulzi/hexa-go/internal/application/article/dto"
	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
)

// ArticleCache defines the interface for article caching
type ArticleCache interface {
	GetArticleList(ctx context.Context, limit, offset int) (*dto.ListArticlesResponse, error)
	SetArticleList(ctx context.Context, limit, offset int, listResp *dto.ListArticlesResponse) error
	InvalidateArticleList(ctx context.Context) error
}

// ListArticlesUseCase handles listing articles with pagination
type ListArticlesUseCase struct {
	articleRepo domainarticle.Repository
	cache       ArticleCache
}

// NewListArticlesUseCase creates a new ListArticlesUseCase
func NewListArticlesUseCase(articleRepo domainarticle.Repository, cache ArticleCache) *ListArticlesUseCase {
	return &ListArticlesUseCase{
		articleRepo: articleRepo,
		cache:       cache,
	}
}

// Execute executes the list articles use case
func (uc *ListArticlesUseCase) Execute(ctx context.Context, limit, offset int) (*dto.ListArticlesResponse, error) {
	// Default pagination
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Try to get from cache first
	if uc.cache != nil {
		cached, err := uc.cache.GetArticleList(ctx, limit, offset)
		if err == nil && cached != nil {
			return cached, nil
		}
	}

	// Get articles from repository
	articles, err := uc.articleRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := uc.articleRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	articleResponses := make([]dto.ArticleResponse, len(articles))
	for i, a := range articles {
		articleResponses[i] = dto.ArticleResponse{
			ID:        a.ID,
			Title:     a.Title,
			Content:   a.Content,
			AuthorID:  a.AuthorID,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		}
	}

	response := &dto.ListArticlesResponse{
		Articles: articleResponses,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}

	// Store in cache
	if uc.cache != nil {
		_ = uc.cache.SetArticleList(ctx, limit, offset, response)
	}

	return response, nil
}
