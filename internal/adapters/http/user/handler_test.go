package user

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
	"github.com/rulzi/hexa-go/internal/application/user/dto"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockCreateUserUseCase is a mock implementation of CreateUserUseCase
type mockCreateUserUseCase struct {
	mock.Mock
}

func (m *mockCreateUserUseCase) Execute(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

// mockGetUserUseCase is a mock implementation of GetUserUseCase
type mockGetUserUseCase struct {
	mock.Mock
}

func (m *mockGetUserUseCase) Execute(ctx context.Context, id int64) (*dto.UserResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

// mockListUsersUseCase is a mock implementation of ListUsersUseCase
type mockListUsersUseCase struct {
	mock.Mock
}

func (m *mockListUsersUseCase) Execute(ctx context.Context, limit, offset int) (*dto.ListUsersResponse, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ListUsersResponse), args.Error(1)
}

// mockUpdateUserUseCase is a mock implementation of UpdateUserUseCase
type mockUpdateUserUseCase struct {
	mock.Mock
}

func (m *mockUpdateUserUseCase) Execute(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

// mockDeleteUserUseCase is a mock implementation of DeleteUserUseCase
type mockDeleteUserUseCase struct {
	mock.Mock
}

func (m *mockDeleteUserUseCase) Execute(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// mockLoginUseCase is a mock implementation of LoginUseCase
type mockLoginUseCase struct {
	mock.Mock
}

func (m *mockLoginUseCase) Execute(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoginResponse), args.Error(1)
}

func setupTestRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestNewHandler(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	assert.NotNil(t, handler)
	assert.Equal(t, createUC, handler.createUseCase)
	assert.Equal(t, getUC, handler.getUseCase)
	assert.Equal(t, listUC, handler.listUseCase)
	assert.Equal(t, updateUC, handler.updateUseCase)
	assert.Equal(t, deleteUC, handler.deleteUseCase)
	assert.Equal(t, loginUC, handler.loginUseCase)
}

func TestHandler_Create_Success(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	reqBody := dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedResp := &dto.UserResponse{
		ID:        1,
		Name:      reqBody.Name,
		Email:     reqBody.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createUC.On("Execute", mock.Anything, reqBody).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.POST("/users", handler.Create)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	createUC.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "User created successfully", response["message"])
}

func TestHandler_Create_BadRequest_InvalidJSON(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	router := setupTestRouter(handler)
	router.POST("/users", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	createUC.AssertNotCalled(t, "Execute")
}

func TestHandler_Create_Conflict_EmailExists(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	reqBody := dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	createUC.On("Execute", mock.Anything, reqBody).Return(nil, domainuser.ErrEmailExists)

	router := setupTestRouter(handler)
	router.POST("/users", handler.Create)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	createUC.AssertExpectations(t)
}

func TestHandler_Create_InternalServerError(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	reqBody := dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	createUC.On("Execute", mock.Anything, reqBody).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.POST("/users", handler.Create)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	createUC.AssertExpectations(t)
}

func TestHandler_Get_Success(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	userID := int64(1)
	expectedResp := &dto.UserResponse{
		ID:        userID,
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	getUC.On("Execute", mock.Anything, userID).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.GET("/users/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	getUC.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "User retrieved successfully", response["message"])
}

func TestHandler_Get_BadRequest_InvalidID(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	router := setupTestRouter(handler)
	router.GET("/users/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/users/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	getUC.AssertNotCalled(t, "Execute")
}

func TestHandler_Get_NotFound(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	userID := int64(999)
	getUC.On("Execute", mock.Anything, userID).Return(nil, domainuser.ErrUserNotFound)

	router := setupTestRouter(handler)
	router.GET("/users/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	getUC.AssertExpectations(t)
}

func TestHandler_Get_InternalServerError(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	userID := int64(1)
	getUC.On("Execute", mock.Anything, userID).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.GET("/users/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	getUC.AssertExpectations(t)
}

func TestHandler_List_Success(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	limit := 10
	offset := 0
	expectedResp := &dto.ListUsersResponse{
		Users: []dto.UserResponse{
			{
				ID:        1,
				Name:      "Test User 1",
				Email:     "test1@example.com",
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
	router.GET("/users", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	listUC.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Users retrieved successfully", response["message"])
}

func TestHandler_List_WithQueryParams(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	limit := 20
	offset := 10
	expectedResp := &dto.ListUsersResponse{
		Users:  []dto.UserResponse{},
		Total:  0,
		Limit:  limit,
		Offset: offset,
	}

	listUC.On("Execute", mock.Anything, limit, offset).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.GET("/users", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/users?limit=20&offset=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	listUC.AssertExpectations(t)
}

func TestHandler_List_InternalServerError(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	limit := 10
	offset := 0
	listUC.On("Execute", mock.Anything, limit, offset).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.GET("/users", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	listUC.AssertExpectations(t)
}

func TestHandler_Update_Success(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	userID := int64(1)
	reqBody := dto.UpdateUserRequest{
		Name:     "Updated User",
		Email:    "updated@example.com",
		Password: "newpassword123",
	}

	expectedResp := &dto.UserResponse{
		ID:        userID,
		Name:      reqBody.Name,
		Email:     reqBody.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	updateUC.On("Execute", mock.Anything, userID, reqBody).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.PUT("/users/:id", handler.Update)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	updateUC.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "User updated successfully", response["message"])
}

func TestHandler_Update_BadRequest_InvalidID(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	reqBody := dto.UpdateUserRequest{
		Name:  "Updated User",
		Email: "updated@example.com",
	}

	router := setupTestRouter(handler)
	router.PUT("/users/:id", handler.Update)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/users/invalid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	updateUC.AssertNotCalled(t, "Execute")
}

func TestHandler_Update_BadRequest_InvalidJSON(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	router := setupTestRouter(handler)
	router.PUT("/users/:id", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	updateUC.AssertNotCalled(t, "Execute")
}

func TestHandler_Update_NotFound(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	userID := int64(999)
	reqBody := dto.UpdateUserRequest{
		Name:  "Updated User",
		Email: "updated@example.com",
	}

	updateUC.On("Execute", mock.Anything, userID, reqBody).Return(nil, domainuser.ErrUserNotFound)

	router := setupTestRouter(handler)
	router.PUT("/users/:id", handler.Update)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/users/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	updateUC.AssertExpectations(t)
}

func TestHandler_Update_Conflict_EmailExists(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	userID := int64(1)
	reqBody := dto.UpdateUserRequest{
		Name:  "Updated User",
		Email: "existing@example.com",
	}

	updateUC.On("Execute", mock.Anything, userID, reqBody).Return(nil, domainuser.ErrEmailExists)

	router := setupTestRouter(handler)
	router.PUT("/users/:id", handler.Update)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	updateUC.AssertExpectations(t)
}

func TestHandler_Update_InternalServerError(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	userID := int64(1)
	reqBody := dto.UpdateUserRequest{
		Name:  "Updated User",
		Email: "updated@example.com",
	}

	updateUC.On("Execute", mock.Anything, userID, reqBody).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.PUT("/users/:id", handler.Update)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	updateUC.AssertExpectations(t)
}

func TestHandler_Delete_Success(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	userID := int64(1)
	deleteUC.On("Execute", mock.Anything, userID).Return(nil)

	router := setupTestRouter(handler)
	router.DELETE("/users/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	deleteUC.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "User deleted successfully", response["message"])
}

func TestHandler_Delete_BadRequest_InvalidID(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	router := setupTestRouter(handler)
	router.DELETE("/users/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/users/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	deleteUC.AssertNotCalled(t, "Execute")
}

func TestHandler_Delete_NotFound(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	userID := int64(999)
	deleteUC.On("Execute", mock.Anything, userID).Return(domainuser.ErrUserNotFound)

	router := setupTestRouter(handler)
	router.DELETE("/users/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/users/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	deleteUC.AssertExpectations(t)
}

func TestHandler_Delete_InternalServerError(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	userID := int64(1)
	deleteUC.On("Execute", mock.Anything, userID).Return(errors.New("database error"))

	router := setupTestRouter(handler)
	router.DELETE("/users/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	deleteUC.AssertExpectations(t)
}

func TestHandler_Register_Success(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	reqBody := dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedResp := &dto.UserResponse{
		ID:        1,
		Name:      reqBody.Name,
		Email:     reqBody.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createUC.On("Execute", mock.Anything, reqBody).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.POST("/users/register", handler.Register)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	createUC.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "User registered successfully", response["message"])
}

func TestHandler_Register_BadRequest_InvalidJSON(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	router := setupTestRouter(handler)
	router.POST("/users/register", handler.Register)

	req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	createUC.AssertNotCalled(t, "Execute")
}

func TestHandler_Register_Conflict_EmailExists(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	reqBody := dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	createUC.On("Execute", mock.Anything, reqBody).Return(nil, domainuser.ErrEmailExists)

	router := setupTestRouter(handler)
	router.POST("/users/register", handler.Register)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	createUC.AssertExpectations(t)
}

func TestHandler_Login_Success(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	reqBody := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedResp := &dto.LoginResponse{
		Token: "jwt_token_here",
		User: dto.UserResponse{
			ID:        1,
			Name:      "Test User",
			Email:     reqBody.Email,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	loginUC.On("Execute", mock.Anything, reqBody).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.POST("/users/login", handler.Login)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	loginUC.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Login successful", response["message"])
}

func TestHandler_Login_BadRequest_InvalidJSON(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	router := setupTestRouter(handler)
	router.POST("/users/login", handler.Login)

	req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	loginUC.AssertNotCalled(t, "Execute")
}

func TestHandler_Login_Unauthorized_InvalidCredentials(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	reqBody := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	loginUC.On("Execute", mock.Anything, reqBody).Return(nil, domainuser.ErrInvalidCredentials)

	router := setupTestRouter(handler)
	router.POST("/users/login", handler.Login)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	loginUC.AssertExpectations(t)
}

func TestHandler_Login_InternalServerError(t *testing.T) {
	createUC := &mockCreateUserUseCase{}
	getUC := &mockGetUserUseCase{}
	listUC := &mockListUsersUseCase{}
	updateUC := &mockUpdateUserUseCase{}
	deleteUC := &mockDeleteUserUseCase{}
	loginUC := &mockLoginUseCase{}

	handler := NewHandler(createUC, getUC, listUC, updateUC, deleteUC, loginUC)

	reqBody := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	loginUC.On("Execute", mock.Anything, reqBody).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.POST("/users/login", handler.Login)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	loginUC.AssertExpectations(t)
}
