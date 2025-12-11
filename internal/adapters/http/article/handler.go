package article

import (
	"strconv"

	"github.com/gin-gonic/gin"
	apparticle "github.com/rulzi/hexa-go/internal/application/article"
	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
	"github.com/rulzi/hexa-go/internal/adapters/http/response"
)

// Handler handles HTTP requests for articles
type Handler struct {
	createUseCase *apparticle.CreateArticleUseCase
	getUseCase    *apparticle.GetArticleUseCase
	listUseCase   *apparticle.ListArticlesUseCase
	updateUseCase *apparticle.UpdateArticleUseCase
	deleteUseCase *apparticle.DeleteArticleUseCase
}

// NewHandler creates a new Handler
func NewHandler(
	createUseCase *apparticle.CreateArticleUseCase,
	getUseCase *apparticle.GetArticleUseCase,
	listUseCase *apparticle.ListArticlesUseCase,
	updateUseCase *apparticle.UpdateArticleUseCase,
	deleteUseCase *apparticle.DeleteArticleUseCase,
) *Handler {
	return &Handler{
		createUseCase: createUseCase,
		getUseCase:    getUseCase,
		listUseCase:   listUseCase,
		updateUseCase: updateUseCase,
		deleteUseCase: deleteUseCase,
	}
}

// Create handles POST /articles
func (h *Handler) Create(c *gin.Context) {
	var req apparticle.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponseBadRequest(c, err.Error())
		return
	}

	resp, err := h.createUseCase.Execute(c.Request.Context(), req)
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

	resp, err := h.getUseCase.Execute(c.Request.Context(), id)
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

	resp, err := h.listUseCase.Execute(c.Request.Context(), limit, offset)
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

	var req apparticle.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponseBadRequest(c, err.Error())
		return
	}

	resp, err := h.updateUseCase.Execute(c.Request.Context(), id, req)
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

	err = h.deleteUseCase.Execute(c.Request.Context(), id)
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

