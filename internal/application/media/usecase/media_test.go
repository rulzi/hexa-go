package usecase

import (
	"testing"

	mediaservice "github.com/rulzi/hexa-go/internal/domain/media/service"
	"github.com/stretchr/testify/assert"
)

func TestNewMediaUseCase(t *testing.T) {
	repo := &mockMediaRepository{}
	service := mediaservice.NewService(repo)
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewMediaUseCase(repo, service, storage, baseURL)

	assert.NotNil(t, uc)
	assert.NotNil(t, uc.Create)
	assert.NotNil(t, uc.Get)
	assert.NotNil(t, uc.List)
	assert.NotNil(t, uc.Update)
	assert.NotNil(t, uc.Delete)
}
