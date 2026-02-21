package article

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rulzi/hexa-go/internal/adapters/http/response"
	"github.com/rulzi/hexa-go/internal/application/article/dto"
	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
)

// ArticleUseCase is the interface for article operations
type ArticleUseCase interface {
	Create(ctx context.Context, req dto.CreateArticleRequest) (*dto.ArticleResponse, error)
	Get(ctx context.Context, id int64) (*dto.ArticleResponse, error)
	List(ctx context.Context, limit, offset int) (*dto.ListArticlesResponse, error)
	Update(ctx context.Context, id int64, req dto.UpdateArticleRequest) (*dto.ArticleResponse, error)
	Delete(ctx context.Context, id int64) error
}

// Handler handles HTTP requests for articles
type Handler struct {
	articleUseCase ArticleUseCase
}

// NewHandler creates a new Handler
func NewHandler(articleUseCase ArticleUseCase) *Handler {
	return &Handler{
		articleUseCase: articleUseCase,
	}
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
		response.ErrorResponseInternalServerError(c, err.Error())
		return
	}

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
		if err == domainarticle.ErrArticleNotFound {
			response.ErrorResponseNotFound(c, err.Error())
		} else {
			response.ErrorResponseInternalServerError(c, err.Error())
		}
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
		response.ErrorResponseInternalServerError(c, err.Error())
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
		if err == domainarticle.ErrArticleNotFound {
			response.ErrorResponseNotFound(c, err.Error())
		} else {
			response.ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

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
		if err == domainarticle.ErrArticleNotFound {
			response.ErrorResponseNotFound(c, err.Error())
		} else {
			response.ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	response.SuccessResponseOK(c, "Article deleted successfully", nil)
}
