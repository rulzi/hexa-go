package usecase

import (
	"context"

	"github.com/rulzi/hexa-go/internal/application/media/dto"
)

// ListMedia handles media listing with pagination
type ListMedia struct {
	deps mediaDeps
}

func newListMedia(deps mediaDeps) *ListMedia {
	return &ListMedia{deps: deps}
}

// Execute lists media with pagination
func (uc *ListMedia) Execute(ctx context.Context, limit, offset int) (*dto.ListMediaResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	mediaList, err := uc.deps.mediaRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.deps.mediaRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	mediaResponses := make([]dto.MediaResponse, len(mediaList))
	for i, m := range mediaList {
		mediaResponses[i] = dto.MediaResponse{
			ID:        m.ID,
			Name:      m.Name,
			Path:      m.Path,
			URL:       dto.BuildURL(uc.deps.baseURL, m.Path),
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
