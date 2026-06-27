package usecase

import (
	"context"

	"github.com/rulzi/hexa-go/internal/application/article/dto"
	articleport "github.com/rulzi/hexa-go/internal/domain/article/port"
)

type listArticleDeps struct {
	articleRepo articleport.Repository
	listCache   ArticleListCache
}

// ListArticle handles article listing with pagination
type ListArticle struct {
	deps listArticleDeps
}

// NewListArticle creates a new ListArticle use case.
func NewListArticle(articleRepo articleport.Repository, listCache ArticleListCache) *ListArticle {
	return &ListArticle{deps: listArticleDeps{
		articleRepo: articleRepo,
		listCache:   listCache,
	}}
}

// Execute lists articles with pagination
func (uc *ListArticle) Execute(ctx context.Context, limit, offset int) (*dto.ListArticlesResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	if uc.deps.listCache != nil {
		cached, err := uc.deps.listCache.GetArticleList(ctx, limit, offset)
		if err == nil && cached != nil {
			return cached, nil
		}
	}

	articles, err := uc.deps.articleRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.deps.articleRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

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

	if uc.deps.listCache != nil {
		_ = uc.deps.listCache.SetArticleList(ctx, limit, offset, response)
	}

	return response, nil
}
