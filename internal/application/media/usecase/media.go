package usecase

import (
	"context"
	"io"
	"time"

	"github.com/rulzi/hexa-go/internal/application/media/dto"
	mediaentity "github.com/rulzi/hexa-go/internal/domain/media/entity"
	mediaport "github.com/rulzi/hexa-go/internal/domain/media/port"
	mediaservice "github.com/rulzi/hexa-go/internal/domain/media/service"
)

// MediaUseCase handles all media operations (create, get, list, update, delete)
type MediaUseCase struct {
	mediaRepo    mediaport.Repository
	mediaService *mediaservice.Service
	storage      mediaport.Storage
	baseURL      string
}

// NewMediaUseCase creates a new MediaUseCase
func NewMediaUseCase(
	mediaRepo mediaport.Repository,
	mediaService *mediaservice.Service,
	storage mediaport.Storage,
	baseURL string,
) *MediaUseCase {
	return &MediaUseCase{
		mediaRepo:    mediaRepo,
		mediaService: mediaService,
		storage:      storage,
		baseURL:      baseURL,
	}
}

// Create creates a new media
func (uc *MediaUseCase) Create(ctx context.Context, filename string, file io.Reader) (*dto.MediaResponse, error) {
	storagePath, err := uc.storage.Save(ctx, filename, file)
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
		_ = uc.storage.Delete(ctx, storagePath)
		return nil, err
	}

	createdMedia, err := uc.mediaRepo.Create(ctx, newMedia)
	if err != nil {
		_ = uc.storage.Delete(ctx, storagePath)
		return nil, err
	}

	return &dto.MediaResponse{
		ID:        createdMedia.ID,
		Name:      createdMedia.Name,
		Path:      createdMedia.Path,
		URL:       dto.BuildURL(uc.baseURL, createdMedia.Path),
		CreatedAt: createdMedia.CreatedAt,
		UpdatedAt: createdMedia.UpdatedAt,
	}, nil
}

// Get retrieves a media by ID
func (uc *MediaUseCase) Get(ctx context.Context, id int64) (*dto.MediaResponse, error) {
	mediaEntity, err := uc.mediaRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if mediaEntity == nil {
		return nil, mediaentity.NewMediaNotFound()
	}

	return &dto.MediaResponse{
		ID:        mediaEntity.ID,
		Name:      mediaEntity.Name,
		Path:      mediaEntity.Path,
		URL:       dto.BuildURL(uc.baseURL, mediaEntity.Path),
		CreatedAt: mediaEntity.CreatedAt,
		UpdatedAt: mediaEntity.UpdatedAt,
	}, nil
}

// List lists media with pagination
func (uc *MediaUseCase) List(ctx context.Context, limit, offset int) (*dto.ListMediaResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	mediaList, err := uc.mediaRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.mediaRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	mediaResponses := make([]dto.MediaResponse, len(mediaList))
	for i, m := range mediaList {
		mediaResponses[i] = dto.MediaResponse{
			ID:        m.ID,
			Name:      m.Name,
			Path:      m.Path,
			URL:       dto.BuildURL(uc.baseURL, m.Path),
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}

	return &dto.ListMediaResponse{
		Media:  mediaResponses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

// Update updates a media
func (uc *MediaUseCase) Update(ctx context.Context, id int64, filename string, file io.Reader) (*dto.MediaResponse, error) {
	existingMedia, err := uc.mediaRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingMedia == nil {
		return nil, mediaentity.NewMediaNotFound()
	}

	oldPath := existingMedia.Path

	storagePath, err := uc.storage.Save(ctx, filename, file)
	if err != nil {
		return nil, err
	}

	existingMedia.Name = filename
	existingMedia.Path = storagePath
	existingMedia.UpdatedAt = time.Now()

	if err := existingMedia.Validate(); err != nil {
		_ = uc.storage.Delete(ctx, storagePath)
		return nil, err
	}

	updatedMedia, err := uc.mediaRepo.Update(ctx, existingMedia)
	if err != nil {
		_ = uc.storage.Delete(ctx, storagePath)
		return nil, err
	}

	_ = uc.storage.Delete(ctx, oldPath)

	return &dto.MediaResponse{
		ID:        updatedMedia.ID,
		Name:      updatedMedia.Name,
		Path:      updatedMedia.Path,
		URL:       dto.BuildURL(uc.baseURL, updatedMedia.Path),
		CreatedAt: updatedMedia.CreatedAt,
		UpdatedAt: updatedMedia.UpdatedAt,
	}, nil
}

// Delete deletes a media
func (uc *MediaUseCase) Delete(ctx context.Context, id int64) error {
	existingMedia, err := uc.mediaRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingMedia == nil {
		return mediaentity.NewMediaNotFound()
	}

	if err := uc.storage.Delete(ctx, existingMedia.Path); err != nil {
		return err
	}

	if err := uc.mediaRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
