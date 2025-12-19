package usecase

import (
	"context"

	"github.com/rulzi/hexa-go/internal/application/media/dto"
	domainmedia "github.com/rulzi/hexa-go/internal/domain/media"
)

// GetMediaUseCase handles retrieving a media by ID
type GetMediaUseCase struct {
	mediaRepo domainmedia.Repository
	baseURL   string
}

// NewGetMediaUseCase creates a new GetMediaUseCase
func NewGetMediaUseCase(mediaRepo domainmedia.Repository, baseURL string) *GetMediaUseCase {
	return &GetMediaUseCase{
		mediaRepo: mediaRepo,
		baseURL:   baseURL,
	}
}

// Execute executes the get media use case
func (uc *GetMediaUseCase) Execute(ctx context.Context, id int64) (*dto.MediaResponse, error) {
	// Get from repository
	mediaEntity, err := uc.mediaRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if mediaEntity == nil {
		return nil, domainmedia.ErrMediaNotFound
	}

	response := &dto.MediaResponse{
		ID:        mediaEntity.ID,
		Name:      mediaEntity.Name,
		Path:      mediaEntity.Path,
		URL:       dto.BuildURL(uc.baseURL, mediaEntity.Path),
		CreatedAt: mediaEntity.CreatedAt,
		UpdatedAt: mediaEntity.UpdatedAt,
	}

	return response, nil
}
