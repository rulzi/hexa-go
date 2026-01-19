package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestStatusConstants(t *testing.T) {
	assert.Equal(t, Status("success"), StatusSuccess)
	assert.Equal(t, Status("error"), StatusError)
}

func TestHTTPStatusCodes_OK(t *testing.T) {
	assert.Equal(t, http.StatusOK, StatusCode.OK())
}

func TestHTTPStatusCodes_Created(t *testing.T) {
	assert.Equal(t, http.StatusCreated, StatusCode.Created())
}

func TestHTTPStatusCodes_Accepted(t *testing.T) {
	assert.Equal(t, http.StatusAccepted, StatusCode.Accepted())
}

func TestHTTPStatusCodes_NoContent(t *testing.T) {
	assert.Equal(t, http.StatusNoContent, StatusCode.NoContent())
}

func TestHTTPStatusCodes_BadRequest(t *testing.T) {
	assert.Equal(t, http.StatusBadRequest, StatusCode.BadRequest())
}

func TestHTTPStatusCodes_Unauthorized(t *testing.T) {
	assert.Equal(t, http.StatusUnauthorized, StatusCode.Unauthorized())
}

func TestHTTPStatusCodes_Forbidden(t *testing.T) {
	assert.Equal(t, http.StatusForbidden, StatusCode.Forbidden())
}

func TestHTTPStatusCodes_NotFound(t *testing.T) {
	assert.Equal(t, http.StatusNotFound, StatusCode.NotFound())
}

func TestHTTPStatusCodes_MethodNotAllowed(t *testing.T) {
	assert.Equal(t, http.StatusMethodNotAllowed, StatusCode.MethodNotAllowed())
}

func TestHTTPStatusCodes_Conflict(t *testing.T) {
	assert.Equal(t, http.StatusConflict, StatusCode.Conflict())
}

func TestHTTPStatusCodes_UnprocessableEntity(t *testing.T) {
	assert.Equal(t, http.StatusUnprocessableEntity, StatusCode.UnprocessableEntity())
}

func TestHTTPStatusCodes_TooManyRequests(t *testing.T) {
	assert.Equal(t, http.StatusTooManyRequests, StatusCode.TooManyRequests())
}

func TestHTTPStatusCodes_InternalServerError(t *testing.T) {
	assert.Equal(t, http.StatusInternalServerError, StatusCode.InternalServerError())
}

func TestHTTPStatusCodes_BadGateway(t *testing.T) {
	assert.Equal(t, http.StatusBadGateway, StatusCode.BadGateway())
}

func TestHTTPStatusCodes_ServiceUnavailable(t *testing.T) {
	assert.Equal(t, http.StatusServiceUnavailable, StatusCode.ServiceUnavailable())
}

func TestSuccessResponse(t *testing.T) {
	c, w := setupTestContext()
	
	data := map[string]string{"key": "value"}
	SuccessResponse(c, http.StatusOK, "Success message", data)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response StandardResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, StatusSuccess, response.Status)
	assert.Equal(t, "Success message", response.Message)
	assert.NotNil(t, response.Data)
}

func TestSuccessResponse_WithNilData(t *testing.T) {
	c, w := setupTestContext()
	
	SuccessResponse(c, http.StatusOK, "Success message", nil)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response StandardResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, StatusSuccess, response.Status)
	assert.Equal(t, "Success message", response.Message)
}

func TestSuccessResponse_WithComplexData(t *testing.T) {
	c, w := setupTestContext()
	
	data := map[string]interface{}{
		"id":    1,
		"name":  "Test",
		"items": []string{"item1", "item2"},
	}
	SuccessResponse(c, http.StatusCreated, "Created", data)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response StandardResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, StatusSuccess, response.Status)
	assert.Equal(t, "Created", response.Message)
}

func TestErrorResponse(t *testing.T) {
	c, w := setupTestContext()
	
	ErrorResponse(c, http.StatusBadRequest, "Error message")
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response StandardResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, StatusError, response.Status)
	assert.Equal(t, "Error message", response.Message)
	assert.Nil(t, response.Data)
}

func TestErrorResponse_WithDifferentStatusCodes(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		message    string
	}{
		{"BadRequest", http.StatusBadRequest, "Bad request"},
		{"Unauthorized", http.StatusUnauthorized, "Unauthorized"},
		{"NotFound", http.StatusNotFound, "Not found"},
		{"InternalServerError", http.StatusInternalServerError, "Internal server error"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, w := setupTestContext()
			
			ErrorResponse(c, tc.statusCode, tc.message)
			
			assert.Equal(t, tc.statusCode, w.Code)
			
			var response StandardResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, StatusError, response.Status)
			assert.Equal(t, tc.message, response.Message)
		})
	}
}

func TestSuccessResponseOK(t *testing.T) {
	c, w := setupTestContext()
	
	data := map[string]string{"key": "value"}
	SuccessResponseOK(c, "OK message", data)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response StandardResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, StatusSuccess, response.Status)
	assert.Equal(t, "OK message", response.Message)
}

func TestSuccessResponseCreated(t *testing.T) {
	c, w := setupTestContext()
	
	data := map[string]string{"id": "1"}
	SuccessResponseCreated(c, "Created message", data)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response StandardResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, StatusSuccess, response.Status)
	assert.Equal(t, "Created message", response.Message)
}

func TestErrorResponseBadRequest(t *testing.T) {
	c, w := setupTestContext()
	
	ErrorResponseBadRequest(c, "Bad request message")
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response StandardResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, StatusError, response.Status)
	assert.Equal(t, "Bad request message", response.Message)
}

func TestErrorResponseUnauthorized(t *testing.T) {
	c, w := setupTestContext()
	
	ErrorResponseUnauthorized(c, "Unauthorized message")
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	var response StandardResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, StatusError, response.Status)
	assert.Equal(t, "Unauthorized message", response.Message)
}

func TestErrorResponseNotFound(t *testing.T) {
	c, w := setupTestContext()
	
	ErrorResponseNotFound(c, "Not found message")
	
	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var response StandardResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, StatusError, response.Status)
	assert.Equal(t, "Not found message", response.Message)
}

func TestErrorResponseConflict(t *testing.T) {
	c, w := setupTestContext()
	
	ErrorResponseConflict(c, "Conflict message")
	
	assert.Equal(t, http.StatusConflict, w.Code)
	
	var response StandardResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, StatusError, response.Status)
	assert.Equal(t, "Conflict message", response.Message)
}

func TestErrorResponseInternalServerError(t *testing.T) {
	c, w := setupTestContext()
	
	ErrorResponseInternalServerError(c, "Internal server error message")
	
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	
	var response StandardResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, StatusError, response.Status)
	assert.Equal(t, "Internal server error message", response.Message)
}

func TestStandardResponse_JSONStructure(t *testing.T) {
	c, w := setupTestContext()
	
	SuccessResponse(c, http.StatusOK, "Test message", map[string]string{"key": "value"})
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// Verify JSON structure
	assert.Contains(t, response, "status")
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "data")
	
	assert.Equal(t, "success", response["status"])
	assert.Equal(t, "Test message", response["message"])
}

func TestStandardResponse_ErrorJSONStructure(t *testing.T) {
	c, w := setupTestContext()
	
	ErrorResponse(c, http.StatusBadRequest, "Error message")
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// Verify JSON structure
	assert.Contains(t, response, "status")
	assert.Contains(t, response, "message")
	
	assert.Equal(t, "error", response["status"])
	assert.Equal(t, "Error message", response["message"])
	
	// Data should be nil or omitted for errors
	if data, exists := response["data"]; exists {
		assert.Nil(t, data)
	}
}

func TestStandardResponse_EmptyMessage(t *testing.T) {
	c, w := setupTestContext()
	
	SuccessResponse(c, http.StatusOK, "", map[string]string{"key": "value"})
	
	var response StandardResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, StatusSuccess, response.Status)
	assert.Empty(t, response.Message)
}

func TestStandardResponse_EmptyData(t *testing.T) {
	c, w := setupTestContext()
	
	SuccessResponse(c, http.StatusOK, "Message", "")
	
	var response StandardResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, StatusSuccess, response.Status)
	assert.Equal(t, "Message", response.Message)
}

