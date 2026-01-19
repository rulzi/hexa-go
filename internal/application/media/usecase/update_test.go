package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/rulzi/hexa-go/internal/application/media/dto"
	domainmedia "github.com/rulzi/hexa-go/internal/domain/media"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Ensure dto is used
var _ *dto.MediaResponse

func TestNewUpdateMediaUseCase(t *testing.T) {
	repo := &mockMediaRepository{}
	service := domainmedia.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewUpdateMediaUseCase(repo, service, storage, baseURL)

	assert.NotNil(t, uc)
	assert.Equal(t, repo, uc.mediaRepo)
	assert.Equal(t, service, uc.mediaService)
	assert.Equal(t, storage, uc.storage)
	assert.Equal(t, baseURL, uc.baseURL)
}

func TestUpdateMediaUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := domainmedia.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewUpdateMediaUseCase(repo, service, storage, baseURL)

	mediaID := int64(1)
	existingMedia := &domainmedia.Media{
		ID:        mediaID,
		Name:      "old.jpg",
		Path:      "2025/12/19/old.jpg",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
	}

	newFilename := "new.jpg"
	newFile := strings.NewReader("new file content")
	newStoragePath := "2025/12/19/new.jpg"

	updatedMedia := &domainmedia.Media{
		ID:        mediaID,
		Name:      newFilename,
		Path:      newStoragePath,
		CreatedAt: existingMedia.CreatedAt,
		UpdatedAt: time.Now(),
	}

	repo.On("GetByID", ctx, mediaID).Return(existingMedia, nil)
	storage.On("Save", ctx, newFilename, newFile).Return(newStoragePath, nil)
	repo.On("Update", ctx, mock.AnythingOfType("*media.Media")).Return(updatedMedia, nil)
	storage.On("Delete", ctx, existingMedia.Path).Return(nil)

	result, err := uc.Execute(ctx, mediaID, newFilename, newFile)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updatedMedia.ID, result.ID)
	assert.Equal(t, newFilename, result.Name)
	assert.Equal(t, newStoragePath, result.Path)
	assert.Contains(t, result.URL, baseURL)
	assert.Contains(t, result.URL, newStoragePath)

	repo.AssertExpectations(t)
	storage.AssertExpectations(t)
}

func TestUpdateMediaUseCase_Execute_MediaNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := domainmedia.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewUpdateMediaUseCase(repo, service, storage, baseURL)

	mediaID := int64(1)
	newFilename := "new.jpg"
	newFile := strings.NewReader("new file content")

	repo.On("GetByID", ctx, mediaID).Return(nil, nil)

	result, err := uc.Execute(ctx, mediaID, newFilename, newFile)

	assert.Error(t, err)
	assert.Equal(t, domainmedia.ErrMediaNotFound, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Update")
	storage.AssertNotCalled(t, "Save")
}

func TestUpdateMediaUseCase_Execute_GetByIDError(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := domainmedia.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewUpdateMediaUseCase(repo, service, storage, baseURL)

	mediaID := int64(1)
	newFilename := "new.jpg"
	newFile := strings.NewReader("new file content")
	repoError := errors.New("database error")

	repo.On("GetByID", ctx, mediaID).Return(nil, repoError)

	result, err := uc.Execute(ctx, mediaID, newFilename, newFile)

	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
	storage.AssertNotCalled(t, "Save")
}

func TestUpdateMediaUseCase_Execute_StorageSaveError(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := domainmedia.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewUpdateMediaUseCase(repo, service, storage, baseURL)

	mediaID := int64(1)
	existingMedia := &domainmedia.Media{
		ID:        mediaID,
		Name:      "old.jpg",
		Path:      "2025/12/19/old.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newFilename := "new.jpg"
	newFile := strings.NewReader("new file content")
	storageError := errors.New("storage save error")

	repo.On("GetByID", ctx, mediaID).Return(existingMedia, nil)
	storage.On("Save", ctx, newFilename, newFile).Return("", storageError)

	result, err := uc.Execute(ctx, mediaID, newFilename, newFile)

	assert.Error(t, err)
	assert.Equal(t, storageError, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
	storage.AssertExpectations(t)
	repo.AssertNotCalled(t, "Update")
	storage.AssertNotCalled(t, "Delete")
}

func TestUpdateMediaUseCase_Execute_ValidationError(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := domainmedia.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewUpdateMediaUseCase(repo, service, storage, baseURL)

	mediaID := int64(1)
	existingMedia := &domainmedia.Media{
		ID:        mediaID,
		Name:      "old.jpg",
		Path:      "2025/12/19/old.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Empty filename will cause validation error
	newFilename := ""
	newFile := strings.NewReader("new file content")
	newStoragePath := "2025/12/19/.jpg"

	repo.On("GetByID", ctx, mediaID).Return(existingMedia, nil)
	storage.On("Save", ctx, newFilename, newFile).Return(newStoragePath, nil)
	storage.On("Delete", ctx, newStoragePath).Return(nil)

	result, err := uc.Execute(ctx, mediaID, newFilename, newFile)

	assert.Error(t, err)
	assert.Equal(t, domainmedia.ErrNameRequired, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
	storage.AssertExpectations(t)
	repo.AssertNotCalled(t, "Update")
	storage.AssertNotCalled(t, "Delete", existingMedia.Path)
}

func TestUpdateMediaUseCase_Execute_RepositoryUpdateError(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := domainmedia.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewUpdateMediaUseCase(repo, service, storage, baseURL)

	mediaID := int64(1)
	existingMedia := &domainmedia.Media{
		ID:        mediaID,
		Name:      "old.jpg",
		Path:      "2025/12/19/old.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newFilename := "new.jpg"
	newFile := strings.NewReader("new file content")
	newStoragePath := "2025/12/19/new.jpg"
	repoError := errors.New("database update error")

	repo.On("GetByID", ctx, mediaID).Return(existingMedia, nil)
	storage.On("Save", ctx, newFilename, newFile).Return(newStoragePath, nil)
	repo.On("Update", ctx, mock.AnythingOfType("*media.Media")).Return(nil, repoError)
	storage.On("Delete", ctx, newStoragePath).Return(nil)

	result, err := uc.Execute(ctx, mediaID, newFilename, newFile)

	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	assert.Nil(t, result)

	repo.AssertExpectations(t)
	storage.AssertExpectations(t)
	// Old file should not be deleted if update fails
	storage.AssertNotCalled(t, "Delete", existingMedia.Path)
}
