package usecase

import (
	"context"
	"time"

	"github.com/rulzi/hexa-go/internal/application/article/dto"
	articleentity "github.com/rulzi/hexa-go/internal/domain/article/entity"
	articleport "github.com/rulzi/hexa-go/internal/domain/article/port"
	articleservice "github.com/rulzi/hexa-go/internal/domain/article/service"
)

// ArticleListCache defines the interface for article list caching (DTO-based for performance)
type ArticleListCache interface {
	GetArticleList(ctx context.Context, limit, offset int) (*dto.ListArticlesResponse, error)
	SetArticleList(ctx context.Context, limit, offset int, listResp *dto.ListArticlesResponse) error
	InvalidateArticleList(ctx context.Context) error
}

// ArticleUseCase handles all article operations (create, get, list, update, delete)
type ArticleUseCase struct {
	articleRepo    articleport.Repository
	articleService *articleservice.Service
	cache          articleport.Cache
	listCache      ArticleListCache
}

// NewArticleUseCase creates a new ArticleUseCase
func NewArticleUseCase(
	articleRepo articleport.Repository,
	articleService *articleservice.Service,
	cache articleport.Cache,
	listCache ArticleListCache,
) *ArticleUseCase {
	return &ArticleUseCase{
		articleRepo:    articleRepo,
		articleService: articleService,
		cache:          cache,
		listCache:      listCache,
	}
}

// Create creates a new article
func (uc *ArticleUseCase) Create(ctx context.Context, req dto.CreateArticleRequest) (*dto.ArticleResponse, error) {
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

	createdArticle, err := uc.articleRepo.Create(ctx, newArticle)
	if err != nil {
		return nil, err
	}

	if uc.cache != nil {
		_ = uc.cache.InvalidateList(ctx)
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

// Get retrieves an article by ID
func (uc *ArticleUseCase) Get(ctx context.Context, id int64) (*dto.ArticleResponse, error) {
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

	articleEntity, err := uc.articleRepo.GetByID(ctx, id)
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

	if uc.cache != nil {
		_ = uc.cache.Set(ctx, id, articleEntity)
	}

	return response, nil
}

// List lists articles with pagination
func (uc *ArticleUseCase) List(ctx context.Context, limit, offset int) (*dto.ListArticlesResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	if uc.listCache != nil {
		cached, err := uc.listCache.GetArticleList(ctx, limit, offset)
		if err == nil && cached != nil {
			return cached, nil
		}
	}

	articles, err := uc.articleRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.articleRepo.Count(ctx)
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

	if uc.listCache != nil {
		_ = uc.listCache.SetArticleList(ctx, limit, offset, response)
	}

	return response, nil
}

// Update updates an article
func (uc *ArticleUseCase) Update(ctx context.Context, id int64, req dto.UpdateArticleRequest) (*dto.ArticleResponse, error) {
	existingArticle, err := uc.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingArticle == nil {
		return nil, articleentity.NewArticleNotFound()
	}

	existingArticle.Title = req.Title
	existingArticle.Content = req.Content
	existingArticle.UpdatedAt = time.Now()

	if err := existingArticle.Validate(); err != nil {
		return nil, err
	}

	updatedArticle, err := uc.articleRepo.Update(ctx, existingArticle)
	if err != nil {
		return nil, err
	}

	response := &dto.ArticleResponse{
		ID:        updatedArticle.ID,
		Title:     updatedArticle.Title,
		Content:   updatedArticle.Content,
		AuthorID:  updatedArticle.AuthorID,
		CreatedAt: updatedArticle.CreatedAt,
		UpdatedAt: updatedArticle.UpdatedAt,
	}

	if uc.cache != nil {
		_ = uc.cache.Delete(ctx, id)
		_ = uc.cache.InvalidateList(ctx)
	}
	if uc.listCache != nil {
		_ = uc.listCache.InvalidateArticleList(ctx)
	}

	return response, nil
}

// Delete deletes an article
func (uc *ArticleUseCase) Delete(ctx context.Context, id int64) error {
	existingArticle, err := uc.articleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingArticle == nil {
		return articleentity.NewArticleNotFound()
	}

	if err := uc.articleRepo.Delete(ctx, id); err != nil {
		return err
	}

	if uc.cache != nil {
		_ = uc.cache.Delete(ctx, id)
		_ = uc.cache.InvalidateList(ctx)
	}
	if uc.listCache != nil {
		_ = uc.listCache.InvalidateArticleList(ctx)
	}

	return nil
}
