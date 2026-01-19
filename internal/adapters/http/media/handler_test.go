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
	domainmedia "github.com/rulzi/hexa-go/internal/domain/media"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockCreateMediaUseCase is a mock implementation of CreateMediaUseCase
type mockCreateMediaUseCase struct {
	mock.Mock
}

func (m *mockCreateMediaUseCase) Execute(ctx context.Context, filename string, file io.Reader) (*dto.MediaResponse, error) {
	args := m.Called(ctx, filename, file)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MediaResponse), args.Error(1)
}

// mockGetMediaUseCase is a mock implementation of GetMediaUseCase
type mockGetMediaUseCase struct {
	mock.Mock
}

func (m *mockGetMediaUseCase) Execute(ctx context.Context, id int64) (*dto.MediaResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MediaResponse), args.Error(1)
}

// mockListMediaUseCase is a mock implementation of ListMediaUseCase
type mockListMediaUseCase struct {
	mock.Mock
}

func (m *mockListMediaUseCase) Execute(ctx context.Context, limit, offset int) (*dto.ListMediaResponse, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ListMediaResponse), args.Error(1)
}

// mockUpdateMediaUseCase is a mock implementation of UpdateMediaUseCase
type mockUpdateMediaUseCase struct {
	mock.Mock
}

func (m *mockUpdateMediaUseCase) Execute(ctx context.Context, id int64, filename string, file io.Reader) (*dto.MediaResponse, error) {
	args := m.Called(ctx, id, filename, file)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MediaResponse), args.Error(1)
}

// mockDeleteMediaUseCase is a mock implementation of DeleteMediaUseCase
type mockDeleteMediaUseCase struct {
	mock.Mock
}

func (m *mockDeleteMediaUseCase) Execute(ctx context.Context, id int64) error {
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
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	assert.NotNil(t, handler)
	assert.Equal(t, createUC, handler.createUseCase)
	assert.Equal(t, getUC, handler.getUseCase)
	assert.Equal(t, listUC, handler.listUseCase)
	assert.Equal(t, updateUC, handler.updateUseCase)
	assert.Equal(t, deleteUC, handler.deleteUseCase)
}

func TestHandler_Create_Success(t *testing.T) {
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

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
	createUC.On("Execute", mock.Anything, filename, mock.Anything).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.POST("/media", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/media", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	createUC.AssertExpectations(t)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Media created successfully", response["message"])
}

func TestHandler_Create_BadRequest_NoFile(t *testing.T) {
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	router := setupTestRouter(handler)
	router.POST("/media", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/media", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	createUC.AssertNotCalled(t, "Execute")
}

func TestHandler_Create_BadRequest_ValidationError(t *testing.T) {
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	filename := "test.jpg"
	fileContent := "test file content"

	// Create multipart form
	body, contentType, err := createMultipartFormData(filename, fileContent)
	assert.NoError(t, err)

	createUC.On("Execute", mock.Anything, filename, mock.Anything).Return(nil, domainmedia.ErrNameRequired)

	router := setupTestRouter(handler)
	router.POST("/media", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/media", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	createUC.AssertExpectations(t)
}

func TestHandler_Create_InternalServerError(t *testing.T) {
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	filename := "test.jpg"
	fileContent := "test file content"

	// Create multipart form
	body, contentType, err := createMultipartFormData(filename, fileContent)
	assert.NoError(t, err)

	createUC.On("Execute", mock.Anything, filename, mock.Anything).Return(nil, errors.New("storage error"))

	router := setupTestRouter(handler)
	router.POST("/media", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/media", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	createUC.AssertExpectations(t)
}

func TestHandler_Get_Success(t *testing.T) {
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	mediaID := int64(1)
	expectedResp := &dto.MediaResponse{
		ID:        mediaID,
		Name:      "test.jpg",
		Path:      "2025/12/19/test.jpg",
		URL:       "http://localhost:8080/api/v1/media/files/2025/12/19/test.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	getUC.On("Execute", mock.Anything, mediaID).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.GET("/media/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/media/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	getUC.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Media retrieved successfully", response["message"])
}

func TestHandler_Get_BadRequest_InvalidID(t *testing.T) {
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	router := setupTestRouter(handler)
	router.GET("/media/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/media/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	getUC.AssertNotCalled(t, "Execute")
}

func TestHandler_Get_NotFound(t *testing.T) {
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	mediaID := int64(999)
	getUC.On("Execute", mock.Anything, mediaID).Return(nil, domainmedia.ErrMediaNotFound)

	router := setupTestRouter(handler)
	router.GET("/media/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/media/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	getUC.AssertExpectations(t)
}

func TestHandler_Get_InternalServerError(t *testing.T) {
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	mediaID := int64(1)
	getUC.On("Execute", mock.Anything, mediaID).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.GET("/media/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/media/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	getUC.AssertExpectations(t)
}

func TestHandler_List_Success(t *testing.T) {
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

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

	listUC.On("Execute", mock.Anything, limit, offset).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.GET("/media", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/media", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	listUC.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Media retrieved successfully", response["message"])
}

func TestHandler_List_WithQueryParams(t *testing.T) {
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	limit := 20
	offset := 10
	expectedResp := &dto.ListMediaResponse{
		Media:  []dto.MediaResponse{},
		Total:  0,
		Limit:  limit,
		Offset: offset,
	}

	listUC.On("Execute", mock.Anything, limit, offset).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.GET("/media", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/media?limit=20&offset=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	listUC.AssertExpectations(t)
}

func TestHandler_List_InternalServerError(t *testing.T) {
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	limit := 10
	offset := 0
	listUC.On("Execute", mock.Anything, limit, offset).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.GET("/media", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/media", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	listUC.AssertExpectations(t)
}

func TestHandler_Update_Success(t *testing.T) {
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

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

	updateUC.On("Execute", mock.Anything, mediaID, filename, mock.Anything).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.PUT("/media/:id", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/media/1", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	updateUC.AssertExpectations(t)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Media updated successfully", response["message"])
}

func TestHandler_Update_BadRequest_InvalidID(t *testing.T) {
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

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
	updateUC.AssertNotCalled(t, "Execute")
}

func TestHandler_Update_BadRequest_NoFile(t *testing.T) {
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	router := setupTestRouter(handler)
	router.PUT("/media/:id", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/media/1", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	updateUC.AssertNotCalled(t, "Execute")
}

func TestHandler_Update_BadRequest_ValidationError(t *testing.T) {
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	mediaID := int64(1)
	filename := "updated.jpg"
	fileContent := "updated file content"

	// Create multipart form
	body, contentType, err := createMultipartFormData(filename, fileContent)
	assert.NoError(t, err)

	updateUC.On("Execute", mock.Anything, mediaID, filename, mock.Anything).Return(nil, domainmedia.ErrNameRequired)

	router := setupTestRouter(handler)
	router.PUT("/media/:id", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/media/1", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	updateUC.AssertExpectations(t)
}

func TestHandler_Update_NotFound(t *testing.T) {
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	mediaID := int64(999)
	filename := "updated.jpg"
	fileContent := "updated file content"

	// Create multipart form
	body, contentType, err := createMultipartFormData(filename, fileContent)
	assert.NoError(t, err)

	updateUC.On("Execute", mock.Anything, mediaID, filename, mock.Anything).Return(nil, domainmedia.ErrMediaNotFound)

	router := setupTestRouter(handler)
	router.PUT("/media/:id", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/media/999", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	updateUC.AssertExpectations(t)
}

func TestHandler_Update_InternalServerError(t *testing.T) {
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	mediaID := int64(1)
	filename := "updated.jpg"
	fileContent := "updated file content"

	// Create multipart form
	body, contentType, err := createMultipartFormData(filename, fileContent)
	assert.NoError(t, err)

	updateUC.On("Execute", mock.Anything, mediaID, filename, mock.Anything).Return(nil, errors.New("storage error"))

	router := setupTestRouter(handler)
	router.PUT("/media/:id", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/media/1", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	updateUC.AssertExpectations(t)
}

func TestHandler_Delete_Success(t *testing.T) {
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	mediaID := int64(1)
	deleteUC.On("Execute", mock.Anything, mediaID).Return(nil)

	router := setupTestRouter(handler)
	router.DELETE("/media/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/media/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	deleteUC.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Media deleted successfully", response["message"])
}

func TestHandler_Delete_BadRequest_InvalidID(t *testing.T) {
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	router := setupTestRouter(handler)
	router.DELETE("/media/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/media/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	deleteUC.AssertNotCalled(t, "Execute")
}

func TestHandler_Delete_NotFound(t *testing.T) {
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	mediaID := int64(999)
	deleteUC.On("Execute", mock.Anything, mediaID).Return(domainmedia.ErrMediaNotFound)

	router := setupTestRouter(handler)
	router.DELETE("/media/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/media/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	deleteUC.AssertExpectations(t)
}

func TestHandler_Delete_InternalServerError(t *testing.T) {
	createUC := &mockCreateMediaUseCase{}
	getUC := &mockGetMediaUseCase{}
	listUC := &mockListMediaUseCase{}
	updateUC := &mockUpdateMediaUseCase{}
	deleteUC := &mockDeleteMediaUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC)

	mediaID := int64(1)
	deleteUC.On("Execute", mock.Anything, mediaID).Return(errors.New("database error"))

	router := setupTestRouter(handler)
	router.DELETE("/media/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/media/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	deleteUC.AssertExpectations(t)
}
