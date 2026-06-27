package usecase

import (
	"context"
	"testing"
	"time"

	mediaentity "github.com/rulzi/hexa-go/internal/domain/media/entity"
	"github.com/stretchr/testify/assert"
)

func TestGetMedia_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewMediaUseCase(repo, storage, baseURL)

	mediaID := int64(1)
	mediaEntity := &mediaentity.Media{
		ID: mediaID, Name: "test.jpg", Path: "2025/12/19/test.jpg",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	repo.On("GetByID", ctx, mediaID).Return(mediaEntity, nil)

	result, err := uc.Get.Execute(ctx, mediaID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mediaEntity.ID, result.ID)
	assert.Contains(t, result.URL, baseURL)
	repo.AssertExpectations(t)
}

func TestGetMedia_Execute_MediaNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewMediaUseCase(repo, storage, baseURL)

	mediaID := int64(999)
	repo.On("GetByID", ctx, mediaID).Return(nil, nil)

	result, err := uc.Get.Execute(ctx, mediaID)

	assert.Error(t, err)
	assert.True(t, mediaentity.IsMediaNotFound(err))
	assert.Nil(t, result)
}
