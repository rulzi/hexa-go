package usecase

import (
	mediaport "github.com/rulzi/hexa-go/internal/domain/media/port"
)

type mediaDeps struct {
	mediaRepo mediaport.Repository
	storage   mediaport.Storage
	baseURL   string
}

// MediaUseCase groups all media use case operations
type MediaUseCase struct {
	Create *CreateMedia
	Get    *GetMedia
	List   *ListMedia
	Update *UpdateMedia
	Delete *DeleteMedia
}

// NewMediaUseCase creates a new MediaUseCase
func NewMediaUseCase(
	mediaRepo mediaport.Repository,
	storage mediaport.Storage,
	baseURL string,
) *MediaUseCase {
	deps := mediaDeps{
		mediaRepo: mediaRepo,
		storage:   storage,
		baseURL:   baseURL,
	}

	return &MediaUseCase{
		Create: newCreateMedia(deps),
		Get:    newGetMedia(deps),
		List:   newListMedia(deps),
		Update: newUpdateMedia(deps),
		Delete: newDeleteMedia(deps),
	}
}
