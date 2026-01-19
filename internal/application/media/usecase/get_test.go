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
var _ *dto.MediaResponse

func TestNewGetMediaUseCase(t *testing.T) {
	repo := &mockMediaRepository{}
	baseURL := "http://localhost:8080"

	uc := NewGetMediaUseCase(repo, baseURL)

	assert.NotNil(t, uc)
	assert.Equal(t, repo, uc.mediaRepo)
	assert.Equal(t, baseURL, uc.baseURL)
}

func TestGetMediaUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	baseURL := "http://localhost:8080"

	uc := NewGetMediaUseCase(repo, baseURL)

	mediaID := int64(1)
	mediaEntity := &domainmedia.Media{
		ID:        mediaID,
		Name:      "test.jpg",
		Path:      "2025/12/19/test.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, mediaID).Return(mediaEntity, nil)

	result, err := uc.Execute(ctx, mediaID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mediaEntity.ID, result.ID)
	assert.Equal(t, mediaEntity.Name, result.Name)
	assert.Equal(t, mediaEntity.Path, result.Path)
	assert.Contains(t, result.URL, baseURL)
	assert.Contains(t, result.URL, mediaEntity.Path)
	assert.Equal(t, "http://localhost:8080/api/v1/media/files/2025/12/19/test.jpg", result.URL)

	repo.AssertExpectations(t)
}

func TestGetMediaUseCase_Execute_MediaNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	baseURL := "http://localhost:8080"

	uc := NewGetMediaUseCase(repo, baseURL)

	mediaID := int64(1)

	repo.On("GetByID", ctx, mediaID).Return(nil, nil)

	result, err := uc.Execute(ctx, mediaID)

	assert.Error(t, err)
	assert.Equal(t, domainmedia.ErrMediaNotFound, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
}

func TestGetMediaUseCase_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	baseURL := "http://localhost:8080"

	uc := NewGetMediaUseCase(repo, baseURL)

	mediaID := int64(1)
	repoError := errors.New("database error")

	repo.On("GetByID", ctx, mediaID).Return(nil, repoError)

	result, err := uc.Execute(ctx, mediaID)

	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
}

func TestGetMediaUseCase_Execute_URIBuilding(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	baseURL := "http://localhost:8080"

	uc := NewGetMediaUseCase(repo, baseURL)

	mediaID := int64(1)
	mediaEntity := &domainmedia.Media{
		ID:        mediaID,
		Name:      "test.jpg",
		Path:      "2025\\12\\19\\test.jpg", // Windows path separator
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, mediaID).Return(mediaEntity, nil)

	result, err := uc.Execute(ctx, mediaID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	// Should normalize path separators
	assert.NotContains(t, result.URL, "\\")
	assert.Contains(t, result.URL, "/2025/12/19/test.jpg")

	repo.AssertExpectations(t)
}

func TestGetMediaUseCase_Execute_PathWithoutLeadingSlash(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	baseURL := "http://localhost:8080"

	uc := NewGetMediaUseCase(repo, baseURL)

	mediaID := int64(1)
	mediaEntity := &domainmedia.Media{
		ID:        mediaID,
		Name:      "test.jpg",
		Path:      "2025/12/19/test.jpg", // No leading slash
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, mediaID).Return(mediaEntity, nil)

	result, err := uc.Execute(ctx, mediaID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	// Should add leading slash
	assert.Contains(t, result.URL, "/2025/12/19/test.jpg")
	assert.Equal(t, "http://localhost:8080/api/v1/media/files/2025/12/19/test.jpg", result.URL)

	repo.AssertExpectations(t)
}

