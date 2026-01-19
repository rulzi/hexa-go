package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	domainmedia "github.com/rulzi/hexa-go/internal/domain/media"
	"github.com/stretchr/testify/assert"
)

func TestNewDeleteMediaUseCase(t *testing.T) {
	repo := &mockMediaRepository{}
	storage := &mockMediaStorage{}

	uc := NewDeleteMediaUseCase(repo, storage)

	assert.NotNil(t, uc)
	assert.Equal(t, repo, uc.mediaRepo)
	assert.Equal(t, storage, uc.storage)
}

func TestDeleteMediaUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	storage := &mockMediaStorage{}

	uc := NewDeleteMediaUseCase(repo, storage)

	mediaID := int64(1)
	existingMedia := &domainmedia.Media{
		ID:        mediaID,
		Name:      "test.jpg",
		Path:      "2025/12/19/test.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, mediaID).Return(existingMedia, nil)
	storage.On("Delete", ctx, existingMedia.Path).Return(nil)
	repo.On("Delete", ctx, mediaID).Return(nil)

	err := uc.Execute(ctx, mediaID)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	storage.AssertExpectations(t)
}

func TestDeleteMediaUseCase_Execute_MediaNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	storage := &mockMediaStorage{}

	uc := NewDeleteMediaUseCase(repo, storage)

	mediaID := int64(1)

	repo.On("GetByID", ctx, mediaID).Return(nil, nil)

	err := uc.Execute(ctx, mediaID)

	assert.Error(t, err)
	assert.Equal(t, domainmedia.ErrMediaNotFound, err)

	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Delete")
	storage.AssertNotCalled(t, "Delete")
}

func TestDeleteMediaUseCase_Execute_GetByIDError(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	storage := &mockMediaStorage{}

	uc := NewDeleteMediaUseCase(repo, storage)

	mediaID := int64(1)
	repoError := errors.New("database error")

	repo.On("GetByID", ctx, mediaID).Return(nil, repoError)

	err := uc.Execute(ctx, mediaID)

	assert.Error(t, err)
	assert.Equal(t, repoError, err)

	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Delete")
	storage.AssertNotCalled(t, "Delete")
}

func TestDeleteMediaUseCase_Execute_StorageDeleteError(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	storage := &mockMediaStorage{}

	uc := NewDeleteMediaUseCase(repo, storage)

	mediaID := int64(1)
	existingMedia := &domainmedia.Media{
		ID:        mediaID,
		Name:      "test.jpg",
		Path:      "2025/12/19/test.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	storageError := errors.New("storage delete error")

	repo.On("GetByID", ctx, mediaID).Return(existingMedia, nil)
	storage.On("Delete", ctx, existingMedia.Path).Return(storageError)

	err := uc.Execute(ctx, mediaID)

	assert.Error(t, err)
	assert.Equal(t, storageError, err)

	repo.AssertExpectations(t)
	storage.AssertExpectations(t)
	repo.AssertNotCalled(t, "Delete")
}

func TestDeleteMediaUseCase_Execute_RepositoryDeleteError(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	storage := &mockMediaStorage{}

	uc := NewDeleteMediaUseCase(repo, storage)

	mediaID := int64(1)
	existingMedia := &domainmedia.Media{
		ID:        mediaID,
		Name:      "test.jpg",
		Path:      "2025/12/19/test.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repoError := errors.New("database delete error")

	repo.On("GetByID", ctx, mediaID).Return(existingMedia, nil)
	storage.On("Delete", ctx, existingMedia.Path).Return(nil)
	repo.On("Delete", ctx, mediaID).Return(repoError)

	err := uc.Execute(ctx, mediaID)

	assert.Error(t, err)
	assert.Equal(t, repoError, err)

	repo.AssertExpectations(t)
	storage.AssertExpectations(t)
}

