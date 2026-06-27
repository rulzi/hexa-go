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

	cached, err := uc.deps.listCache.GetList(ctx, limit, offset)
	if err == nil && cached != nil {
		return listPageToResponse(cached, limit, offset), nil
	}

	articles, err := uc.deps.articleRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.deps.articleRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	response := listPageToResponse(&ArticleListPage{Articles: articles, Total: total}, limit, offset)
	_ = uc.deps.listCache.SetList(ctx, limit, offset, &ArticleListPage{Articles: articles, Total: total})

	return response, nil
}

func listPageToResponse(page *ArticleListPage, limit, offset int) *dto.ListArticlesResponse {
	articleResponses := make([]dto.ArticleResponse, len(page.Articles))
	for i, a := range page.Articles {
		articleResponses[i] = dto.ArticleResponse{
			ID:        a.ID,
			Title:     a.Title,
			Content:   a.Content,
			AuthorID:  a.AuthorID,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		}
	}

	return &dto.ListArticlesResponse{
		Articles: articleResponses,
		Total:    page.Total,
		Limit:    limit,
		Offset:   offset,
	}
}
