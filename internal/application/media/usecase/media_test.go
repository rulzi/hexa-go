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

func TestNewMediaUseCase(t *testing.T) {
	repo := &mockMediaRepository{}
	service := mediaservice.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewMediaUseCase(repo, service, storage, baseURL)

	assert.NotNil(t, uc)
	assert.Equal(t, repo, uc.mediaRepo)
	assert.Equal(t, service, uc.mediaService)
	assert.Equal(t, storage, uc.storage)
	assert.Equal(t, baseURL, uc.baseURL)
}

// --- Create ---

func TestMediaUseCase_Create_Success(t *testing.T) {
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

	result, err := uc.Create(ctx, filename, file)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedMedia.ID, result.ID)
	assert.Equal(t, expectedMedia.Name, result.Name)
	assert.Contains(t, result.URL, storagePath)
	storage.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestMediaUseCase_Create_StorageError(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := mediaservice.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewMediaUseCase(repo, service, storage, baseURL)

	filename := "test.jpg"
	file := strings.NewReader("content")
	storage.On("Save", ctx, filename, mock.Anything).Return("", errors.New("storage error"))

	result, err := uc.Create(ctx, filename, file)

	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertNotCalled(t, "Create")
}

// --- Get ---

func TestMediaUseCase_Get_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := mediaservice.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewMediaUseCase(repo, service, storage, baseURL)

	mediaID := int64(1)
	mediaEntity := &mediaentity.Media{
		ID: mediaID, Name: "test.jpg", Path: "2025/12/19/test.jpg",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	repo.On("GetByID", ctx, mediaID).Return(mediaEntity, nil)

	result, err := uc.Get(ctx, mediaID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mediaEntity.ID, result.ID)
	assert.Contains(t, result.URL, baseURL)
	repo.AssertExpectations(t)
}

func TestMediaUseCase_Get_MediaNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := mediaservice.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewMediaUseCase(repo, service, storage, baseURL)

	mediaID := int64(999)
	repo.On("GetByID", ctx, mediaID).Return(nil, nil)

	result, err := uc.Get(ctx, mediaID)

	assert.Error(t, err)
	assert.True(t, mediaentity.IsMediaNotFound(err))
	assert.Nil(t, result)
}

// --- List ---

func TestMediaUseCase_List_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := mediaservice.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewMediaUseCase(repo, service, storage, baseURL)

	limit, offset := 10, 0
	mediaList := []*mediaentity.Media{
		{ID: 1, Name: "a.jpg", Path: "path/a.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	total := int64(1)
	repo.On("List", ctx, limit, offset).Return(mediaList, nil)
	repo.On("Count", ctx).Return(total, nil)

	result, err := uc.List(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, total, result.Total)
	assert.Equal(t, limit, result.Limit)
	assert.Equal(t, offset, result.Offset)
	assert.Len(t, result.Media, 1)
}

func TestMediaUseCase_List_DefaultPagination(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := mediaservice.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewMediaUseCase(repo, service, storage, baseURL)

	repo.On("List", ctx, 10, 0).Return([]*mediaentity.Media{}, nil)
	repo.On("Count", ctx).Return(int64(0), nil)

	result, err := uc.List(ctx, -1, -1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 0, result.Offset)
}

// --- Update ---

func TestMediaUseCase_Update_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := mediaservice.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewMediaUseCase(repo, service, storage, baseURL)

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

	result, err := uc.Update(ctx, mediaID, "new.jpg", strings.NewReader("new content"))

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "new.jpg", result.Name)
	repo.AssertExpectations(t)
	storage.AssertExpectations(t)
}

func TestMediaUseCase_Update_MediaNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := mediaservice.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewMediaUseCase(repo, service, storage, baseURL)

	mediaID := int64(999)
	repo.On("GetByID", ctx, mediaID).Return(nil, nil)

	result, err := uc.Update(ctx, mediaID, "new.jpg", strings.NewReader("x"))

	assert.Error(t, err)
	assert.True(t, mediaentity.IsMediaNotFound(err))
	assert.Nil(t, result)
	storage.AssertNotCalled(t, "Save")
}

// --- Delete ---

func TestMediaUseCase_Delete_Success(t *testing.T) {
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

	err := uc.Delete(ctx, mediaID)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	storage.AssertExpectations(t)
}

func TestMediaUseCase_Delete_MediaNotFound(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	service := mediaservice.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewMediaUseCase(repo, service, storage, baseURL)

	mediaID := int64(999)
	repo.On("GetByID", ctx, mediaID).Return(nil, nil)

	err := uc.Delete(ctx, mediaID)

	assert.Error(t, err)
	assert.True(t, mediaentity.IsMediaNotFound(err))
	storage.AssertNotCalled(t, "Delete")
	repo.AssertNotCalled(t, "Delete")
}
