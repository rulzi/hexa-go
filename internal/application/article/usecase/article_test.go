package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewArticleUseCase(t *testing.T) {
	repo := &mockArticleRepository{}
	cache := &mockArticleCache{}
	listCache := &mockArticleListCache{}

	uc := NewArticleUseCase(repo, cache, listCache)

	assert.NotNil(t, uc)
	assert.NotNil(t, uc.Create)
	assert.NotNil(t, uc.Get)
	assert.NotNil(t, uc.List)
	assert.NotNil(t, uc.Update)
	assert.NotNil(t, uc.Delete)
}
