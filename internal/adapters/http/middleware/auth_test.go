package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockTokenValidator is a mock implementation of TokenValidator
type mockTokenValidator struct {
	mock.Mock
}

func (m *mockTokenValidator) Validate(token string) (*domainuser.TokenClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainuser.TokenClaims), args.Error(1)
}

func setupTestRouter(middleware gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	return router
}

func TestAuthMiddleware_Success(t *testing.T) {
	mockValidator := &mockTokenValidator{}
	middleware := AuthMiddleware(mockValidator)

	expectedClaims := &domainuser.TokenClaims{
		UserID: 1,
		Email:  "test@example.com",
	}

	mockValidator.On("Validate", "valid-token").Return(expectedClaims, nil)

	router := setupTestRouter(middleware)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockValidator.AssertExpectations(t)

	// Verify that user info is set in context
	// We can't directly access context in the test, but we can verify the response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["message"])
}

func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	mockValidator := &mockTokenValidator{}
	middleware := AuthMiddleware(mockValidator)

	router := setupTestRouter(middleware)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// No Authorization header
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockValidator.AssertNotCalled(t, "Validate")

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "authorization header is required", response["message"])
}

func TestAuthMiddleware_EmptyAuthorizationHeader(t *testing.T) {
	mockValidator := &mockTokenValidator{}
	middleware := AuthMiddleware(mockValidator)

	router := setupTestRouter(middleware)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockValidator.AssertNotCalled(t, "Validate")

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "authorization header is required", response["message"])
}

func TestAuthMiddleware_InvalidFormat_NoBearer(t *testing.T) {
	mockValidator := &mockTokenValidator{}
	middleware := AuthMiddleware(mockValidator)

	router := setupTestRouter(middleware)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockValidator.AssertNotCalled(t, "Validate")

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "invalid authorization header format", response["message"])
}

func TestAuthMiddleware_InvalidFormat_NoSpace(t *testing.T) {
	mockValidator := &mockTokenValidator{}
	middleware := AuthMiddleware(mockValidator)

	router := setupTestRouter(middleware)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearertoken")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockValidator.AssertNotCalled(t, "Validate")

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "invalid authorization header format", response["message"])
}

func TestAuthMiddleware_InvalidFormat_TooManyParts(t *testing.T) {
	mockValidator := &mockTokenValidator{}
	middleware := AuthMiddleware(mockValidator)

	router := setupTestRouter(middleware)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer token extra")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockValidator.AssertNotCalled(t, "Validate")

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "invalid authorization header format", response["message"])
}

func TestAuthMiddleware_InvalidFormat_WrongPrefix(t *testing.T) {
	mockValidator := &mockTokenValidator{}
	middleware := AuthMiddleware(mockValidator)

	router := setupTestRouter(middleware)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Basic token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockValidator.AssertNotCalled(t, "Validate")

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "invalid authorization header format", response["message"])
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	mockValidator := &mockTokenValidator{}
	middleware := AuthMiddleware(mockValidator)

	mockValidator.On("Validate", "invalid-token").Return(nil, errors.New("token expired"))

	router := setupTestRouter(middleware)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockValidator.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "invalid or expired token", response["message"])
}

func TestAuthMiddleware_ContextValuesSet(t *testing.T) {
	mockValidator := &mockTokenValidator{}
	middleware := AuthMiddleware(mockValidator)

	expectedClaims := &domainuser.TokenClaims{
		UserID: 123,
		Email:  "user@example.com",
	}

	mockValidator.On("Validate", "valid-token").Return(expectedClaims, nil)

	// Create a custom handler to verify context values
	router := gin.New()
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		assert.True(t, exists, "user_id should be set in context")
		assert.Equal(t, int64(123), userID)

		userEmail, exists := c.Get("user_email")
		assert.True(t, exists, "user_email should be set in context")
		assert.Equal(t, "user@example.com", userEmail)

		c.JSON(http.StatusOK, gin.H{
			"user_id":    userID,
			"user_email": userEmail,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockValidator.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(123), response["user_id"])
	assert.Equal(t, "user@example.com", response["user_email"])
}

func TestAuthMiddleware_AbortsOnError(t *testing.T) {
	mockValidator := &mockTokenValidator{}
	middleware := AuthMiddleware(mockValidator)

	router := gin.New()
	router.Use(middleware)
	// This handler should not be called if middleware aborts
	router.GET("/test", func(c *gin.Context) {
		t.Error("Handler should not be called when middleware aborts")
		c.JSON(http.StatusOK, gin.H{"message": "should not reach here"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// No Authorization header - should abort
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockValidator.AssertNotCalled(t, "Validate")

	// Verify the handler was not called by checking response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])
	assert.NotEqual(t, "should not reach here", response["message"])
}

func TestAuthMiddleware_CaseSensitiveBearer(t *testing.T) {
	mockValidator := &mockTokenValidator{}
	middleware := AuthMiddleware(mockValidator)

	router := setupTestRouter(middleware)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "bearer token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockValidator.AssertNotCalled(t, "Validate")

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "invalid authorization header format", response["message"])
}

func TestAuthMiddleware_EmptyToken(t *testing.T) {
	mockValidator := &mockTokenValidator{}
	middleware := AuthMiddleware(mockValidator)

	// Empty token should still call Validate with empty string
	mockValidator.On("Validate", "").Return(nil, errors.New("empty token"))

	router := setupTestRouter(middleware)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer ")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockValidator.AssertExpectations(t)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "invalid or expired token", response["message"])
}

