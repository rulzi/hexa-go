package usecase

import (
	"context"
	"io"
	"time"

	"github.com/rulzi/hexa-go/internal/application/media/dto"
	domainmedia "github.com/rulzi/hexa-go/internal/domain/media"
)

// CreateMediaUseCase handles the creation of a new media
type CreateMediaUseCase struct {
	mediaRepo    domainmedia.Repository
	mediaService *domainmedia.Service
	storage      domainmedia.Storage
	baseURL      string
}

// NewCreateMediaUseCase creates a new CreateMediaUseCase
func NewCreateMediaUseCase(
	mediaRepo domainmedia.Repository,
	mediaService *domainmedia.Service,
	storage domainmedia.Storage,
	baseURL string,
) *CreateMediaUseCase {
	return &CreateMediaUseCase{
		mediaRepo:    mediaRepo,
		mediaService: mediaService,
		storage:      storage,
		baseURL:      baseURL,
	}
}

// Execute executes the create media use case
func (uc *CreateMediaUseCase) Execute(ctx context.Context, filename string, file io.Reader) (*dto.MediaResponse, error) {
	// Save file to storage
	storagePath, err := uc.storage.Save(ctx, filename, file)
	if err != nil {
		return nil, err
	}

	// Create media entity
	newMedia := &domainmedia.Media{
		Name:      filename,
		Path:      storagePath,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Validate entity
	if err := newMedia.Validate(); err != nil {
		// Clean up file if validation fails
		_ = uc.storage.Delete(ctx, storagePath)
		return nil, err
	}

	// Save to repository
	createdMedia, err := uc.mediaRepo.Create(ctx, newMedia)
	if err != nil {
		// Clean up file if database save fails
		_ = uc.storage.Delete(ctx, storagePath)
		return nil, err
	}

	// Return response DTO
	return &dto.MediaResponse{
		ID:        createdMedia.ID,
		Name:      createdMedia.Name,
		Path:      createdMedia.Path,
		URL:       dto.BuildURL(uc.baseURL, createdMedia.Path),
		CreatedAt: createdMedia.CreatedAt,
		UpdatedAt: createdMedia.UpdatedAt,
	}, nil
}
