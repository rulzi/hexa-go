package usecase

import (
	"context"
	"io"
	"time"

	"github.com/rulzi/hexa-go/internal/application/media/dto"
	mediaentity "github.com/rulzi/hexa-go/internal/domain/media/entity"
	mediaport "github.com/rulzi/hexa-go/internal/domain/media/port"
)

type createMediaDeps struct {
	mediaRepo mediaport.Repository
	storage   mediaport.Storage
	baseURL   string
}

// CreateMedia handles media creation
type CreateMedia struct {
	deps createMediaDeps
}

// NewCreateMedia creates a new CreateMedia use case.
func NewCreateMedia(mediaRepo mediaport.Repository, storage mediaport.Storage, baseURL string) *CreateMedia {
	return &CreateMedia{deps: createMediaDeps{
		mediaRepo: mediaRepo,
		storage:   storage,
		baseURL:   baseURL,
	}}
}

// Execute creates a new media
func (uc *CreateMedia) Execute(ctx context.Context, filename string, file io.Reader) (*dto.MediaResponse, error) {
	storagePath, err := uc.deps.storage.Save(ctx, filename, file)
	if err != nil {
		return nil, err
	}

	newMedia := &mediaentity.Media{
		Name:      filename,
		Path:      storagePath,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := newMedia.Validate(); err != nil {
		_ = uc.deps.storage.Delete(ctx, storagePath)
		return nil, err
	}

	createdMedia, err := uc.deps.mediaRepo.Create(ctx, newMedia)
	if err != nil {
		_ = uc.deps.storage.Delete(ctx, storagePath)
		return nil, err
	}

	return &dto.MediaResponse{
		ID:        createdMedia.ID,
		Name:      createdMedia.Name,
		Path:      createdMedia.Path,
		URL:       dto.BuildURL(uc.deps.baseURL, createdMedia.Path),
		CreatedAt: createdMedia.CreatedAt,
		UpdatedAt: createdMedia.UpdatedAt,
	}, nil
}
