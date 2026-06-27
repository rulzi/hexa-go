package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMediaUseCase(t *testing.T) {
	repo := &mockMediaRepository{}
	storage := &mockMediaStorage{}
	baseURL := "http://localhost:8080"

	uc := NewMediaUseCase(repo, storage, baseURL)

	assert.NotNil(t, uc)
	assert.NotNil(t, uc.Create)
	assert.NotNil(t, uc.Get)
	assert.NotNil(t, uc.List)
	assert.NotNil(t, uc.Update)
	assert.NotNil(t, uc.Delete)
}
