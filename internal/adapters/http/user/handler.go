package user

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rulzi/hexa-go/internal/adapters/http/response"
	"github.com/rulzi/hexa-go/internal/application/user/dto"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
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
	logger      logger.Logger
}

// NewHandler creates a new Handler
func NewHandler(userUseCase UseCase, appLogger logger.Logger) *Handler {
	return &Handler{
		userUseCase: userUseCase,
		logger:      appLogger,
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
			h.logger.WarnWithFields("user create conflict", map[string]interface{}{"email": req.Email})
			response.ErrorResponseConflict(c, err.Error())
		} else {
			h.logger.ErrorWithFields("user create failed", map[string]interface{}{"error": err.Error()})
			response.ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	h.logger.InfoWithFields("user created", map[string]interface{}{"user_id": resp.ID, "email": resp.Email})
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
			h.logger.DebugWithFields("user not found", map[string]interface{}{"user_id": id})
			response.ErrorResponseNotFound(c, err.Error())
		} else {
			h.logger.ErrorWithFields("user get failed", map[string]interface{}{"user_id": id, "error": err.Error()})
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
		h.logger.ErrorWithFields("user list failed", map[string]interface{}{"error": err.Error()})
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
			h.logger.DebugWithFields("user not found on update", map[string]interface{}{"user_id": id})
			response.ErrorResponseNotFound(c, err.Error())
		case domainuser.ErrEmailExists:
			h.logger.WarnWithFields("user update conflict", map[string]interface{}{"user_id": id, "email": req.Email})
			response.ErrorResponseConflict(c, err.Error())
		default:
			h.logger.ErrorWithFields("user update failed", map[string]interface{}{"user_id": id, "error": err.Error()})
			response.ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	h.logger.InfoWithFields("user updated", map[string]interface{}{"user_id": id})
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
			h.logger.DebugWithFields("user not found on delete", map[string]interface{}{"user_id": id})
			response.ErrorResponseNotFound(c, err.Error())
		} else {
			h.logger.ErrorWithFields("user delete failed", map[string]interface{}{"user_id": id, "error": err.Error()})
			response.ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	h.logger.InfoWithFields("user deleted", map[string]interface{}{"user_id": id})
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
			h.logger.WarnWithFields("register conflict", map[string]interface{}{"email": req.Email})
			response.ErrorResponseConflict(c, err.Error())
		} else {
			h.logger.ErrorWithFields("register failed", map[string]interface{}{"error": err.Error()})
			response.ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	h.logger.InfoWithFields("user registered", map[string]interface{}{"user_id": resp.ID, "email": resp.Email})
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
			h.logger.WarnWithFields("login failed invalid credentials", map[string]interface{}{"email": req.Email})
			response.ErrorResponseUnauthorized(c, err.Error())
		} else {
			h.logger.ErrorWithFields("login failed", map[string]interface{}{"email": req.Email, "error": err.Error()})
			response.ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	h.logger.InfoWithFields("login successful", map[string]interface{}{"user_id": resp.User.ID, "email": req.Email})
	response.SuccessResponseOK(c, "Login successful", resp)
}
