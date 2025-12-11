package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rulzi/hexa-go/internal/adapters/http/response"
	"github.com/rulzi/hexa-go/internal/application/user/dto"
	"github.com/rulzi/hexa-go/internal/application/user/usecase"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
)

// Handler handles HTTP requests for users
type Handler struct {
	createUseCase *usecase.CreateUserUseCase
	getUseCase    *usecase.GetUserUseCase
	listUseCase   *usecase.ListUsersUseCase
	updateUseCase *usecase.UpdateUserUseCase
	deleteUseCase *usecase.DeleteUserUseCase
	loginUseCase  *usecase.LoginUseCase
}

// NewHandler creates a new Handler
func NewHandler(
	createUseCase *usecase.CreateUserUseCase,
	getUseCase *usecase.GetUserUseCase,
	listUseCase *usecase.ListUsersUseCase,
	updateUseCase *usecase.UpdateUserUseCase,
	deleteUseCase *usecase.DeleteUserUseCase,
	loginUseCase *usecase.LoginUseCase,
) *Handler {
	return &Handler{
		createUseCase: createUseCase,
		getUseCase:    getUseCase,
		listUseCase:   listUseCase,
		updateUseCase: updateUseCase,
		deleteUseCase: deleteUseCase,
		loginUseCase:  loginUseCase,
	}
}

// Create handles POST /users
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponseBadRequest(c, err.Error())
		return
	}

	resp, err := h.createUseCase.Execute(c.Request.Context(), req)
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

	resp, err := h.getUseCase.Execute(c.Request.Context(), id)
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

	resp, err := h.listUseCase.Execute(c.Request.Context(), limit, offset)
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

	resp, err := h.updateUseCase.Execute(c.Request.Context(), id, req)
	if err != nil {
		if err == domainuser.ErrUserNotFound {
			response.ErrorResponseNotFound(c, err.Error())
		} else if err == domainuser.ErrEmailExists {
			response.ErrorResponseConflict(c, err.Error())
		} else {
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

	err = h.deleteUseCase.Execute(c.Request.Context(), id)
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

	resp, err := h.createUseCase.Execute(c.Request.Context(), req)
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

	resp, err := h.loginUseCase.Execute(c.Request.Context(), req)
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
