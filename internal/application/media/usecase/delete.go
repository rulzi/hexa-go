package usecase

import (
	"context"

	mediaentity "github.com/rulzi/hexa-go/internal/domain/media/entity"
	mediaport "github.com/rulzi/hexa-go/internal/domain/media/port"
)

type deleteMediaDeps struct {
	mediaRepo mediaport.Repository
	storage   mediaport.Storage
}

// DeleteMedia handles media deletion
type DeleteMedia struct {
	deps deleteMediaDeps
}

// NewDeleteMedia creates a new DeleteMedia use case.
func NewDeleteMedia(mediaRepo mediaport.Repository, storage mediaport.Storage) *DeleteMedia {
	return &DeleteMedia{deps: deleteMediaDeps{
		mediaRepo: mediaRepo,
		storage:   storage,
	}}
}

// Execute deletes a media
func (uc *DeleteMedia) Execute(ctx context.Context, id int64) error {
	existingMedia, err := uc.deps.mediaRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingMedia == nil {
		return mediaentity.NewMediaNotFound()
	}

	if err := uc.deps.storage.Delete(ctx, existingMedia.Path); err != nil {
		return err
	}

	if err := uc.deps.mediaRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
