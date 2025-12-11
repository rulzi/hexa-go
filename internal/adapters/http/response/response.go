package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseStatus represents the status in the response
type ResponseStatus string

const (
	StatusSuccess ResponseStatus = "success"
	StatusError   ResponseStatus = "error"
)

// StandardResponse represents the standard API response structure
type StandardResponse struct {
	Status  ResponseStatus `json:"status"`
	Message string         `json:"message"`
	Data    interface{}    `json:"data,omitempty"`
}

// HTTPStatusCodes contains standardized HTTP status codes
type HTTPStatusCodes struct{}

var StatusCode = HTTPStatusCodes{}

// Success status codes
func (HTTPStatusCodes) OK() int        { return http.StatusOK }        // 200
func (HTTPStatusCodes) Created() int   { return http.StatusCreated }   // 201
func (HTTPStatusCodes) Accepted() int  { return http.StatusAccepted }  // 202
func (HTTPStatusCodes) NoContent() int { return http.StatusNoContent } // 204

// Client error status codes
func (HTTPStatusCodes) BadRequest() int          { return http.StatusBadRequest }          // 400
func (HTTPStatusCodes) Unauthorized() int        { return http.StatusUnauthorized }        // 401
func (HTTPStatusCodes) Forbidden() int           { return http.StatusForbidden }           // 403
func (HTTPStatusCodes) NotFound() int            { return http.StatusNotFound }            // 404
func (HTTPStatusCodes) MethodNotAllowed() int    { return http.StatusMethodNotAllowed }    // 405
func (HTTPStatusCodes) Conflict() int            { return http.StatusConflict }            // 409
func (HTTPStatusCodes) UnprocessableEntity() int { return http.StatusUnprocessableEntity } // 422
func (HTTPStatusCodes) TooManyRequests() int     { return http.StatusTooManyRequests }     // 429

// Server error status codes
func (HTTPStatusCodes) InternalServerError() int { return http.StatusInternalServerError } // 500
func (HTTPStatusCodes) BadGateway() int          { return http.StatusBadGateway }          // 502
func (HTTPStatusCodes) ServiceUnavailable() int  { return http.StatusServiceUnavailable }  // 503

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

