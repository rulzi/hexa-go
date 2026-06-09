package media

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rulzi/hexa-go/internal/adapters/http/errmapper"
	"github.com/rulzi/hexa-go/internal/adapters/http/response"
	"github.com/rulzi/hexa-go/internal/application/media/dto"
	domainerrs "github.com/rulzi/hexa-go/internal/domain/errs"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
)

// UseCase is the interface for media operations
type UseCase interface {
	Create(ctx context.Context, filename string, file io.Reader) (*dto.MediaResponse, error)
	Get(ctx context.Context, id int64) (*dto.MediaResponse, error)
	List(ctx context.Context, limit, offset int) (*dto.ListMediaResponse, error)
	Update(ctx context.Context, id int64, filename string, file io.Reader) (*dto.MediaResponse, error)
	Delete(ctx context.Context, id int64) error
}

// Handler handles HTTP requests for media
type Handler struct {
	mediaUseCase UseCase
	logger       logger.Logger
}

// NewHandler creates a new Handler
func NewHandler(mediaUseCase UseCase, appLogger logger.Logger) *Handler {
	return &Handler{
		mediaUseCase: mediaUseCase,
		logger:       appLogger,
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

// Create handles POST /media (multipart/form-data with file field)
func (h *Handler) Create(c *gin.Context) {
	filename, reader, err := h.parseUploadFile(c)
	if err != nil {
		response.ErrorResponseBadRequest(c, err.Error())
		return
	}

	resp, err := h.mediaUseCase.Create(c.Request.Context(), filename, reader)
	if err != nil {
		h.handleUseCaseError(c, err, map[string]interface{}{"filename": filename}, "media create failed")
		return
	}

	h.logger.InfoWithFields("media created", map[string]interface{}{"media_id": resp.ID, "filename": filename})
	response.SuccessResponseCreated(c, "Media created successfully", resp)
}

// Get handles GET /media/:id
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorResponseBadRequest(c, "invalid media id")
		return
	}

	resp, err := h.mediaUseCase.Get(c.Request.Context(), id)
	if err != nil {
		h.handleUseCaseError(c, err, map[string]interface{}{"media_id": id}, "media get failed")
		return
	}

	response.SuccessResponseOK(c, "Media retrieved successfully", resp)
}

// List handles GET /media
func (h *Handler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := h.mediaUseCase.List(c.Request.Context(), limit, offset)
	if err != nil {
		h.handleUseCaseError(c, err, nil, "media list failed")
		return
	}

	response.SuccessResponseOK(c, "Media retrieved successfully", resp)
}

// Update handles PUT /media/:id (multipart/form-data with file field)
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorResponseBadRequest(c, "invalid media id")
		return
	}

	filename, reader, err := h.parseUploadFile(c)
	if err != nil {
		response.ErrorResponseBadRequest(c, err.Error())
		return
	}

	resp, err := h.mediaUseCase.Update(c.Request.Context(), id, filename, reader)
	if err != nil {
		h.handleUseCaseError(c, err, map[string]interface{}{"media_id": id, "filename": filename}, "media update failed")
		return
	}

	h.logger.InfoWithFields("media updated", map[string]interface{}{"media_id": id, "filename": filename})
	response.SuccessResponseOK(c, "Media updated successfully", resp)
}

func (h *Handler) parseUploadFile(c *gin.Context) (string, io.Reader, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", nil, fmt.Errorf("file is required")
	}

	src, err := file.Open()
	if err != nil {
		return "", nil, fmt.Errorf("failed to open file")
	}
	defer func() {
		if closeErr := src.Close(); closeErr != nil {
			h.logger.WarnWithFields("failed to close uploaded file", map[string]interface{}{"error": closeErr.Error()})
		}
	}()

	return validateUpload(file, src)
}

// Delete handles DELETE /media/:id
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorResponseBadRequest(c, "invalid media id")
		return
	}

	err = h.mediaUseCase.Delete(c.Request.Context(), id)
	if err != nil {
		h.handleUseCaseError(c, err, map[string]interface{}{"media_id": id}, "media delete failed")
		return
	}

	h.logger.InfoWithFields("media deleted", map[string]interface{}{"media_id": id})
	response.SuccessResponseOK(c, "Media deleted successfully", nil)
}
