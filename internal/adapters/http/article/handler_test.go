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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockCreateArticleUseCase is a mock implementation of CreateArticleUseCase
type mockCreateArticleUseCase struct {
	mock.Mock
}

func (m *mockCreateArticleUseCase) Execute(ctx context.Context, req dto.CreateArticleRequest) (*dto.ArticleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ArticleResponse), args.Error(1)
}

// mockGetArticleUseCase is a mock implementation of GetArticleUseCase
type mockGetArticleUseCase struct {
	mock.Mock
}

func (m *mockGetArticleUseCase) Execute(ctx context.Context, id int64) (*dto.ArticleResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ArticleResponse), args.Error(1)
}

// mockListArticlesUseCase is a mock implementation of ListArticlesUseCase
type mockListArticlesUseCase struct {
	mock.Mock
}

func (m *mockListArticlesUseCase) Execute(ctx context.Context, limit, offset int) (*dto.ListArticlesResponse, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ListArticlesResponse), args.Error(1)
}

// mockUpdateArticleUseCase is a mock implementation of UpdateArticleUseCase
type mockUpdateArticleUseCase struct {
	mock.Mock
}

func (m *mockUpdateArticleUseCase) Execute(ctx context.Context, id int64, req dto.UpdateArticleRequest) (*dto.ArticleResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ArticleResponse), args.Error(1)
}

// mockDeleteArticleUseCase is a mock implementation of DeleteArticleUseCase
type mockDeleteArticleUseCase struct {
	mock.Mock
}

func (m *mockDeleteArticleUseCase) Execute(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupTestRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestNewHandler(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	assert.NotNil(t, handler)
	assert.Equal(t, createUC, handler.createUseCase)
	assert.Equal(t, getUC, handler.getUseCase)
	assert.Equal(t, listUC, handler.listUseCase)
	assert.Equal(t, updateUC, handler.updateUseCase)
	assert.Equal(t, deleteUC, handler.deleteUseCase)
}

func TestHandler_Create_Success(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

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

	createUC.On("Execute", mock.Anything, reqBody).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.POST("/articles", handler.Create)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	createUC.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Article created successfully", response["message"])
}

func TestHandler_Create_BadRequest_InvalidJSON(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	router := setupTestRouter(handler)
	router.POST("/articles", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	createUC.AssertNotCalled(t, "Execute")
}

func TestHandler_Create_BadRequest_MissingFields(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

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
	createUC.AssertNotCalled(t, "Execute")
}

func TestHandler_Create_InternalServerError(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	reqBody := dto.CreateArticleRequest{
		Title:    "Test Article",
		Content:  "Test Content",
		AuthorID: 1,
	}

	createUC.On("Execute", mock.Anything, reqBody).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.POST("/articles", handler.Create)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	createUC.AssertExpectations(t)
}

func TestHandler_Get_Success(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	articleID := int64(1)
	expectedResp := &dto.ArticleResponse{
		ID:        articleID,
		Title:     "Test Article",
		Content:   "Test Content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	getUC.On("Execute", mock.Anything, articleID).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.GET("/articles/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/articles/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	getUC.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Article retrieved successfully", response["message"])
}

func TestHandler_Get_BadRequest_InvalidID(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	router := setupTestRouter(handler)
	router.GET("/articles/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/articles/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	getUC.AssertNotCalled(t, "Execute")
}

func TestHandler_Get_NotFound(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	articleID := int64(999)
	getUC.On("Execute", mock.Anything, articleID).Return(nil, domainarticle.ErrArticleNotFound)

	router := setupTestRouter(handler)
	router.GET("/articles/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/articles/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	getUC.AssertExpectations(t)
}

func TestHandler_Get_InternalServerError(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	articleID := int64(1)
	getUC.On("Execute", mock.Anything, articleID).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.GET("/articles/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/articles/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	getUC.AssertExpectations(t)
}

func TestHandler_List_Success(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

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

	listUC.On("Execute", mock.Anything, limit, offset).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.GET("/articles", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/articles", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	listUC.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Articles retrieved successfully", response["message"])
}

func TestHandler_List_WithQueryParams(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	limit := 20
	offset := 10
	expectedResp := &dto.ListArticlesResponse{
		Articles: []dto.ArticleResponse{},
		Total:    0,
		Limit:    limit,
		Offset:   offset,
	}

	listUC.On("Execute", mock.Anything, limit, offset).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.GET("/articles", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/articles?limit=20&offset=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	listUC.AssertExpectations(t)
}

func TestHandler_List_InternalServerError(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	limit := 10
	offset := 0
	listUC.On("Execute", mock.Anything, limit, offset).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.GET("/articles", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/articles", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	listUC.AssertExpectations(t)
}

func TestHandler_Update_Success(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

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

	updateUC.On("Execute", mock.Anything, articleID, reqBody).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.PUT("/articles/:id", handler.Update)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/articles/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	updateUC.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Article updated successfully", response["message"])
}

func TestHandler_Update_BadRequest_InvalidID(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

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
	updateUC.AssertNotCalled(t, "Execute")
}

func TestHandler_Update_BadRequest_InvalidJSON(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	router := setupTestRouter(handler)
	router.PUT("/articles/:id", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/articles/1", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	updateUC.AssertNotCalled(t, "Execute")
}

func TestHandler_Update_BadRequest_MissingFields(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

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
	updateUC.AssertNotCalled(t, "Execute")
}

func TestHandler_Update_NotFound(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	articleID := int64(999)
	reqBody := dto.UpdateArticleRequest{
		Title:   "Updated Article",
		Content: "Updated Content",
	}

	updateUC.On("Execute", mock.Anything, articleID, reqBody).Return(nil, domainarticle.ErrArticleNotFound)

	router := setupTestRouter(handler)
	router.PUT("/articles/:id", handler.Update)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/articles/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	updateUC.AssertExpectations(t)
}

func TestHandler_Update_InternalServerError(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	articleID := int64(1)
	reqBody := dto.UpdateArticleRequest{
		Title:   "Updated Article",
		Content: "Updated Content",
	}

	updateUC.On("Execute", mock.Anything, articleID, reqBody).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.PUT("/articles/:id", handler.Update)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/articles/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	updateUC.AssertExpectations(t)
}

func TestHandler_Delete_Success(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	articleID := int64(1)
	deleteUC.On("Execute", mock.Anything, articleID).Return(nil)

	router := setupTestRouter(handler)
	router.DELETE("/articles/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/articles/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	deleteUC.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Article deleted successfully", response["message"])
}

func TestHandler_Delete_BadRequest_InvalidID(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	router := setupTestRouter(handler)
	router.DELETE("/articles/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/articles/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	deleteUC.AssertNotCalled(t, "Execute")
}

func TestHandler_Delete_NotFound(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	articleID := int64(999)
	deleteUC.On("Execute", mock.Anything, articleID).Return(domainarticle.ErrArticleNotFound)

	router := setupTestRouter(handler)
	router.DELETE("/articles/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/articles/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	deleteUC.AssertExpectations(t)
}

func TestHandler_Delete_InternalServerError(t *testing.T) {
	createUC := &mockCreateArticleUseCase{}
	getUC := &mockGetArticleUseCase{}
	listUC := &mockListArticlesUseCase{}
	updateUC := &mockUpdateArticleUseCase{}
	deleteUC := &mockDeleteArticleUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	articleID := int64(1)
	deleteUC.On("Execute", mock.Anything, articleID).Return(errors.New("database error"))

	router := setupTestRouter(handler)
	router.DELETE("/articles/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/articles/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	deleteUC.AssertExpectations(t)
}

