package usecase

import (
	"context"

	"github.com/rulzi/hexa-go/internal/application/media/dto"
	domainmedia "github.com/rulzi/hexa-go/internal/domain/media"
)

// ListMediaUseCase handles listing media with pagination
type ListMediaUseCase struct {
	mediaRepo domainmedia.Repository
	baseURL   string
}

// NewListMediaUseCase creates a new ListMediaUseCase
func NewListMediaUseCase(mediaRepo domainmedia.Repository, baseURL string) *ListMediaUseCase {
	return &ListMediaUseCase{
		mediaRepo: mediaRepo,
		baseURL:   baseURL,
	}
}

// Execute executes the list media use case
func (uc *ListMediaUseCase) Execute(ctx context.Context, limit, offset int) (*dto.ListMediaResponse, error) {
	// Default pagination
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Get media from repository
	mediaList, err := uc.mediaRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := uc.mediaRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
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

	response := &dto.ListMediaResponse{
		Media:  mediaResponses,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	return response, nil
}
