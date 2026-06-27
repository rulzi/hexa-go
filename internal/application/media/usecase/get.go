package usecase

import (
	"context"

	"github.com/rulzi/hexa-go/internal/application/media/dto"
	mediaentity "github.com/rulzi/hexa-go/internal/domain/media/entity"
	mediaport "github.com/rulzi/hexa-go/internal/domain/media/port"
)

type getMediaDeps struct {
	mediaRepo mediaport.Repository
	baseURL   string
}

// GetMedia handles media retrieval by ID
type GetMedia struct {
	deps getMediaDeps
}

// NewGetMedia creates a new GetMedia use case.
func NewGetMedia(mediaRepo mediaport.Repository, baseURL string) *GetMedia {
	return &GetMedia{deps: getMediaDeps{
		mediaRepo: mediaRepo,
		baseURL:   baseURL,
	}}
}

// Execute retrieves a media by ID
func (uc *GetMedia) Execute(ctx context.Context, id int64) (*dto.MediaResponse, error) {
	mediaEntity, err := uc.deps.mediaRepo.GetByID(ctx, id)
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
		URL:       dto.BuildURL(uc.deps.baseURL, mediaEntity.Path),
		CreatedAt: mediaEntity.CreatedAt,
		UpdatedAt: mediaEntity.UpdatedAt,
	}, nil
}
