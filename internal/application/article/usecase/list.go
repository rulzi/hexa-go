package usecase

import (
	"context"

	"github.com/rulzi/hexa-go/internal/application/article/dto"
	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
)

// ListArticlesUseCase handles listing articles with pagination
type ListArticlesUseCase struct {
	articleRepo domainarticle.Repository
	cache       domainarticle.Cache
	dtoCache    ArticleListCache // Keep DTO cache for list caching (performance optimization)
}

// ArticleListCache defines the interface for article list caching (DTO-based for performance)
// This is a secondary adapter interface for list caching
type ArticleListCache interface {
	GetArticleList(ctx context.Context, limit, offset int) (*dto.ListArticlesResponse, error)
	SetArticleList(ctx context.Context, limit, offset int, listResp *dto.ListArticlesResponse) error
	InvalidateArticleList(ctx context.Context) error
}

// NewListArticlesUseCase creates a new ListArticlesUseCase
func NewListArticlesUseCase(articleRepo domainarticle.Repository, cache domainarticle.Cache, dtoCache ArticleListCache) *ListArticlesUseCase {
	return &ListArticlesUseCase{
		articleRepo: articleRepo,
		cache:       cache,
		dtoCache:    dtoCache,
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

	// Try to get from cache first (using DTO cache for performance)
	if uc.dtoCache != nil {
		cached, err := uc.dtoCache.GetArticleList(ctx, limit, offset)
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

	// Store in cache (using DTO cache for performance)
	if uc.dtoCache != nil {
		_ = uc.dtoCache.SetArticleList(ctx, limit, offset, response)
	}

	return response, nil
}
