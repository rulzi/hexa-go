package article

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rulzi/hexa-go/internal/adapters/http/errmapper"
	"github.com/rulzi/hexa-go/internal/adapters/http/response"
	"github.com/rulzi/hexa-go/internal/application/article/dto"
	domainerrs "github.com/rulzi/hexa-go/internal/domain/errs"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
)

// UseCase is the interface for article operations
type UseCase interface {
	Create(ctx context.Context, req dto.CreateArticleRequest) (*dto.ArticleResponse, error)
	Get(ctx context.Context, id int64) (*dto.ArticleResponse, error)
	List(ctx context.Context, limit, offset int) (*dto.ListArticlesResponse, error)
	Update(ctx context.Context, id int64, req dto.UpdateArticleRequest) (*dto.ArticleResponse, error)
	Delete(ctx context.Context, id int64) error
}

// Handler handles HTTP requests for articles
type Handler struct {
	articleUseCase UseCase
	logger         logger.Logger
}

// NewHandler creates a new Handler
func NewHandler(articleUseCase UseCase, appLogger logger.Logger) *Handler {
	return &Handler{
		articleUseCase: articleUseCase,
		logger:         appLogger,
	}
}

func (h *Handler) handleUseCaseError(c *gin.Context, err error, fields map[string]interface{}, logMsg string) {
	if fields == nil {
		fields = make(map[string]interface{})
	}

	if domainerrs.IsNotFound(err) || domainerrs.IsValidation(err) || domainerrs.IsConflict(err) || domainerrs.IsUnauthorized(err) {
		h.logger.DebugWithFields(logMsg, fields)
	} else {
		fields["error"] = err.Error()
		h.logger.ErrorWithFields(logMsg, fields)
	}
	errmapper.Respond(c, err)
}

// Create handles POST /articles
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponseBadRequest(c, err.Error())
		return
	}

	resp, err := h.articleUseCase.Create(c.Request.Context(), req)
	if err != nil {
		h.handleUseCaseError(c, err, map[string]interface{}{"title": req.Title}, "article create failed")
		return
	}

	h.logger.InfoWithFields("article created", map[string]interface{}{"article_id": resp.ID, "title": resp.Title})
	response.SuccessResponseCreated(c, "Article created successfully", resp)
}

// Get handles GET /articles/:id
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorResponseBadRequest(c, "invalid article id")
		return
	}

	resp, err := h.articleUseCase.Get(c.Request.Context(), id)
	if err != nil {
		h.handleUseCaseError(c, err, map[string]interface{}{"article_id": id}, "article get failed")
		return
	}

	response.SuccessResponseOK(c, "Article retrieved successfully", resp)
}

// List handles GET /articles
func (h *Handler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := h.articleUseCase.List(c.Request.Context(), limit, offset)
	if err != nil {
		h.handleUseCaseError(c, err, nil, "article list failed")
		return
	}

	response.SuccessResponseOK(c, "Articles retrieved successfully", resp)
}

// Update handles PUT /articles/:id
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorResponseBadRequest(c, "invalid article id")
		return
	}

	var req dto.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponseBadRequest(c, err.Error())
		return
	}

	resp, err := h.articleUseCase.Update(c.Request.Context(), id, req)
	if err != nil {
		h.handleUseCaseError(c, err, map[string]interface{}{"article_id": id}, "article update failed")
		return
	}

	h.logger.InfoWithFields("article updated", map[string]interface{}{"article_id": id})
	response.SuccessResponseOK(c, "Article updated successfully", resp)
}

// Delete handles DELETE /articles/:id
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorResponseBadRequest(c, "invalid article id")
		return
	}

	err = h.articleUseCase.Delete(c.Request.Context(), id)
	if err != nil {
		h.handleUseCaseError(c, err, map[string]interface{}{"article_id": id}, "article delete failed")
		return
	}

	h.logger.InfoWithFields("article deleted", map[string]interface{}{"article_id": id})
	response.SuccessResponseOK(c, "Article deleted successfully", nil)
}
