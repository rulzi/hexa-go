package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	mediaentity "github.com/rulzi/hexa-go/internal/domain/media/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateMedia_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewCreateMedia(repo, storage, baseURL)

	filename := "test.jpg"
	file := strings.NewReader("test file content")
	storagePath := "2025/12/19/test.jpg"
	expectedMedia := &mediaentity.Media{
		ID: 1, Name: filename, Path: storagePath,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	storage.On("Save", ctx, filename, mock.Anything).Return(storagePath, nil)
	repo.On("Create", ctx, mock.Anything).Return(expectedMedia, nil)

	result, err := uc.Execute(ctx, filename, file)

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
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewCreateMedia(repo, storage, baseURL)

	filename := "test.jpg"
	file := strings.NewReader("content")
	storage.On("Save", ctx, filename, mock.Anything).Return("", errors.New("storage error"))

	result, err := uc.Execute(ctx, filename, file)

	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertNotCalled(t, "Create")
}

func TestCreateMedia_Execute_RepoError(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewCreateMedia(repo, storage, baseURL)

	filename := "test.jpg"
	file := strings.NewReader("content")
	storagePath := "2025/01/01/test.jpg"

	storage.On("Save", ctx, filename, mock.Anything).Return(storagePath, nil)
	storage.On("Delete", ctx, storagePath).Return(nil)
	repo.On("Create", ctx, mock.Anything).Return(nil, errors.New("db error"))

	result, err := uc.Execute(ctx, filename, file)

	assert.Error(t, err)
	assert.Nil(t, result)
	storage.AssertExpectations(t)
}

func TestCreateMedia_Execute_ValidationError(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewCreateMedia(repo, storage, baseURL)

	storagePath := "2025/01/01/test.jpg"
	storage.On("Save", ctx, "", mock.Anything).Return(storagePath, nil)
	storage.On("Delete", ctx, storagePath).Return(nil)

	result, err := uc.Execute(ctx, "", strings.NewReader("content"))

	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertNotCalled(t, "Create")
}
