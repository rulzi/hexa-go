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
	articleentity "github.com/rulzi/hexa-go/internal/domain/article/entity"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testLogger logger.Logger = logger.NewSimpleLogger()

type mockCreateArticle struct{ mock.Mock }

func (m *mockCreateArticle) Execute(ctx context.Context, req dto.CreateArticleRequest) (*dto.ArticleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ArticleResponse), args.Error(1)
}

type mockGetArticle struct{ mock.Mock }

func (m *mockGetArticle) Execute(ctx context.Context, id int64) (*dto.ArticleResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ArticleResponse), args.Error(1)
}

type mockListArticle struct{ mock.Mock }

func (m *mockListArticle) Execute(ctx context.Context, limit, offset int) (*dto.ListArticlesResponse, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ListArticlesResponse), args.Error(1)
}

type mockUpdateArticle struct{ mock.Mock }

func (m *mockUpdateArticle) Execute(ctx context.Context, id int64, req dto.UpdateArticleRequest) (*dto.ArticleResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ArticleResponse), args.Error(1)
}

type mockDeleteArticle struct{ mock.Mock }

func (m *mockDeleteArticle) Execute(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type articleTestMocks struct {
	Create *mockCreateArticle
	Get    *mockGetArticle
	List   *mockListArticle
	Update *mockUpdateArticle
	Delete *mockDeleteArticle
}

func newTestHandler() (*Handler, *articleTestMocks) {
	mocks := &articleTestMocks{
		Create: &mockCreateArticle{},
		Get:    &mockGetArticle{},
		List:   &mockListArticle{},
		Update: &mockUpdateArticle{},
		Delete: &mockDeleteArticle{},
	}
	handler := NewHandlerWithDeps(Deps{
		Create: mocks.Create,
		Get:    mocks.Get,
		List:   mocks.List,
		Update: mocks.Update,
		Delete: mocks.Delete,
	}, testLogger)
	return handler, mocks
}

func setupTestRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestNewHandler(t *testing.T) {
	handler, mocks := newTestHandler()
	assert.NotNil(t, handler)
	assert.NotNil(t, mocks.Create)
}

func TestHandler_Create_Success(t *testing.T) {
	handler, mocks := newTestHandler()

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

	mocks.Create.On("Execute", mock.Anything, reqBody).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.POST("/articles", handler.Create)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mocks.Create.AssertExpectations(t); mocks.Get.AssertExpectations(t); mocks.List.AssertExpectations(t); mocks.Update.AssertExpectations(t); mocks.Delete.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Article created successfully", response["message"])
}

func TestHandler_Create_BadRequest_InvalidJSON(t *testing.T) {
	handler, mocks := newTestHandler()

	router := setupTestRouter(handler)
	router.POST("/articles", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mocks.Create.AssertNotCalled(t, "Execute")
}

func TestHandler_Create_BadRequest_MissingFields(t *testing.T) {
	handler, mocks := newTestHandler()

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
	mocks.Create.AssertNotCalled(t, "Execute")
}

func TestHandler_Create_InternalServerError(t *testing.T) {
	handler, mocks := newTestHandler()

	reqBody := dto.CreateArticleRequest{
		Title:    "Test Article",
		Content:  "Test Content",
		AuthorID: 1,
	}

	mocks.Create.On("Execute", mock.Anything, reqBody).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.POST("/articles", handler.Create)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mocks.Create.AssertExpectations(t); mocks.Get.AssertExpectations(t); mocks.List.AssertExpectations(t); mocks.Update.AssertExpectations(t); mocks.Delete.AssertExpectations(t)
}

func TestHandler_Get_Success(t *testing.T) {
	handler, mocks := newTestHandler()

	articleID := int64(1)
	expectedResp := &dto.ArticleResponse{
		ID:        articleID,
		Title:     "Test Article",
		Content:   "Test Content",
		AuthorID:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mocks.Get.On("Execute", mock.Anything, articleID).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.GET("/articles/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/articles/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mocks.Create.AssertExpectations(t); mocks.Get.AssertExpectations(t); mocks.List.AssertExpectations(t); mocks.Update.AssertExpectations(t); mocks.Delete.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Article retrieved successfully", response["message"])
}

func TestHandler_Get_BadRequest_InvalidID(t *testing.T) {
	handler, mocks := newTestHandler()

	router := setupTestRouter(handler)
	router.GET("/articles/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/articles/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mocks.Get.AssertNotCalled(t, "Execute")
}

func TestHandler_Get_NotFound(t *testing.T) {
	handler, mocks := newTestHandler()

	articleID := int64(999)
	mocks.Get.On("Execute", mock.Anything, articleID).Return(nil, articleentity.NewArticleNotFound())

	router := setupTestRouter(handler)
	router.GET("/articles/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/articles/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mocks.Create.AssertExpectations(t); mocks.Get.AssertExpectations(t); mocks.List.AssertExpectations(t); mocks.Update.AssertExpectations(t); mocks.Delete.AssertExpectations(t)
}

func TestHandler_Get_InternalServerError(t *testing.T) {
	handler, mocks := newTestHandler()

	articleID := int64(1)
	mocks.Get.On("Execute", mock.Anything, articleID).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.GET("/articles/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/articles/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mocks.Create.AssertExpectations(t); mocks.Get.AssertExpectations(t); mocks.List.AssertExpectations(t); mocks.Update.AssertExpectations(t); mocks.Delete.AssertExpectations(t)
}

func TestHandler_List_Success(t *testing.T) {
	handler, mocks := newTestHandler()

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

	mocks.List.On("Execute", mock.Anything, limit, offset).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.GET("/articles", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/articles", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mocks.Create.AssertExpectations(t); mocks.Get.AssertExpectations(t); mocks.List.AssertExpectations(t); mocks.Update.AssertExpectations(t); mocks.Delete.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Articles retrieved successfully", response["message"])
}

func TestHandler_List_WithQueryParams(t *testing.T) {
	handler, mocks := newTestHandler()

	limit := 20
	offset := 10
	expectedResp := &dto.ListArticlesResponse{
		Articles: []dto.ArticleResponse{},
		Total:    0,
		Limit:    limit,
		Offset:   offset,
	}

	mocks.List.On("Execute", mock.Anything, limit, offset).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.GET("/articles", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/articles?limit=20&offset=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mocks.Create.AssertExpectations(t); mocks.Get.AssertExpectations(t); mocks.List.AssertExpectations(t); mocks.Update.AssertExpectations(t); mocks.Delete.AssertExpectations(t)
}

func TestHandler_List_InternalServerError(t *testing.T) {
	handler, mocks := newTestHandler()

	limit := 10
	offset := 0
	mocks.List.On("Execute", mock.Anything, limit, offset).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.GET("/articles", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/articles", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mocks.Create.AssertExpectations(t); mocks.Get.AssertExpectations(t); mocks.List.AssertExpectations(t); mocks.Update.AssertExpectations(t); mocks.Delete.AssertExpectations(t)
}

func TestHandler_Update_Success(t *testing.T) {
	handler, mocks := newTestHandler()

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

	mocks.Update.On("Execute", mock.Anything, articleID, reqBody).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.PUT("/articles/:id", handler.Update)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/articles/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mocks.Create.AssertExpectations(t); mocks.Get.AssertExpectations(t); mocks.List.AssertExpectations(t); mocks.Update.AssertExpectations(t); mocks.Delete.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Article updated successfully", response["message"])
}

func TestHandler_Update_BadRequest_InvalidID(t *testing.T) {
	handler, mocks := newTestHandler()

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
	mocks.Update.AssertNotCalled(t, "Execute")
}

func TestHandler_Update_BadRequest_InvalidJSON(t *testing.T) {
	handler, mocks := newTestHandler()

	router := setupTestRouter(handler)
	router.PUT("/articles/:id", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/articles/1", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mocks.Update.AssertNotCalled(t, "Execute")
}

func TestHandler_Update_BadRequest_MissingFields(t *testing.T) {
	handler, mocks := newTestHandler()

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
	mocks.Update.AssertNotCalled(t, "Execute")
}

func TestHandler_Update_NotFound(t *testing.T) {
	handler, mocks := newTestHandler()

	articleID := int64(999)
	reqBody := dto.UpdateArticleRequest{
		Title:   "Updated Article",
		Content: "Updated Content",
	}

	mocks.Update.On("Execute", mock.Anything, articleID, reqBody).Return(nil, articleentity.NewArticleNotFound())

	router := setupTestRouter(handler)
	router.PUT("/articles/:id", handler.Update)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/articles/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mocks.Create.AssertExpectations(t); mocks.Get.AssertExpectations(t); mocks.List.AssertExpectations(t); mocks.Update.AssertExpectations(t); mocks.Delete.AssertExpectations(t)
}

func TestHandler_Update_InternalServerError(t *testing.T) {
	handler, mocks := newTestHandler()

	articleID := int64(1)
	reqBody := dto.UpdateArticleRequest{
		Title:   "Updated Article",
		Content: "Updated Content",
	}

	mocks.Update.On("Execute", mock.Anything, articleID, reqBody).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.PUT("/articles/:id", handler.Update)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/articles/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mocks.Create.AssertExpectations(t); mocks.Get.AssertExpectations(t); mocks.List.AssertExpectations(t); mocks.Update.AssertExpectations(t); mocks.Delete.AssertExpectations(t)
}

func TestHandler_Delete_Success(t *testing.T) {
	handler, mocks := newTestHandler()

	articleID := int64(1)
	mocks.Delete.On("Execute", mock.Anything, articleID).Return(nil)

	router := setupTestRouter(handler)
	router.DELETE("/articles/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/articles/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mocks.Create.AssertExpectations(t); mocks.Get.AssertExpectations(t); mocks.List.AssertExpectations(t); mocks.Update.AssertExpectations(t); mocks.Delete.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Article deleted successfully", response["message"])
}

func TestHandler_Delete_BadRequest_InvalidID(t *testing.T) {
	handler, mocks := newTestHandler()

	router := setupTestRouter(handler)
	router.DELETE("/articles/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/articles/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mocks.Delete.AssertNotCalled(t, "Execute")
}

func TestHandler_Delete_NotFound(t *testing.T) {
	handler, mocks := newTestHandler()

	articleID := int64(999)
	mocks.Delete.On("Execute", mock.Anything, articleID).Return(articleentity.NewArticleNotFound())

	router := setupTestRouter(handler)
	router.DELETE("/articles/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/articles/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mocks.Create.AssertExpectations(t); mocks.Get.AssertExpectations(t); mocks.List.AssertExpectations(t); mocks.Update.AssertExpectations(t); mocks.Delete.AssertExpectations(t)
}

func TestHandler_Delete_InternalServerError(t *testing.T) {
	handler, mocks := newTestHandler()

	articleID := int64(1)
	mocks.Delete.On("Execute", mock.Anything, articleID).Return(errors.New("database error"))

	router := setupTestRouter(handler)
	router.DELETE("/articles/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/articles/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mocks.Create.AssertExpectations(t); mocks.Get.AssertExpectations(t); mocks.List.AssertExpectations(t); mocks.Update.AssertExpectations(t); mocks.Delete.AssertExpectations(t)
}

