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

func TestNewCreateMediaUseCase(t *testing.T) {
	repo := &mockMediaRepository{}
	service := domainmedia.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewCreateMediaUseCase(repo, service, storage, baseURL)

	assert.NotNil(t, uc)
	assert.Equal(t, repo, uc.mediaRepo)
	assert.Equal(t, service, uc.mediaService)
	assert.Equal(t, storage, uc.storage)
	assert.Equal(t, baseURL, uc.baseURL)
}

func TestCreateMediaUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := domainmedia.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewCreateMediaUseCase(repo, service, storage, baseURL)

	filename := "test.jpg"
	file := strings.NewReader("test file content")
	storagePath := "2025/12/19/test.jpg"

	expectedMedia := &domainmedia.Media{
		ID:        1,
		Name:      filename,
		Path:      storagePath,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	storage.On("Save", ctx, filename, file).Return(storagePath, nil)
	repo.On("Create", ctx, mock.AnythingOfType("*media.Media")).Return(expectedMedia, nil)

	result, err := uc.Execute(ctx, filename, file)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedMedia.ID, result.ID)
	assert.Equal(t, expectedMedia.Name, result.Name)
	assert.Equal(t, expectedMedia.Path, result.Path)
	assert.Contains(t, result.URL, storagePath)
	assert.Contains(t, result.URL, baseURL)

	storage.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestCreateMediaUseCase_Execute_StorageError(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := domainmedia.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewCreateMediaUseCase(repo, service, storage, baseURL)

	filename := "test.jpg"
	file := strings.NewReader("test file content")
	storageError := errors.New("storage error")

	storage.On("Save", ctx, filename, file).Return("", storageError)

	result, err := uc.Execute(ctx, filename, file)

	assert.Error(t, err)
	assert.Equal(t, storageError, err)
	assert.Nil(t, result)

	storage.AssertExpectations(t)
	repo.AssertNotCalled(t, "Create")
}

func TestCreateMediaUseCase_Execute_ValidationError(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := domainmedia.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewCreateMediaUseCase(repo, service, storage, baseURL)

	// Test with empty filename (will cause validation error after storage save)
	filename := ""
	file := strings.NewReader("test file content")
	storagePath := "2025/12/19/.jpg"

	storage.On("Save", ctx, filename, file).Return(storagePath, nil)
	storage.On("Delete", ctx, storagePath).Return(nil)

	result, err := uc.Execute(ctx, filename, file)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domainmedia.ErrNameRequired, err)

	storage.AssertExpectations(t)
	repo.AssertNotCalled(t, "Create")
}

func TestCreateMediaUseCase_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := domainmedia.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewCreateMediaUseCase(repo, service, storage, baseURL)

	filename := "test.jpg"
	file := strings.NewReader("test file content")
	storagePath := "2025/12/19/test.jpg"
	repoError := errors.New("database error")

	storage.On("Save", ctx, filename, file).Return(storagePath, nil)
	repo.On("Create", ctx, mock.AnythingOfType("*media.Media")).Return(nil, repoError)
	storage.On("Delete", ctx, storagePath).Return(nil)

	result, err := uc.Execute(ctx, filename, file)

	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	assert.Nil(t, result)

	storage.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestCreateMediaUseCase_Execute_URIBuilding(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := domainmedia.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewCreateMediaUseCase(repo, service, storage, baseURL)

	filename := "test.jpg"
	file := strings.NewReader("test file content")
	storagePath := "2025/12/19/test.jpg"

	expectedMedia := &domainmedia.Media{
		ID:        1,
		Name:      filename,
		Path:      storagePath,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	storage.On("Save", ctx, filename, file).Return(storagePath, nil)
	repo.On("Create", ctx, mock.AnythingOfType("*media.Media")).Return(expectedMedia, nil)

	result, err := uc.Execute(ctx, filename, file)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "http://localhost:8080/api/v1/media/files/2025/12/19/test.jpg", result.URL)

	storage.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestCreateMediaUseCase_Execute_BaseURLWithTrailingSlash(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := domainmedia.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080/"

	uc := NewCreateMediaUseCase(repo, service, storage, baseURL)

	filename := "test.jpg"
	file := strings.NewReader("test file content")
	storagePath := "2025/12/19/test.jpg"

	expectedMedia := &domainmedia.Media{
		ID:        1,
		Name:      filename,
		Path:      storagePath,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	storage.On("Save", ctx, filename, file).Return(storagePath, nil)
	repo.On("Create", ctx, mock.AnythingOfType("*media.Media")).Return(expectedMedia, nil)

	result, err := uc.Execute(ctx, filename, file)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	// Should not have double slashes
	assert.NotContains(t, result.URL, "//api")
	assert.Equal(t, "http://localhost:8080/api/v1/media/files/2025/12/19/test.jpg", result.URL)

	storage.AssertExpectations(t)
	repo.AssertExpectations(t)
}
