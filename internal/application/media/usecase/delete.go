package usecase

import (
	"context"

	domainmedia "github.com/rulzi/hexa-go/internal/domain/media"
)

// DeleteMediaUseCase handles deleting a media
type DeleteMediaUseCase struct {
	mediaRepo domainmedia.Repository
	storage   domainmedia.Storage
}

// NewDeleteMediaUseCase creates a new DeleteMediaUseCase
func NewDeleteMediaUseCase(mediaRepo domainmedia.Repository, storage domainmedia.Storage) *DeleteMediaUseCase {
	return &DeleteMediaUseCase{
		mediaRepo: mediaRepo,
		storage:   storage,
	}
}

// Execute executes the delete media use case
func (uc *DeleteMediaUseCase) Execute(ctx context.Context, id int64) error {
	// Check if media exists
	existingMedia, err := uc.mediaRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingMedia == nil {
		return domainmedia.ErrMediaNotFound
	}

	// Delete file from storage
	if err := uc.storage.Delete(ctx, existingMedia.Path); err != nil {
		return err
	}

	// Delete media from database
	if err := uc.mediaRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
