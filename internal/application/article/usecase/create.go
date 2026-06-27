package usecase

import (
	"context"
	"time"

	"github.com/rulzi/hexa-go/internal/application/article/dto"
	articleentity "github.com/rulzi/hexa-go/internal/domain/article/entity"
	articleport "github.com/rulzi/hexa-go/internal/domain/article/port"
)

type createArticleDeps struct {
	articleRepo articleport.Repository
	cache       articleport.Cache
}

// CreateArticle handles article creation
type CreateArticle struct {
	deps createArticleDeps
}

// NewCreateArticle creates a new CreateArticle use case.
func NewCreateArticle(articleRepo articleport.Repository, cache articleport.Cache) *CreateArticle {
	return &CreateArticle{deps: createArticleDeps{
		articleRepo: articleRepo,
		cache:       cache,
	}}
}

// Execute creates a new article
func (uc *CreateArticle) Execute(ctx context.Context, req dto.CreateArticleRequest) (*dto.ArticleResponse, error) {
	newArticle := &articleentity.Article{
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  req.AuthorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := newArticle.Validate(); err != nil {
		return nil, err
	}

	createdArticle, err := uc.deps.articleRepo.Create(ctx, newArticle)
	if err != nil {
		return nil, err
	}

	if uc.deps.cache != nil {
		_ = uc.deps.cache.InvalidateList(ctx)
	}

	return &dto.ArticleResponse{
		ID:        createdArticle.ID,
		Title:     createdArticle.Title,
		Content:   createdArticle.Content,
		AuthorID:  createdArticle.AuthorID,
		CreatedAt: createdArticle.CreatedAt,
		UpdatedAt: createdArticle.UpdatedAt,
	}, nil
}
