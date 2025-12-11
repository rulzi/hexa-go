package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Status represents the status in the response
type Status string

const (
	// StatusSuccess represents a successful response status
	StatusSuccess Status = "success"
	// StatusError represents an error response status
	StatusError Status = "error"
)

// StandardResponse represents the standard API response structure
type StandardResponse struct {
	Status  Status      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// HTTPStatusCodes contains standardized HTTP status codes
type HTTPStatusCodes struct{}

// StatusCode provides access to HTTP status code methods
var StatusCode = HTTPStatusCodes{}

// Success status codes

// OK returns HTTP 200 status code
func (HTTPStatusCodes) OK() int { return http.StatusOK } // 200

// Created returns HTTP 201 status code
func (HTTPStatusCodes) Created() int { return http.StatusCreated } // 201

// Accepted returns HTTP 202 status code
func (HTTPStatusCodes) Accepted() int { return http.StatusAccepted } // 202

// NoContent returns HTTP 204 status code
func (HTTPStatusCodes) NoContent() int { return http.StatusNoContent } // 204

// Client error status codes

// BadRequest returns HTTP 400 status code
func (HTTPStatusCodes) BadRequest() int { return http.StatusBadRequest } // 400

// Unauthorized returns HTTP 401 status code
func (HTTPStatusCodes) Unauthorized() int { return http.StatusUnauthorized } // 401

// Forbidden returns HTTP 403 status code
func (HTTPStatusCodes) Forbidden() int { return http.StatusForbidden } // 403

// NotFound returns HTTP 404 status code
func (HTTPStatusCodes) NotFound() int { return http.StatusNotFound } // 404

// MethodNotAllowed returns HTTP 405 status code
func (HTTPStatusCodes) MethodNotAllowed() int { return http.StatusMethodNotAllowed } // 405

// Conflict returns HTTP 409 status code
func (HTTPStatusCodes) Conflict() int { return http.StatusConflict } // 409

// UnprocessableEntity returns HTTP 422 status code
func (HTTPStatusCodes) UnprocessableEntity() int { return http.StatusUnprocessableEntity } // 422

// TooManyRequests returns HTTP 429 status code
func (HTTPStatusCodes) TooManyRequests() int { return http.StatusTooManyRequests } // 429

// Server error status codes

// InternalServerError returns HTTP 500 status code
func (HTTPStatusCodes) InternalServerError() int { return http.StatusInternalServerError } // 500

// BadGateway returns HTTP 502 status code
func (HTTPStatusCodes) BadGateway() int { return http.StatusBadGateway } // 502

// ServiceUnavailable returns HTTP 503 status code
func (HTTPStatusCodes) ServiceUnavailable() int { return http.StatusServiceUnavailable } // 503

// SuccessResponse sends a successful response with data
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, StandardResponse{
		Status:  StatusSuccess,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, StandardResponse{
		Status:  StatusError,
		Message: message,
		Data:    nil,
	})
}

// SuccessResponseOK sends a 200 OK success response
func SuccessResponseOK(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, StatusCode.OK(), message, data)
}

// SuccessResponseCreated sends a 201 Created success response
func SuccessResponseCreated(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, StatusCode.Created(), message, data)
}

// ErrorResponseBadRequest sends a 400 Bad Request error response
func ErrorResponseBadRequest(c *gin.Context, message string) {
	ErrorResponse(c, StatusCode.BadRequest(), message)
}

// ErrorResponseUnauthorized sends a 401 Unauthorized error response
func ErrorResponseUnauthorized(c *gin.Context, message string) {
	ErrorResponse(c, StatusCode.Unauthorized(), message)
}

// ErrorResponseNotFound sends a 404 Not Found error response
func ErrorResponseNotFound(c *gin.Context, message string) {
	ErrorResponse(c, StatusCode.NotFound(), message)
}

// ErrorResponseConflict sends a 409 Conflict error response
func ErrorResponseConflict(c *gin.Context, message string) {
	ErrorResponse(c, StatusCode.Conflict(), message)
}

// ErrorResponseInternalServerError sends a 500 Internal Server Error response
func ErrorResponseInternalServerError(c *gin.Context, message string) {
	ErrorResponse(c, StatusCode.InternalServerError(), message)
}
