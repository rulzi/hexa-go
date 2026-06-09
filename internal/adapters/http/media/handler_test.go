package media

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rulzi/hexa-go/internal/application/media/dto"
	mediaentity "github.com/rulzi/hexa-go/internal/domain/media/entity"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testLogger logger.Logger = logger.NewSimpleLogger()

// mockMediaUseCase is a mock implementation of MediaUseCase
type mockMediaUseCase struct {
	mock.Mock
}

func (m *mockMediaUseCase) Create(ctx context.Context, filename string, file io.Reader) (*dto.MediaResponse, error) {
	args := m.Called(ctx, filename, file)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MediaResponse), args.Error(1)
}

func (m *mockMediaUseCase) Get(ctx context.Context, id int64) (*dto.MediaResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MediaResponse), args.Error(1)
}

func (m *mockMediaUseCase) List(ctx context.Context, limit, offset int) (*dto.ListMediaResponse, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ListMediaResponse), args.Error(1)
}

func (m *mockMediaUseCase) Update(ctx context.Context, id int64, filename string, file io.Reader) (*dto.MediaResponse, error) {
	args := m.Called(ctx, id, filename, file)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MediaResponse), args.Error(1)
}

func (m *mockMediaUseCase) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupTestRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// createMultipartFormData creates a multipart form with a file field
func createMultipartFormData(filename, content string) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, "", err
	}

	_, err = part.Write([]byte(content))
	if err != nil {
		return nil, "", err
	}

	err = writer.Close()
	if err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}

func TestNewHandler(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)
	assert.NotNil(t, handler)
	assert.Equal(t, uc, handler.mediaUseCase)
}

func TestHandler_Create_Success(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)

	filename := "test.jpg"
	fileContent := "test file content"
	expectedResp := &dto.MediaResponse{
		ID:        1,
		Name:      filename,
		Path:      "2025/12/19/test.jpg",
		URL:       "http://localhost:8080/api/v1/media/files/2025/12/19/test.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create multipart form
	body, contentType, err := createMultipartFormData(filename, fileContent)
	assert.NoError(t, err)

	// Mock expects the file content to be read
	uc.On("Create", mock.Anything, filename, mock.Anything).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.POST("/media", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/media", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	uc.AssertExpectations(t)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Media created successfully", response["message"])
}

func TestHandler_Create_BadRequest_NoFile(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)

	router := setupTestRouter(handler)
	router.POST("/media", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/media", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "Create")
}

func TestHandler_Create_BadRequest_ValidationError(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)

	filename := "test.jpg"
	fileContent := "test file content"

	// Create multipart form
	body, contentType, err := createMultipartFormData(filename, fileContent)
	assert.NoError(t, err)

	uc.On("Create", mock.Anything, filename, mock.Anything).Return(nil, mediaentity.NewNameRequired())

	router := setupTestRouter(handler)
	router.POST("/media", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/media", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Create_InternalServerError(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)

	filename := "test.jpg"
	fileContent := "test file content"

	// Create multipart form
	body, contentType, err := createMultipartFormData(filename, fileContent)
	assert.NoError(t, err)

	uc.On("Create", mock.Anything, filename, mock.Anything).Return(nil, errors.New("storage error"))

	router := setupTestRouter(handler)
	router.POST("/media", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/media", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Get_Success(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)

	mediaID := int64(1)
	expectedResp := &dto.MediaResponse{
		ID:        mediaID,
		Name:      "test.jpg",
		Path:      "2025/12/19/test.jpg",
		URL:       "http://localhost:8080/api/v1/media/files/2025/12/19/test.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	uc.On("Get", mock.Anything, mediaID).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.GET("/media/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/media/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Media retrieved successfully", response["message"])
}

func TestHandler_Get_BadRequest_InvalidID(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)

	router := setupTestRouter(handler)
	router.GET("/media/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/media/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "Get")
}

func TestHandler_Get_NotFound(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)

	mediaID := int64(999)
	uc.On("Get", mock.Anything, mediaID).Return(nil, mediaentity.NewMediaNotFound())

	router := setupTestRouter(handler)
	router.GET("/media/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/media/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Get_InternalServerError(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)

	mediaID := int64(1)
	uc.On("Get", mock.Anything, mediaID).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.GET("/media/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/media/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_List_Success(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)

	limit := 10
	offset := 0
	expectedResp := &dto.ListMediaResponse{
		Media: []dto.MediaResponse{
			{
				ID:        1,
				Name:      "test1.jpg",
				Path:      "2025/12/19/test1.jpg",
				URL:       "http://localhost:8080/api/v1/media/files/2025/12/19/test1.jpg",
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
	router.GET("/media", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/media", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Media retrieved successfully", response["message"])
}

func TestHandler_List_WithQueryParams(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)

	limit := 20
	offset := 10
	expectedResp := &dto.ListMediaResponse{
		Media:  []dto.MediaResponse{},
		Total:  0,
		Limit:  limit,
		Offset: offset,
	}

	uc.On("List", mock.Anything, limit, offset).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.GET("/media", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/media?limit=20&offset=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_List_InternalServerError(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)

	limit := 10
	offset := 0
	uc.On("List", mock.Anything, limit, offset).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.GET("/media", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/media", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Update_Success(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)

	mediaID := int64(1)
	filename := "updated.jpg"
	fileContent := "updated file content"
	expectedResp := &dto.MediaResponse{
		ID:        mediaID,
		Name:      filename,
		Path:      "2025/12/19/updated.jpg",
		URL:       "http://localhost:8080/api/v1/media/files/2025/12/19/updated.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create multipart form
	body, contentType, err := createMultipartFormData(filename, fileContent)
	assert.NoError(t, err)

	uc.On("Update", mock.Anything, mediaID, filename, mock.Anything).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.PUT("/media/:id", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/media/1", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Media updated successfully", response["message"])
}

func TestHandler_Update_BadRequest_InvalidID(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)

	filename := "updated.jpg"
	fileContent := "updated file content"

	// Create multipart form
	body, contentType, err := createMultipartFormData(filename, fileContent)
	assert.NoError(t, err)

	router := setupTestRouter(handler)
	router.PUT("/media/:id", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/media/invalid", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "Update")
}

func TestHandler_Update_BadRequest_NoFile(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)

	router := setupTestRouter(handler)
	router.PUT("/media/:id", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/media/1", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "Update")
}

func TestHandler_Update_BadRequest_ValidationError(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)

	mediaID := int64(1)
	filename := "updated.jpg"
	fileContent := "updated file content"

	// Create multipart form
	body, contentType, err := createMultipartFormData(filename, fileContent)
	assert.NoError(t, err)

	uc.On("Update", mock.Anything, mediaID, filename, mock.Anything).Return(nil, mediaentity.NewNameRequired())

	router := setupTestRouter(handler)
	router.PUT("/media/:id", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/media/1", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Update_NotFound(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)

	mediaID := int64(999)
	filename := "updated.jpg"
	fileContent := "updated file content"

	// Create multipart form
	body, contentType, err := createMultipartFormData(filename, fileContent)
	assert.NoError(t, err)

	uc.On("Update", mock.Anything, mediaID, filename, mock.Anything).Return(nil, mediaentity.NewMediaNotFound())

	router := setupTestRouter(handler)
	router.PUT("/media/:id", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/media/999", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Update_InternalServerError(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)

	mediaID := int64(1)
	filename := "updated.jpg"
	fileContent := "updated file content"

	// Create multipart form
	body, contentType, err := createMultipartFormData(filename, fileContent)
	assert.NoError(t, err)

	uc.On("Update", mock.Anything, mediaID, filename, mock.Anything).Return(nil, errors.New("storage error"))

	router := setupTestRouter(handler)
	router.PUT("/media/:id", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/media/1", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Delete_Success(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)

	mediaID := int64(1)
	uc.On("Delete", mock.Anything, mediaID).Return(nil)

	router := setupTestRouter(handler)
	router.DELETE("/media/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/media/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Media deleted successfully", response["message"])
}

func TestHandler_Delete_BadRequest_InvalidID(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)

	router := setupTestRouter(handler)
	router.DELETE("/media/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/media/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "Delete")
}

func TestHandler_Delete_NotFound(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)

	mediaID := int64(999)
	uc.On("Delete", mock.Anything, mediaID).Return(mediaentity.NewMediaNotFound())

	router := setupTestRouter(handler)
	router.DELETE("/media/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/media/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Delete_InternalServerError(t *testing.T) {
	uc := &mockMediaUseCase{}
	handler := NewHandler(uc, testLogger)

	mediaID := int64(1)
	uc.On("Delete", mock.Anything, mediaID).Return(errors.New("database error"))

	router := setupTestRouter(handler)
	router.DELETE("/media/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/media/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	uc.AssertExpectations(t)
}
