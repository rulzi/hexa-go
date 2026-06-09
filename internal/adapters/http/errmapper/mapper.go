package errmapper

import (
	"github.com/gin-gonic/gin"
	"github.com/rulzi/hexa-go/internal/adapters/http/response"
	domainerrs "github.com/rulzi/hexa-go/internal/domain/errs"
)

// StatusCode maps a domain error to the appropriate HTTP status code.
// Unknown errors default to 500 — infrastructure errors stay opaque to clients.
func StatusCode(err error) int {
	switch {
	case domainerrs.IsValidation(err):
		return response.StatusCode.BadRequest()
	case domainerrs.IsUnauthorized(err):
		return response.StatusCode.Unauthorized()
	case domainerrs.IsNotFound(err):
		return response.StatusCode.NotFound()
	case domainerrs.IsConflict(err):
		return response.StatusCode.Conflict()
	default:
		return response.StatusCode.InternalServerError()
	}
}

// Respond writes a standardized error response for the given domain/infrastructure error.
func Respond(c *gin.Context, err error) {
	response.ErrorResponse(c, StatusCode(err), domainerrs.Message(err))
}
