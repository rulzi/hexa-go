package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rulzi/hexa-go/internal/application/article"
	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
)

// ArticleHandler handles HTTP requests for articles
type ArticleHandler struct {
	createUseCase *article.CreateArticleUseCase
	getUseCase    *article.GetArticleUseCase
	listUseCase   *article.ListArticlesUseCase
	updateUseCase *article.UpdateArticleUseCase
	deleteUseCase *article.DeleteArticleUseCase
}

// NewArticleHandler creates a new ArticleHandler
func NewArticleHandler(
	createUseCase *article.CreateArticleUseCase,
	getUseCase *article.GetArticleUseCase,
	listUseCase *article.ListArticlesUseCase,
	updateUseCase *article.UpdateArticleUseCase,
	deleteUseCase *article.DeleteArticleUseCase,
) *ArticleHandler {
	return &ArticleHandler{
		createUseCase: createUseCase,
		getUseCase:    getUseCase,
		listUseCase:   listUseCase,
		updateUseCase: updateUseCase,
		deleteUseCase: deleteUseCase,
	}
}

// Create handles POST /articles
func (h *ArticleHandler) Create(c *gin.Context) {
	var req article.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponseBadRequest(c, err.Error())
		return
	}

	response, err := h.createUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		ErrorResponseInternalServerError(c, err.Error())
		return
	}

	SuccessResponseCreated(c, "Article created successfully", response)
}

// Get handles GET /articles/:id
func (h *ArticleHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ErrorResponseBadRequest(c, "invalid article id")
		return
	}

	response, err := h.getUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		if err == domainarticle.ErrArticleNotFound {
			ErrorResponseNotFound(c, err.Error())
		} else {
			ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponseOK(c, "Article retrieved successfully", response)
}

// List handles GET /articles
func (h *ArticleHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	response, err := h.listUseCase.Execute(c.Request.Context(), limit, offset)
	if err != nil {
		ErrorResponseInternalServerError(c, err.Error())
		return
	}

	SuccessResponseOK(c, "Articles retrieved successfully", response)
}

// Update handles PUT /articles/:id
func (h *ArticleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ErrorResponseBadRequest(c, "invalid article id")
		return
	}

	var req article.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponseBadRequest(c, err.Error())
		return
	}

	response, err := h.updateUseCase.Execute(c.Request.Context(), id, req)
	if err != nil {
		if err == domainarticle.ErrArticleNotFound {
			ErrorResponseNotFound(c, err.Error())
		} else {
			ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponseOK(c, "Article updated successfully", response)
}

// Delete handles DELETE /articles/:id
func (h *ArticleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ErrorResponseBadRequest(c, "invalid article id")
		return
	}

	err = h.deleteUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		if err == domainarticle.ErrArticleNotFound {
			ErrorResponseNotFound(c, err.Error())
		} else {
			ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponseOK(c, "Article deleted successfully", nil)
}
