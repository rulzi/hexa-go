package article

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rulzi/hexa-go/internal/application/article/dto"
	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testLogger logger.Logger = logger.NewSimpleLogger()

// mockArticleUseCase is a mock implementation of ArticleUseCase
type mockArticleUseCase struct {
	mock.Mock
}

func (m *mockArticleUseCase) Create(ctx context.Context, req dto.CreateArticleRequest) (*dto.ArticleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ArticleResponse), args.Error(1)
}

func (m *mockArticleUseCase) Get(ctx context.Context, id int64) (*dto.ArticleResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ArticleResponse), args.Error(1)
}

func (m *mockArticleUseCase) List(ctx context.Context, limit, offset int) (*dto.ListArticlesResponse, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ListArticlesResponse), args.Error(1)
}

func (m *mockArticleUseCase) Update(ctx context.Context, id int64, req dto.UpdateArticleRequest) (*dto.ArticleResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ArticleResponse), args.Error(1)
}

func (m *mockArticleUseCase) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupTestRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestNewHandler(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)
	assert.NotNil(t, handler)
	assert.Equal(t, uc, handler.articleUseCase)
}

func TestHandler_Create_Success(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)

	reqBody := dto.CreateArticleRequest{
		Title:    "Test Article",
		Content:  "Test Content",
		AuthorID: 1,
	}

	expectedResp := &dto.ArticleResponse{
		ID:        1,
		Title:     reqBody.Title,
		Content:   reqBody.Content,
		AuthorID:  reqBody.AuthorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	uc.On("Create", mock.Anything, reqBody).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.POST("/articles", handler.Create)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	uc.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Article created successfully", response["message"])
}

func TestHandler_Create_BadRequest_InvalidJSON(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)

	router := setupTestRouter(handler)
	router.POST("/articles", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "Create")
}

func TestHandler_Create_BadRequest_MissingFields(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)

	reqBody := map[string]interface{}{
		"title": "Test Article",
		// Missing content and author_id
	}

	router := setupTestRouter(handler)
	router.POST("/articles", handler.Create)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "Create")
}

func TestHandler_Create_InternalServerError(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)

	reqBody := dto.CreateArticleRequest{
		Title:    "Test Article",
		Content:  "Test Content",
		AuthorID: 1,
	}

	uc.On("Create", mock.Anything, reqBody).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.POST("/articles", handler.Create)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Get_Success(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)

	articleID := int64(1)
	expectedResp := &dto.ArticleResponse{
		ID:        articleID,
		Title:     "Test Article",
		Content:   "Test Content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	uc.On("Get", mock.Anything, articleID).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.GET("/articles/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/articles/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Article retrieved successfully", response["message"])
}

func TestHandler_Get_BadRequest_InvalidID(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)

	router := setupTestRouter(handler)
	router.GET("/articles/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/articles/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "Get")
}

func TestHandler_Get_NotFound(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)

	articleID := int64(999)
	uc.On("Get", mock.Anything, articleID).Return(nil, domainarticle.ErrArticleNotFound)

	router := setupTestRouter(handler)
	router.GET("/articles/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/articles/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Get_InternalServerError(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)

	articleID := int64(1)
	uc.On("Get", mock.Anything, articleID).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.GET("/articles/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/articles/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_List_Success(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)

	limit := 10
	offset := 0
	expectedResp := &dto.ListArticlesResponse{
		Articles: []dto.ArticleResponse{
			{
				ID:        1,
				Title:     "Test Article 1",
				Content:   "Test Content 1",
				AuthorID:  1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		Total:  1,
		Limit:  limit,
		Offset: offset,
	}

	uc.On("List", mock.Anything, limit, offset).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.GET("/articles", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/articles", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Articles retrieved successfully", response["message"])
}

func TestHandler_List_WithQueryParams(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)

	limit := 20
	offset := 10
	expectedResp := &dto.ListArticlesResponse{
		Articles: []dto.ArticleResponse{},
		Total:    0,
		Limit:    limit,
		Offset:   offset,
	}

	uc.On("List", mock.Anything, limit, offset).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.GET("/articles", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/articles?limit=20&offset=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_List_InternalServerError(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)

	limit := 10
	offset := 0
	uc.On("List", mock.Anything, limit, offset).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.GET("/articles", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/articles", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Update_Success(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)

	articleID := int64(1)
	reqBody := dto.UpdateArticleRequest{
		Title:   "Updated Article",
		Content: "Updated Content",
	}

	expectedResp := &dto.ArticleResponse{
		ID:        articleID,
		Title:     reqBody.Title,
		Content:   reqBody.Content,
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	uc.On("Update", mock.Anything, articleID, reqBody).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.PUT("/articles/:id", handler.Update)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/articles/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Article updated successfully", response["message"])
}

func TestHandler_Update_BadRequest_InvalidID(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)

	reqBody := dto.UpdateArticleRequest{
		Title:   "Updated Article",
		Content: "Updated Content",
	}

	router := setupTestRouter(handler)
	router.PUT("/articles/:id", handler.Update)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/articles/invalid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "Update")
}

func TestHandler_Update_BadRequest_InvalidJSON(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)

	router := setupTestRouter(handler)
	router.PUT("/articles/:id", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/articles/1", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "Update")
}

func TestHandler_Update_BadRequest_MissingFields(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)

	reqBody := map[string]interface{}{
		"title": "Updated Article",
		// Missing content
	}

	router := setupTestRouter(handler)
	router.PUT("/articles/:id", handler.Update)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/articles/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "Update")
}

func TestHandler_Update_NotFound(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)

	articleID := int64(999)
	reqBody := dto.UpdateArticleRequest{
		Title:   "Updated Article",
		Content: "Updated Content",
	}

	uc.On("Update", mock.Anything, articleID, reqBody).Return(nil, domainarticle.ErrArticleNotFound)

	router := setupTestRouter(handler)
	router.PUT("/articles/:id", handler.Update)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/articles/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Update_InternalServerError(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)

	articleID := int64(1)
	reqBody := dto.UpdateArticleRequest{
		Title:   "Updated Article",
		Content: "Updated Content",
	}

	uc.On("Update", mock.Anything, articleID, reqBody).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.PUT("/articles/:id", handler.Update)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/articles/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Delete_Success(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)

	articleID := int64(1)
	uc.On("Delete", mock.Anything, articleID).Return(nil)

	router := setupTestRouter(handler)
	router.DELETE("/articles/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/articles/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Article deleted successfully", response["message"])
}

func TestHandler_Delete_BadRequest_InvalidID(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)

	router := setupTestRouter(handler)
	router.DELETE("/articles/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/articles/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "Delete")
}

func TestHandler_Delete_NotFound(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)

	articleID := int64(999)
	uc.On("Delete", mock.Anything, articleID).Return(domainarticle.ErrArticleNotFound)

	router := setupTestRouter(handler)
	router.DELETE("/articles/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/articles/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Delete_InternalServerError(t *testing.T) {
	uc := &mockArticleUseCase{}
	handler := NewHandler(uc, testLogger)

	articleID := int64(1)
	uc.On("Delete", mock.Anything, articleID).Return(errors.New("database error"))

	router := setupTestRouter(handler)
	router.DELETE("/articles/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/articles/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	uc.AssertExpectations(t)
}

