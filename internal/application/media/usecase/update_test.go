package usecase

import (
	"context"
	"strings"
	"testing"
	"time"

	mediaentity "github.com/rulzi/hexa-go/internal/domain/media/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateMedia_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewUpdateMedia(repo, storage, baseURL)

	mediaID := int64(1)
	existingMedia := &mediaentity.Media{
		ID: mediaID, Name: "old.jpg", Path: "old/path.jpg",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	newPath := "new/path.jpg"
	updatedMedia := &mediaentity.Media{
		ID: mediaID, Name: "new.jpg", Path: newPath,
		CreatedAt: existingMedia.CreatedAt, UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, mediaID).Return(existingMedia, nil)
	storage.On("Save", ctx, "new.jpg", mock.Anything).Return(newPath, nil)
	repo.On("Update", ctx, mock.Anything).Return(updatedMedia, nil)
	storage.On("Delete", ctx, existingMedia.Path).Return(nil)

	result, err := uc.Execute(ctx, mediaID, "new.jpg", strings.NewReader("new content"))

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "new.jpg", result.Name)
	repo.AssertExpectations(t)
	storage.AssertExpectations(t)
}

func TestUpdateMedia_Execute_MediaNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewUpdateMedia(repo, storage, baseURL)

	mediaID := int64(999)
	repo.On("GetByID", ctx, mediaID).Return(nil, nil)

	result, err := uc.Execute(ctx, mediaID, "new.jpg", strings.NewReader("x"))

	assert.Error(t, err)
	assert.True(t, mediaentity.IsMediaNotFound(err))
	assert.Nil(t, result)
	storage.AssertNotCalled(t, "Save")
}
