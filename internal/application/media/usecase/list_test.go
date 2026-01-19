package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rulzi/hexa-go/internal/application/media/dto"
	domainmedia "github.com/rulzi/hexa-go/internal/domain/media"
	"github.com/stretchr/testify/assert"
)

// Ensure dto is used
var _ *dto.ListMediaResponse

func TestNewListMediaUseCase(t *testing.T) {
	repo := &mockMediaRepository{}
	baseURL := "http://localhost:8080"

	uc := NewListMediaUseCase(repo, baseURL)

	assert.NotNil(t, uc)
	assert.Equal(t, repo, uc.mediaRepo)
	assert.Equal(t, baseURL, uc.baseURL)
}

func TestListMediaUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	baseURL := "http://localhost:8080"

	uc := NewListMediaUseCase(repo, baseURL)

	limit := 10
	offset := 0
	mediaList := []*domainmedia.Media{
		{
			ID:        1,
			Name:      "test1.jpg",
			Path:      "2025/12/19/test1.jpg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Name:      "test2.jpg",
			Path:      "2025/12/19/test2.jpg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	total := int64(2)

	repo.On("List", ctx, limit, offset).Return(mediaList, nil)
	repo.On("Count", ctx).Return(total, nil)

	result, err := uc.Execute(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, total, result.Total)
	assert.Equal(t, limit, result.Limit)
	assert.Equal(t, offset, result.Offset)
	assert.Equal(t, len(mediaList), len(result.Media))
	assert.Equal(t, mediaList[0].ID, result.Media[0].ID)
	assert.Equal(t, mediaList[1].ID, result.Media[1].ID)
	assert.Contains(t, result.Media[0].URL, baseURL)
	assert.Contains(t, result.Media[0].URL, mediaList[0].Path)

	repo.AssertExpectations(t)
}

func TestListMediaUseCase_Execute_DefaultPagination(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	baseURL := "http://localhost:8080"

	uc := NewListMediaUseCase(repo, baseURL)

	mediaList := []*domainmedia.Media{}
	total := int64(0)

	repo.On("List", ctx, 10, 0).Return(mediaList, nil)
	repo.On("Count", ctx).Return(total, nil)

	result, err := uc.Execute(ctx, -1, -1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 0, result.Offset)

	repo.AssertExpectations(t)
}

func TestListMediaUseCase_Execute_ListError(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	baseURL := "http://localhost:8080"

	uc := NewListMediaUseCase(repo, baseURL)

	limit := 10
	offset := 0
	listError := errors.New("list error")

	repo.On("List", ctx, limit, offset).Return(nil, listError)

	result, err := uc.Execute(ctx, limit, offset)

	assert.Error(t, err)
	assert.Equal(t, listError, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Count")
}

func TestListMediaUseCase_Execute_CountError(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	baseURL := "http://localhost:8080"

	uc := NewListMediaUseCase(repo, baseURL)

	limit := 10
	offset := 0
	mediaList := []*domainmedia.Media{
		{
			ID:        1,
			Name:      "test1.jpg",
			Path:      "2025/12/19/test1.jpg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	countError := errors.New("count error")

	repo.On("List", ctx, limit, offset).Return(mediaList, nil)
	repo.On("Count", ctx).Return(int64(0), countError)

	result, err := uc.Execute(ctx, limit, offset)

	assert.Error(t, err)
	assert.Equal(t, countError, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
}

func TestListMediaUseCase_Execute_EmptyList(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	baseURL := "http://localhost:8080"

	uc := NewListMediaUseCase(repo, baseURL)

	limit := 10
	offset := 0
	mediaList := []*domainmedia.Media{}
	total := int64(0)

	repo.On("List", ctx, limit, offset).Return(mediaList, nil)
	repo.On("Count", ctx).Return(total, nil)

	result, err := uc.Execute(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(0), result.Total)
	assert.Equal(t, 0, len(result.Media))

	repo.AssertExpectations(t)
}

func TestListMediaUseCase_Execute_URIBuilding(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	baseURL := "http://localhost:8080"

	uc := NewListMediaUseCase(repo, baseURL)

	limit := 10
	offset := 0
	mediaList := []*domainmedia.Media{
		{
			ID:        1,
			Name:      "test.jpg",
			Path:      "2025/12/19/test.jpg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	total := int64(1)

	repo.On("List", ctx, limit, offset).Return(mediaList, nil)
	repo.On("Count", ctx).Return(total, nil)

	result, err := uc.Execute(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "http://localhost:8080/api/v1/media/files/2025/12/19/test.jpg", result.Media[0].URL)

	repo.AssertExpectations(t)
}

