package media

import (
	"context"
	"io"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rulzi/hexa-go/internal/adapters/http/response"
	"github.com/rulzi/hexa-go/internal/application/media/dto"
	domainmedia "github.com/rulzi/hexa-go/internal/domain/media"
)

// UseCase interfaces for dependency injection and testing
type CreateMediaUseCase interface {
	Execute(ctx context.Context, filename string, file io.Reader) (*dto.MediaResponse, error)
}

type GetMediaUseCase interface {
	Execute(ctx context.Context, id int64) (*dto.MediaResponse, error)
}

type ListMediaUseCase interface {
	Execute(ctx context.Context, limit, offset int) (*dto.ListMediaResponse, error)
}

type UpdateMediaUseCase interface {
	Execute(ctx context.Context, id int64, filename string, file io.Reader) (*dto.MediaResponse, error)
}

type DeleteMediaUseCase interface {
	Execute(ctx context.Context, id int64) error
}

// Handler handles HTTP requests for media
type Handler struct {
	createUseCase CreateMediaUseCase
	getUseCase    GetMediaUseCase
	listUseCase   ListMediaUseCase
	updateUseCase UpdateMediaUseCase
	deleteUseCase DeleteMediaUseCase
}

// NewHandler creates a new Handler
func NewHandler(
	createUseCase CreateMediaUseCase,
	getUseCase GetMediaUseCase,
	listUseCase ListMediaUseCase,
	updateUseCase UpdateMediaUseCase,
	deleteUseCase DeleteMediaUseCase,
) *Handler {
	return &Handler{
		createUseCase: createUseCase,
		getUseCase:    getUseCase,
		listUseCase:   listUseCase,
		updateUseCase: updateUseCase,
		deleteUseCase: deleteUseCase,
	}
}

// Create handles POST /media (multipart/form-data with file field)
func (h *Handler) Create(c *gin.Context) {
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		response.ErrorResponseBadRequest(c, "file is required")
		return
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		response.ErrorResponseBadRequest(c, "failed to open file")
		return
	}
	defer func() {
		if err := src.Close(); err != nil {
			log.Printf("Failed to close file: %v", err)
		}
	}()

	// Execute use case
	resp, err := h.createUseCase.Execute(c.Request.Context(), file.Filename, src)
	if err != nil {
		if err == domainmedia.ErrNameRequired || err == domainmedia.ErrPathRequired {
			response.ErrorResponseBadRequest(c, err.Error())
		} else {
			response.ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	response.SuccessResponseCreated(c, "Media created successfully", resp)
}

// Get handles GET /media/:id
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorResponseBadRequest(c, "invalid media id")
		return
	}

	resp, err := h.getUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		if err == domainmedia.ErrMediaNotFound {
			response.ErrorResponseNotFound(c, err.Error())
		} else {
			response.ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	response.SuccessResponseOK(c, "Media retrieved successfully", resp)
}

// List handles GET /media
func (h *Handler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := h.listUseCase.Execute(c.Request.Context(), limit, offset)
	if err != nil {
		response.ErrorResponseInternalServerError(c, err.Error())
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

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		response.ErrorResponseBadRequest(c, "file is required")
		return
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		response.ErrorResponseBadRequest(c, "failed to open file")
		return
	}
	defer func() {
		if err := src.Close(); err != nil {
			log.Printf("Failed to close file: %v", err)
		}
	}()

	resp, err := h.updateUseCase.Execute(c.Request.Context(), id, file.Filename, src)
	if err != nil {
		switch err {
		case domainmedia.ErrMediaNotFound:
			response.ErrorResponseNotFound(c, err.Error())
		case domainmedia.ErrNameRequired, domainmedia.ErrPathRequired:
			response.ErrorResponseBadRequest(c, err.Error())
		default:
			response.ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	response.SuccessResponseOK(c, "Media updated successfully", resp)
}

// Delete handles DELETE /media/:id
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorResponseBadRequest(c, "invalid media id")
		return
	}

	err = h.deleteUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		if err == domainmedia.ErrMediaNotFound {
			response.ErrorResponseNotFound(c, err.Error())
		} else {
			response.ErrorResponseInternalServerError(c, err.Error())
		}
		return
	}

	response.SuccessResponseOK(c, "Media deleted successfully", nil)
}
