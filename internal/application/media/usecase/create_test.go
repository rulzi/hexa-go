package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	mediaentity "github.com/rulzi/hexa-go/internal/domain/media/entity"
	mediaservice "github.com/rulzi/hexa-go/internal/domain/media/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateMedia_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := mediaservice.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewMediaUseCase(repo, service, storage, baseURL)

	filename := "test.jpg"
	file := strings.NewReader("test file content")
	storagePath := "2025/12/19/test.jpg"
	expectedMedia := &mediaentity.Media{
		ID: 1, Name: filename, Path: storagePath,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	storage.On("Save", ctx, filename, mock.Anything).Return(storagePath, nil)
	repo.On("Create", ctx, mock.Anything).Return(expectedMedia, nil)

	result, err := uc.Create.Execute(ctx, filename, file)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedMedia.ID, result.ID)
	assert.Equal(t, expectedMedia.Name, result.Name)
	assert.Contains(t, result.URL, storagePath)
	storage.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestCreateMedia_Execute_StorageError(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := mediaservice.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewMediaUseCase(repo, service, storage, baseURL)

	filename := "test.jpg"
	file := strings.NewReader("content")
	storage.On("Save", ctx, filename, mock.Anything).Return("", errors.New("storage error"))

	result, err := uc.Create.Execute(ctx, filename, file)

	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertNotCalled(t, "Create")
}
