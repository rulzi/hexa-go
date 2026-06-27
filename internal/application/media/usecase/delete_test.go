package usecase

import (
	"context"
	"testing"
	"time"

	mediaentity "github.com/rulzi/hexa-go/internal/domain/media/entity"
	mediaservice "github.com/rulzi/hexa-go/internal/domain/media/service"
	"github.com/stretchr/testify/assert"
)

func TestDeleteMedia_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := mediaservice.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewMediaUseCase(repo, service, storage, baseURL)

	mediaID := int64(1)
	existingMedia := &mediaentity.Media{
		ID: mediaID, Name: "test.jpg", Path: "path/test.jpg",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	repo.On("GetByID", ctx, mediaID).Return(existingMedia, nil)
	storage.On("Delete", ctx, existingMedia.Path).Return(nil)
	repo.On("Delete", ctx, mediaID).Return(nil)

	err := uc.Delete.Execute(ctx, mediaID)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	storage.AssertExpectations(t)
}

func TestDeleteMedia_Execute_MediaNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := mediaservice.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewMediaUseCase(repo, service, storage, baseURL)

	mediaID := int64(999)
	repo.On("GetByID", ctx, mediaID).Return(nil, nil)

	err := uc.Delete.Execute(ctx, mediaID)

	assert.Error(t, err)
	assert.True(t, mediaentity.IsMediaNotFound(err))
	storage.AssertNotCalled(t, "Delete")
	repo.AssertNotCalled(t, "Delete")
}
