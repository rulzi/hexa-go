package user

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rulzi/hexa-go/internal/adapters/http/errmapper"
	"github.com/rulzi/hexa-go/internal/adapters/http/response"
	"github.com/rulzi/hexa-go/internal/application/user/dto"
	domainerrs "github.com/rulzi/hexa-go/internal/domain/errs"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
	"github.com/rulzi/hexa-go/internal/pkg/pagination"
)

type createUseCase interface {
	Execute(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)
}

type getUseCase interface {
	Execute(ctx context.Context, id int64) (*dto.UserResponse, error)
}

type listUseCase interface {
	Execute(ctx context.Context, limit, offset int) (*dto.ListUsersResponse, error)
}

type updateUseCase interface {
	Execute(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error)
}

type deleteUseCase interface {
	Execute(ctx context.Context, id int64) error
}

type loginUseCase interface {
	Execute(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
}

// Deps holds user use case operations
type Deps struct {
	Create createUseCase
	Get    getUseCase
	List   listUseCase
	Update updateUseCase
	Delete deleteUseCase
	Login  loginUseCase
}

// Handler handles HTTP requests for users
type Handler struct {
	deps   Deps
	logger logger.Logger
}

// NewHandler creates a new Handler
func NewHandler(deps Deps, appLogger logger.Logger) *Handler {
	return &Handler{
		deps:   deps,
		logger: appLogger,
	}
}

func (h *Handler) handleUseCaseError(c *gin.Context, err error, fields map[string]interface{}, logMsg string) {
	if fields == nil {
		fields = make(map[string]interface{})
	}

	if domainerrs.IsNotFound(err) || domainerrs.IsValidation(err) || domainerrs.IsConflict(err) || domainerrs.IsUnauthorized(err) {
		h.logger.DebugWithFields(logMsg, fields)
	} else {
		fields["error"] = err.Error()
		h.logger.ErrorWithFields(logMsg, fields)
	}
	errmapper.Respond(c, err)
}

// Create handles POST /users
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponseBadRequest(c, err.Error())
		return
	}

	resp, err := h.deps.Create.Execute(c.Request.Context(), req)
	if err != nil {
		h.handleUseCaseError(c, err, map[string]interface{}{"email": req.Email}, "user create failed")
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

	resp, err := h.deps.Get.Execute(c.Request.Context(), id)
	if err != nil {
		h.handleUseCaseError(c, err, map[string]interface{}{"user_id": id}, "user get failed")
		return
	}

	response.SuccessResponseOK(c, "User retrieved successfully", resp)
}

// List handles GET /users
func (h *Handler) List(c *gin.Context) {
	limit, offset := pagination.Parse(c.DefaultQuery("limit", "10"), c.DefaultQuery("offset", "0"))

	resp, err := h.deps.List.Execute(c.Request.Context(), limit, offset)
	if err != nil {
		h.handleUseCaseError(c, err, nil, "user list failed")
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

	resp, err := h.deps.Update.Execute(c.Request.Context(), id, req)
	if err != nil {
		h.handleUseCaseError(c, err, map[string]interface{}{"user_id": id, "email": req.Email}, "user update failed")
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

	err = h.deps.Delete.Execute(c.Request.Context(), id)
	if err != nil {
		h.handleUseCaseError(c, err, map[string]interface{}{"user_id": id}, "user delete failed")
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

	resp, err := h.deps.Create.Execute(c.Request.Context(), req)
	if err != nil {
		h.handleUseCaseError(c, err, map[string]interface{}{"email": req.Email}, "register failed")
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

	resp, err := h.deps.Login.Execute(c.Request.Context(), req)
	if err != nil {
		h.handleUseCaseError(c, err, map[string]interface{}{"email": req.Email}, "login failed")
		return
	}

	h.logger.InfoWithFields("login successful", map[string]interface{}{"user_id": resp.User.ID, "email": req.Email})
	response.SuccessResponseOK(c, "Login successful", resp)
}
