package http

import (
	"net/http"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.createUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Get handles GET /articles/:id
func (h *ArticleHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid article id"})
		return
	}

	response, err := h.getUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == domainarticle.ErrArticleNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// List handles GET /articles
func (h *ArticleHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	response, err := h.listUseCase.Execute(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Update handles PUT /articles/:id
func (h *ArticleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid article id"})
		return
	}

	var req article.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.updateUseCase.Execute(c.Request.Context(), id, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == domainarticle.ErrArticleNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Delete handles DELETE /articles/:id
func (h *ArticleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid article id"})
		return
	}

	err = h.deleteUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == domainarticle.ErrArticleNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "article deleted successfully"})
}
