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
	userentity "github.com/rulzi/hexa-go/internal/domain/user/entity"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testLogger logger.Logger = logger.NewSimpleLogger()

// mockUserUseCase is a mock implementation of UserUseCase
type mockUserUseCase struct {
	mock.Mock
}

func (m *mockUserUseCase) Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *mockUserUseCase) Get(ctx context.Context, id int64) (*dto.UserResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *mockUserUseCase) List(ctx context.Context, limit, offset int) (*dto.ListUsersResponse, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ListUsersResponse), args.Error(1)
}

func (m *mockUserUseCase) Update(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *mockUserUseCase) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockUserUseCase) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
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
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)
	assert.NotNil(t, handler)
	assert.Equal(t, uc, handler.userUseCase)
}

func TestHandler_Create_Success(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

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

	uc.On("Create", mock.Anything, reqBody).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.POST("/users", handler.Create)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	uc.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "User created successfully", response["message"])
}

func TestHandler_Create_BadRequest_InvalidJSON(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	router := setupTestRouter(handler)
	router.POST("/users", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "Create")
}

func TestHandler_Create_Conflict_EmailExists(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	reqBody := dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	uc.On("Create", mock.Anything, reqBody).Return(nil, userentity.NewEmailExists())

	router := setupTestRouter(handler)
	router.POST("/users", handler.Create)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Create_InternalServerError(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	reqBody := dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	uc.On("Create", mock.Anything, reqBody).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.POST("/users", handler.Create)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Get_Success(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	userID := int64(1)
	expectedResp := &dto.UserResponse{
		ID:        userID,
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	uc.On("Get", mock.Anything, userID).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.GET("/users/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "User retrieved successfully", response["message"])
}

func TestHandler_Get_BadRequest_InvalidID(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	router := setupTestRouter(handler)
	router.GET("/users/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/users/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "Get")
}

func TestHandler_Get_NotFound(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	userID := int64(999)
	uc.On("Get", mock.Anything, userID).Return(nil, userentity.NewUserNotFound())

	router := setupTestRouter(handler)
	router.GET("/users/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Get_InternalServerError(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	userID := int64(1)
	uc.On("Get", mock.Anything, userID).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.GET("/users/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_List_Success(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

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

	uc.On("List", mock.Anything, limit, offset).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.GET("/users", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Users retrieved successfully", response["message"])
}

func TestHandler_List_WithQueryParams(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	limit := 20
	offset := 10
	expectedResp := &dto.ListUsersResponse{
		Users:  []dto.UserResponse{},
		Total:  0,
		Limit:  limit,
		Offset: offset,
	}

	uc.On("List", mock.Anything, limit, offset).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.GET("/users", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/users?limit=20&offset=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_List_InternalServerError(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	limit := 10
	offset := 0
	uc.On("List", mock.Anything, limit, offset).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.GET("/users", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Update_Success(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

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

	uc.On("Update", mock.Anything, userID, reqBody).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.PUT("/users/:id", handler.Update)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "User updated successfully", response["message"])
}

func TestHandler_Update_BadRequest_InvalidID(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

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
	uc.AssertNotCalled(t, "Update")
}

func TestHandler_Update_BadRequest_InvalidJSON(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	router := setupTestRouter(handler)
	router.PUT("/users/:id", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "Update")
}

func TestHandler_Update_NotFound(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	userID := int64(999)
	reqBody := dto.UpdateUserRequest{
		Name:  "Updated User",
		Email: "updated@example.com",
	}

	uc.On("Update", mock.Anything, userID, reqBody).Return(nil, userentity.NewUserNotFound())

	router := setupTestRouter(handler)
	router.PUT("/users/:id", handler.Update)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/users/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Update_Conflict_EmailExists(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	userID := int64(1)
	reqBody := dto.UpdateUserRequest{
		Name:  "Updated User",
		Email: "existing@example.com",
	}

	uc.On("Update", mock.Anything, userID, reqBody).Return(nil, userentity.NewEmailExists())

	router := setupTestRouter(handler)
	router.PUT("/users/:id", handler.Update)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Update_InternalServerError(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	userID := int64(1)
	reqBody := dto.UpdateUserRequest{
		Name:  "Updated User",
		Email: "updated@example.com",
	}

	uc.On("Update", mock.Anything, userID, reqBody).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.PUT("/users/:id", handler.Update)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Delete_Success(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	userID := int64(1)
	uc.On("Delete", mock.Anything, userID).Return(nil)

	router := setupTestRouter(handler)
	router.DELETE("/users/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "User deleted successfully", response["message"])
}

func TestHandler_Delete_BadRequest_InvalidID(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	router := setupTestRouter(handler)
	router.DELETE("/users/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/users/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "Delete")
}

func TestHandler_Delete_NotFound(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	userID := int64(999)
	uc.On("Delete", mock.Anything, userID).Return(userentity.NewUserNotFound())

	router := setupTestRouter(handler)
	router.DELETE("/users/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/users/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Delete_InternalServerError(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	userID := int64(1)
	uc.On("Delete", mock.Anything, userID).Return(errors.New("database error"))

	router := setupTestRouter(handler)
	router.DELETE("/users/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Register_Success(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

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

	uc.On("Create", mock.Anything, reqBody).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.POST("/users/register", handler.Register)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	uc.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "User registered successfully", response["message"])
}

func TestHandler_Register_BadRequest_InvalidJSON(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	router := setupTestRouter(handler)
	router.POST("/users/register", handler.Register)

	req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "Create")
}

func TestHandler_Register_Conflict_EmailExists(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	reqBody := dto.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	uc.On("Create", mock.Anything, reqBody).Return(nil, userentity.NewEmailExists())

	router := setupTestRouter(handler)
	router.POST("/users/register", handler.Register)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Login_Success(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

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

	uc.On("Login", mock.Anything, reqBody).Return(expectedResp, nil)

	router := setupTestRouter(handler)
	router.POST("/users/login", handler.Login)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Login successful", response["message"])
}

func TestHandler_Login_BadRequest_InvalidJSON(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	router := setupTestRouter(handler)
	router.POST("/users/login", handler.Login)

	req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "Login")
}

func TestHandler_Login_Unauthorized_InvalidCredentials(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	reqBody := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	uc.On("Login", mock.Anything, reqBody).Return(nil, userentity.NewInvalidCredentials())

	router := setupTestRouter(handler)
	router.POST("/users/login", handler.Login)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	uc.AssertExpectations(t)
}

func TestHandler_Login_InternalServerError(t *testing.T) {
	uc := &mockUserUseCase{}
	handler := NewHandler(uc, testLogger)

	reqBody := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	uc.On("Login", mock.Anything, reqBody).Return(nil, errors.New("database error"))

	router := setupTestRouter(handler)
	router.POST("/users/login", handler.Login)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	uc.AssertExpectations(t)
}
