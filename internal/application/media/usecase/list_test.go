package usecase

import (
	"context"
	"testing"
	"time"

	mediaentity "github.com/rulzi/hexa-go/internal/domain/media/entity"
	"github.com/stretchr/testify/assert"
)

func TestListMedia_Execute_Success(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewMediaUseCase(repo, storage, baseURL)

	limit, offset := 10, 0
	mediaList := []*mediaentity.Media{
		{ID: 1, Name: "a.jpg", Path: "path/a.jpg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	total := int64(1)
	repo.On("List", ctx, limit, offset).Return(mediaList, nil)
	repo.On("Count", ctx).Return(total, nil)

	result, err := uc.List.Execute(ctx, limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, total, result.Total)
	assert.Equal(t, limit, result.Limit)
	assert.Equal(t, offset, result.Offset)
	assert.Len(t, result.Media, 1)
}

func TestListMedia_Execute_DefaultPagination(t *testing.T) {
	ctx := context.Background()
	repo := &mockMediaRepository{}
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewMediaUseCase(repo, storage, baseURL)

	repo.On("List", ctx, 10, 0).Return([]*mediaentity.Media{}, nil)
	repo.On("Count", ctx).Return(int64(0), nil)

	result, err := uc.List.Execute(ctx, -1, -1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 0, result.Offset)
}
