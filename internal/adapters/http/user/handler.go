package user

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rulzi/hexa-go/internal/adapters/http/response"
	"github.com/rulzi/hexa-go/internal/application/user/dto"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
)

// UseCase is the interface for user operations
type UseCase interface {
	Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
	Get(ctx context.Context, id int64) (*dto.UserResponse, error)
	List(ctx context.Context, limit, offset int) (*dto.ListUsersResponse, error)
	Update(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(ctx context.Context, id int64) error
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
}

// Handler handles HTTP requests for users
type Handler struct {
	userUseCase UseCase
}

// NewHandler creates a new Handler
func NewHandler(userUseCase UseCase) *Handler {
	return &Handler{
		userUseCase: userUseCase,
	}
}

// Create handles POST /users
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponseBadRequest(c, err.Error())
		return
	}

	resp, err := h.userUseCase.Create(c.Request.Context(), req)
	if err != nil {
		if err == domainuser.ErrEmailExists {
			response.ErrorResponseConflict(c, err.Error())
		} else {
			response.ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	response.SuccessResponseCreated(c, "User created successfully", resp)
}

// Get handles GET /users/:id
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorResponseBadRequest(c, "invalid user id")
		return
	}

	resp, err := h.userUseCase.Get(c.Request.Context(), id)
	if err != nil {
		if err == domainuser.ErrUserNotFound {
			response.ErrorResponseNotFound(c, err.Error())
		} else {
			response.ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	response.SuccessResponseOK(c, "User retrieved successfully", resp)
}

// List handles GET /users
func (h *Handler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := h.userUseCase.List(c.Request.Context(), limit, offset)
	if err != nil {
		response.ErrorResponseInternalServerError(c, err.Error())
		return
	}

	response.SuccessResponseOK(c, "Users retrieved successfully", resp)
}

// Update handles PUT /users/:id
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorResponseBadRequest(c, "invalid user id")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponseBadRequest(c, err.Error())
		return
	}

	resp, err := h.userUseCase.Update(c.Request.Context(), id, req)
	if err != nil {
		switch err {
		case domainuser.ErrUserNotFound:
			response.ErrorResponseNotFound(c, err.Error())
		case domainuser.ErrEmailExists:
			response.ErrorResponseConflict(c, err.Error())
		default:
			response.ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	response.SuccessResponseOK(c, "User updated successfully", resp)
}

// Delete handles DELETE /users/:id
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorResponseBadRequest(c, "invalid user id")
		return
	}

	err = h.userUseCase.Delete(c.Request.Context(), id)
	if err != nil {
		if err == domainuser.ErrUserNotFound {
			response.ErrorResponseNotFound(c, err.Error())
		} else {
			response.ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	response.SuccessResponseOK(c, "User deleted successfully", nil)
}

// Register handles POST /users/register
func (h *Handler) Register(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponseBadRequest(c, err.Error())
		return
	}

	resp, err := h.userUseCase.Create(c.Request.Context(), req)
	if err != nil {
		if err == domainuser.ErrEmailExists {
			response.ErrorResponseConflict(c, err.Error())
		} else {
			response.ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	response.SuccessResponseCreated(c, "User registered successfully", resp)
}

// Login handles POST /users/login
func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponseBadRequest(c, err.Error())
		return
	}

	resp, err := h.userUseCase.Login(c.Request.Context(), req)
	if err != nil {
		if err == domainuser.ErrInvalidCredentials {
			response.ErrorResponseUnauthorized(c, err.Error())
		} else {
			response.ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	response.SuccessResponseOK(c, "Login successful", resp)
}
