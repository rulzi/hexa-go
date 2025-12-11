package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rulzi/hexa-go/internal/application/user"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	createUseCase *user.CreateUserUseCase
	getUseCase    *user.GetUserUseCase
	listUseCase   *user.ListUsersUseCase
	updateUseCase *user.UpdateUserUseCase
	deleteUseCase *user.DeleteUserUseCase
	loginUseCase  *user.LoginUseCase
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(
	createUseCase *user.CreateUserUseCase,
	getUseCase *user.GetUserUseCase,
	listUseCase *user.ListUsersUseCase,
	updateUseCase *user.UpdateUserUseCase,
	deleteUseCase *user.DeleteUserUseCase,
	loginUseCase *user.LoginUseCase,
) *UserHandler {
	return &UserHandler{
		createUseCase: createUseCase,
		getUseCase:    getUseCase,
		listUseCase:   listUseCase,
		updateUseCase: updateUseCase,
		deleteUseCase: deleteUseCase,
		loginUseCase:  loginUseCase,
	}
}

// Create handles POST /users
func (h *UserHandler) Create(c *gin.Context) {
	var req user.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponseBadRequest(c, err.Error())
		return
	}

	response, err := h.createUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		if err == domainuser.ErrEmailExists {
			ErrorResponseConflict(c, err.Error())
		} else {
			ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponseCreated(c, "User created successfully", response)
}

// Get handles GET /users/:id
func (h *UserHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ErrorResponseBadRequest(c, "invalid user id")
		return
	}

	response, err := h.getUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		if err == domainuser.ErrUserNotFound {
			ErrorResponseNotFound(c, err.Error())
		} else {
			ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponseOK(c, "User retrieved successfully", response)
}

// List handles GET /users
func (h *UserHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	response, err := h.listUseCase.Execute(c.Request.Context(), limit, offset)
	if err != nil {
		ErrorResponseInternalServerError(c, err.Error())
		return
	}

	SuccessResponseOK(c, "Users retrieved successfully", response)
}

// Update handles PUT /users/:id
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ErrorResponseBadRequest(c, "invalid user id")
		return
	}

	var req user.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponseBadRequest(c, err.Error())
		return
	}

	response, err := h.updateUseCase.Execute(c.Request.Context(), id, req)
	if err != nil {
		if err == domainuser.ErrUserNotFound {
			ErrorResponseNotFound(c, err.Error())
		} else if err == domainuser.ErrEmailExists {
			ErrorResponseConflict(c, err.Error())
		} else {
			ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponseOK(c, "User updated successfully", response)
}

// Delete handles DELETE /users/:id
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ErrorResponseBadRequest(c, "invalid user id")
		return
	}

	err = h.deleteUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		if err == domainuser.ErrUserNotFound {
			ErrorResponseNotFound(c, err.Error())
		} else {
			ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponseOK(c, "User deleted successfully", nil)
}

// Register handles POST /users/register
func (h *UserHandler) Register(c *gin.Context) {
	var req user.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponseBadRequest(c, err.Error())
		return
	}

	response, err := h.createUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		if err == domainuser.ErrEmailExists {
			ErrorResponseConflict(c, err.Error())
		} else {
			ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponseCreated(c, "User registered successfully", response)
}

// Login handles POST /users/login
func (h *UserHandler) Login(c *gin.Context) {
	var req user.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponseBadRequest(c, err.Error())
		return
	}

	response, err := h.loginUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		if err == domainuser.ErrInvalidCredentials {
			ErrorResponseUnauthorized(c, err.Error())
		} else {
			ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponseOK(c, "Login successful", response)
}
