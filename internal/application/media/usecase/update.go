package usecase

import (
	"context"
	"io"
	"time"

	"github.com/rulzi/hexa-go/internal/application/media/dto"
	mediaentity "github.com/rulzi/hexa-go/internal/domain/media/entity"
	mediaport "github.com/rulzi/hexa-go/internal/domain/media/port"
)

type updateMediaDeps struct {
	mediaRepo mediaport.Repository
	storage   mediaport.Storage
	baseURL   string
}

// UpdateMedia handles media updates
type UpdateMedia struct {
	deps updateMediaDeps
}

// NewUpdateMedia creates a new UpdateMedia use case.
func NewUpdateMedia(mediaRepo mediaport.Repository, storage mediaport.Storage, baseURL string) *UpdateMedia {
	return &UpdateMedia{deps: updateMediaDeps{
		mediaRepo: mediaRepo,
		storage:   storage,
		baseURL:   baseURL,
	}}
}

// Execute updates a media
func (uc *UpdateMedia) Execute(ctx context.Context, id int64, filename string, file io.Reader) (*dto.MediaResponse, error) {
	existingMedia, err := uc.deps.mediaRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingMedia == nil {
		return nil, mediaentity.NewMediaNotFound()
	}

	oldPath := existingMedia.Path

	storagePath, err := uc.deps.storage.Save(ctx, filename, file)
	if err != nil {
		return nil, err
	}

	existingMedia.Name = filename
	existingMedia.Path = storagePath
	existingMedia.UpdatedAt = time.Now()

	if err := existingMedia.Validate(); err != nil {
		_ = uc.deps.storage.Delete(ctx, storagePath)
		return nil, err
	}

	updatedMedia, err := uc.deps.mediaRepo.Update(ctx, existingMedia)
	if err != nil {
		_ = uc.deps.storage.Delete(ctx, storagePath)
		return nil, err
	}

	_ = uc.deps.storage.Delete(ctx, oldPath)

	return &dto.MediaResponse{
		ID:        updatedMedia.ID,
		Name:      updatedMedia.Name,
		Path:      updatedMedia.Path,
		URL:       dto.BuildURL(uc.deps.baseURL, updatedMedia.Path),
		CreatedAt: updatedMedia.CreatedAt,
		UpdatedAt: updatedMedia.UpdatedAt,
	}, nil
}
