package article

import (
	"context"
	"time"

	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
)

// UpdateArticleUseCase handles updating an article
type UpdateArticleUseCase struct {
	articleRepo    domainarticle.Repository
	articleService *domainarticle.Service
	singleCache    ArticleSingleCache
	listCache      ArticleCache
}

// NewUpdateArticleUseCase creates a new UpdateArticleUseCase
func NewUpdateArticleUseCase(
	articleRepo domainarticle.Repository,
	articleService *domainarticle.Service,
	singleCache ArticleSingleCache,
	listCache ArticleCache,
) *UpdateArticleUseCase {
	return &UpdateArticleUseCase{
		articleRepo:    articleRepo,
		articleService: articleService,
		singleCache:    singleCache,
		listCache:      listCache,
	}
}

// Execute executes the update article use case
func (uc *UpdateArticleUseCase) Execute(ctx context.Context, id int64, req UpdateArticleRequest) (*ArticleResponse, error) {
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

	response := &ArticleResponse{
		ID:        updatedArticle.ID,
		Title:     updatedArticle.Title,
		Content:   updatedArticle.Content,
		AuthorID:  updatedArticle.AuthorID,
		CreatedAt: updatedArticle.CreatedAt,
		UpdatedAt: updatedArticle.UpdatedAt,
	}

	// Invalidate cache
	if uc.singleCache != nil {
		_ = uc.singleCache.DeleteArticle(ctx, id)
	}
	if uc.listCache != nil {
		_ = uc.listCache.InvalidateArticleList(ctx)
	}

	return response, nil
}

