package usecase

import (
	"context"
	"io"
	"time"

	"github.com/rulzi/hexa-go/internal/application/media/dto"
	domainmedia "github.com/rulzi/hexa-go/internal/domain/media"
)

// UpdateMediaUseCase handles updating a media
type UpdateMediaUseCase struct {
	mediaRepo    domainmedia.Repository
	mediaService *domainmedia.Service
	storage      domainmedia.Storage
	baseURL      string
}

// NewUpdateMediaUseCase creates a new UpdateMediaUseCase
func NewUpdateMediaUseCase(
	mediaRepo domainmedia.Repository,
	mediaService *domainmedia.Service,
	storage domainmedia.Storage,
	baseURL string,
) *UpdateMediaUseCase {
	return &UpdateMediaUseCase{
		mediaRepo:    mediaRepo,
		mediaService: mediaService,
		storage:      storage,
		baseURL:      baseURL,
	}
}

// Execute executes the update media use case
func (uc *UpdateMediaUseCase) Execute(ctx context.Context, id int64, filename string, file io.Reader) (*dto.MediaResponse, error) {
	// Get existing media
	existingMedia, err := uc.mediaRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingMedia == nil {
		return nil, domainmedia.ErrMediaNotFound
	}

	oldPath := existingMedia.Path

	// Save new file to storage
	storagePath, err := uc.storage.Save(ctx, filename, file)
	if err != nil {
		return nil, err
	}

	// Update fields
	existingMedia.Name = filename
	existingMedia.Path = storagePath
	existingMedia.UpdatedAt = time.Now()

	// Validate entity
	if err := existingMedia.Validate(); err != nil {
		// Clean up new file if validation fails
		_ = uc.storage.Delete(ctx, storagePath)
		return nil, err
	}

	// Update in repository
	updatedMedia, err := uc.mediaRepo.Update(ctx, existingMedia)
	if err != nil {
		// Clean up new file if database update fails
		_ = uc.storage.Delete(ctx, storagePath)
		return nil, err
	}

	// Delete old file from storage
	_ = uc.storage.Delete(ctx, oldPath)

	response := &dto.MediaResponse{
		ID:        updatedMedia.ID,
		Name:      updatedMedia.Name,
		Path:      updatedMedia.Path,
		URL:       dto.BuildURL(uc.baseURL, updatedMedia.Path),
		CreatedAt: updatedMedia.CreatedAt,
		UpdatedAt: updatedMedia.UpdatedAt,
	}

	return response, nil
}
