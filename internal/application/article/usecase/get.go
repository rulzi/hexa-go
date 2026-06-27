package usecase

import (
	"context"

	"github.com/rulzi/hexa-go/internal/application/article/dto"
	articleentity "github.com/rulzi/hexa-go/internal/domain/article/entity"
)

// GetArticle handles article retrieval by ID
type GetArticle struct {
	deps articleDeps
}

func newGetArticle(deps articleDeps) *GetArticle {
	return &GetArticle{deps: deps}
}

// Execute retrieves an article by ID
func (uc *GetArticle) Execute(ctx context.Context, id int64) (*dto.ArticleResponse, error) {
	if uc.deps.cache != nil {
		cached, err := uc.deps.cache.Get(ctx, id)
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

	articleEntity, err := uc.deps.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if articleEntity == nil {
		return nil, articleentity.NewArticleNotFound()
	}

	response := &dto.ArticleResponse{
		ID:        articleEntity.ID,
		Title:     articleEntity.Title,
		Content:   articleEntity.Content,
		AuthorID:  articleEntity.AuthorID,
		CreatedAt: articleEntity.CreatedAt,
		UpdatedAt: articleEntity.UpdatedAt,
	}

	if uc.deps.cache != nil {
		_ = uc.deps.cache.Set(ctx, id, articleEntity)
	}

	return response, nil
}
