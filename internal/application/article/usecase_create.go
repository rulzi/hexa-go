package article

import (
	"context"
	"time"

	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
)

// CreateArticleUseCase handles the creation of a new article
type CreateArticleUseCase struct {
	articleRepo    domainarticle.Repository
	articleService *domainarticle.Service
	cache          ArticleCache
}

// NewCreateArticleUseCase creates a new CreateArticleUseCase
func NewCreateArticleUseCase(
	articleRepo domainarticle.Repository,
	articleService *domainarticle.Service,
	cache ArticleCache,
) *CreateArticleUseCase {
	return &CreateArticleUseCase{
		articleRepo:    articleRepo,
		articleService: articleService,
		cache:          cache,
	}
}

// Execute executes the create article use case
func (uc *CreateArticleUseCase) Execute(ctx context.Context, req CreateArticleRequest) (*ArticleResponse, error) {
	// Create article entity
	newArticle := &domainarticle.Article{
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  req.AuthorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Validate entity
	if err := newArticle.Validate(); err != nil {
		return nil, err
	}

	// Save to repository
	createdArticle, err := uc.articleRepo.Create(ctx, newArticle)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	if uc.cache != nil {
		_ = uc.cache.InvalidateArticleList(ctx)
	}

	// Return response DTO
	return &ArticleResponse{
		ID:        createdArticle.ID,
		Title:     createdArticle.Title,
		Content:   createdArticle.Content,
		AuthorID:  createdArticle.AuthorID,
		CreatedAt: createdArticle.CreatedAt,
		UpdatedAt: createdArticle.UpdatedAt,
	}, nil
}
